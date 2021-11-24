package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	tcpserver "lsp/server"
	"lsp/server/parse"
	"net"
	"os"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"kythe.io/kythe/go/languageserver"
	"kythe.io/kythe/go/services/xrefs"
	"encoding/hex"
)

const (
	maxContentLength = 1 << 20
)

const (
	serverInitialize  string = "initialize"
	serverInitialized string = "initialized"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	iface := flag.String("iface", "127.0.0.1", "interface to bind to, defaults to localhost")
	port := flag.String("port", "", "port to bind to")
	logfile := flag.String("logfile", "", "also log to this file (in additional to stderr) under logger/")
	flag.Parse()

	if iface == nil || *iface == "" {
		return errors.New("-iface is required")
	}

	if port == nil || *port == "" {
		return errors.New("-port is required")
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(*iface, *port))
	if err != nil {
		return errors.Wrap(err, "creating listener")
	}
	defer listener.Close()

	var logWriter io.Writer
	if *logfile == "" {
		logWriter = os.Stderr
	} else {
		logDir := fmt.Sprint("logger/", *logfile)
		file, err := os.Create(logDir)
		if err != nil {
			return err
		}
		defer file.Close()
		logWriter = io.MultiWriter(os.Stderr, file)
	}
	log.SetOutput(logWriter)

	addr := listener.Addr().(*net.TCPAddr)

	log.Printf("listening on %q", addr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err, "accepting client connection")
			return errors.Wrap(err, "accepting client connection")
		}
		go func() {
			err := handleClientConn(conn)
			if err != nil {
				log.Printf("handling client: %v", err)
			}
		}()
	}
}

func handleClientConn(conn io.ReadWriteCloser) error {
	defer conn.Close()

	fmt.Print("\n")
	var xref xrefs.Service
	options := &languageserver.Options{}
	server := languageserver.NewServer(xref, options)

	more := true
	for more {
		req, last, err := parseRequest(conn)
		// req, last, err := parseRequest(io.TeeReader(conn, os.Stderr))
		if err != nil {
			log.Println(err, "parsing request")
			return errors.Wrap(err, "parsing request")
		}

		if last {
			more = false
		}

		// handle request and respond
		if err := serveReq(conn, req, server); err != nil {
			log.Println(err, "serving request")
			return errors.Wrap(err, "serving request")
		}
	}
	return nil
}

func serveReq(conn io.Writer, req *parse.LspRequest, server languageserver.Server) error {
	// write to `resp` according to what `req` contains
	body := req.Body
	var result interface{}
	var err error

	switch body.Method {
	case serverInitialize:
		result, err = tcpserver.Initialize(body, server)
	case serverInitialized:
	default:
		err = errors.Errorf("unsupported method: %q", body.Method)
	}
	if err != nil {
		return errors.Wrap(err, "handling method")
	}

	response, err := NewResponse(body.Id, result, err)
	if err != nil {
		fmt.Println(response, "ERROR!!!!")
		return errors.Wrap(err, "preparing response")
	}

	requestBody := Resp {
		Jsonrpc: "2.0",
		Id: 0,
		Result: ResultValue {
			Capabilities: CapabilitiesValue{
				TextDocumentSync: 2,
				CompletionProvider: ResolveProviderValue{
					ResolveProvider: true,
				},
				Workspace: WorkspaceValue{
					WorkspaceFolders: WorkspaceFoldersValue{
						Supported: true,
					},
				},
			},
		},
	}

	marshalledBodyRequest, err := json.Marshal(&requestBody)
	if err != nil {
		return errors.Wrap(err, "marshaling response body")
	}
	// contentLengthRespBody := int(len(marshalledBodyRequest))

	// requestMid := "\r\n"
	// marshalledMidRequest, err := json.Marshal(&requestMid)
	// if err != nil {
	// 	return errors.Wrap(err, "marshaling response header")
	// }

	// requestHeader := fmt.Sprint("Content-Length: ", contentLengthRespBody, "%0D%0A")
	// str := ``
	// buf := bytes.NewBufferString(str).Bytes()

	// fmt.Println(buf)
	
	// CR LF -> %0D%0A to seperate header and body
	s := fmt.Sprintf("%x", "Content-Length: 171\r\n\r\n")
	// returns binary value of "string" + %0D%0A
	b, err := hex.DecodeString(s) 


	// buf := append(buf)
	
	marshalledHeaderRequest, err := json.Marshal(&requestBody)
	marshalledHeaderRequest = json.RawMessage(marshalledHeaderRequest)
	// fmt.Println(marshalledHeaderRequest)

	if err != nil {
		return errors.Wrap(err, "marshaling response header")
	}

	log.Print("\n")
	// log.Println("stringified", string(marshalledHeaderRequest))

	// fmt.Println(buf)
	// fmt.Println(buf.Bytes())
	if _, err := conn.Write(b); err != nil {
		return errors.Wrap(err, "writing response to connection")
	}
	// time.Sleep(4 * time.Second)

	log.Print("\n")
	log.Println("stringified", string(marshalledBodyRequest))
	if _, err := conn.Write(marshalledBodyRequest); err != nil {
		return errors.Wrap(err, "writing response to connection")
	}
	// time.Sleep(4 * time.Second)


	return nil
}

type Resp struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  ResultValue `json:"result"`
}

type ResultValue struct {
	Capabilities CapabilitiesValue `json:"capabilities"`
}

type CapabilitiesValue struct {
	TextDocumentSync   int                  `json:"textDocumentSync"`
	CompletionProvider ResolveProviderValue `json:"completionProvider"`
	Workspace          WorkspaceValue       `json:"workspace"`
}

type ResolveProviderValue struct {
	ResolveProvider bool `json:"resolveProvider"`
}

type WorkspaceValue struct {
	WorkspaceFolders WorkspaceFoldersValue `json:"workspaceFolders"`
}

type WorkspaceFoldersValue struct {
	Supported bool `json:"supported"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	// result is the content of the response.
	Result json.RawMessage `json:"result"`
	// err is set only if the call failed.
	Error error `json:"error"`
	// ID of the request this is a response to.
	Id int `json:"id"`
}

func NewResponse(id int, result interface{}, err error) (*Response, error) {
	r, merr := marshalInterface(result)
	resp := &Response{
		Jsonrpc: "2.0",
		Result:  r,
		Error:   err,
		Id:      id,
	}
	return resp, merr
}

func marshalInterface(obj interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("failed to marshal json: %w", err)
		return json.RawMessage{}, fmt.Errorf("failed to marshal json: %w", err)
	}
	return json.RawMessage(data), nil
}

func parseRequest(in io.Reader) (_ *parse.LspRequest, last bool, err error) {
	header, err := parseHeader(in)
	if err != nil {
		log.Println(err, "parsing header")
		return nil, false, errors.Wrap(err, "parsing header")
	}

	switch header.ContentType {
	case "application/vscode-jsonrpc; charset=utf-8":
		// continue
	case "":

	default:
		return nil, false, errors.Errorf("unsupported content type: %q", header.ContentType)
	}

	parsedBody, err := parseBody(in, header.ContentLength)
	if err != nil {
		log.Println(err, "parsing body")
		return nil, false, errors.Wrap(err, "parsing body")
	}

	body := new(parse.LspBody)
	err = json.Unmarshal([]byte(parsedBody), &body)

	switch err {
	case io.EOF:
		// no more requests are coming
		last = true
	case nil:
		// no problem
	default:
		log.Println(err, "decoding body")
		return nil, false, errors.Wrap(err, "decoding body")
	}

	// do something with `req`
	return &parse.LspRequest{Header: header, Body: body}, last, nil
}

func parseHeader(in io.Reader) (*parse.LspHeader, error) {
	var lsp parse.LspHeader
	scan := bufio.NewScanner(in)
	fmt.Println("received header... ")

	for scan.Scan() {
		header := scan.Text()
		fmt.Println(header)
		if header == "" {
			// last header
			return &lsp, nil
		}
		name, value, err := splitOnce(header, ": ")
		if err != nil {
			log.Println(err, "parsing an header entry")
			return nil, errors.Wrap(err, "parsing an header entry")
		}
		switch name {
		case "Content-Length":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Println(err, "invalid Content-Length: %q", value)
				return nil, errors.Wrapf(err, "invalid Content-Length: %q", value)
			}
			lsp.ContentLength = v
		case "Content-Type":
			lsp.ContentType = value
		}
	}
	if err := scan.Err(); err != nil {
		log.Println(err, "scanning header entries")
		return nil, errors.Wrap(err, "scanning header entries")
	}
	log.Println("no body contained")
	return nil, errors.New("no body contained")
}

func splitOnce(in, sep string) (prefix, suffix string, err error) {
	sepIdx := strings.Index(in, sep)
	if sepIdx < 0 {
		log.Printf("separator %q not found", sep)
		return "", "", errors.Errorf("separator %q not found", sep)
	}
	prefix = in[:sepIdx]
	suffix = in[sepIdx+len(sep):]
	return prefix, suffix, nil
}

func parseBody(in io.Reader, contentLength int64) (string, error) {
	var body string
	scanner := bufio.NewScanner(in)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF || (len(data) == int(contentLength)) {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		if len(scanner.Bytes()) == int(contentLength) {
			body = scanner.Text()
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err, "scanning body entries")
		return "", errors.Wrap(err, "scanning body entries")
	}

	fmt.Println("body is printed here")

	return body, nil
}
