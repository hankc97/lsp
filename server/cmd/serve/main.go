package main

import (
	"bufio"
	"context"
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
)

const (
	maxContentLength = 1 << 20
)

const (
	serverInitialize  string = "initialize"
	serverInitialized string = "initialized"
)

var (
	Iface   = flag.String("iface", "127.0.0.1", "interface to bind to, defaults to localhost")
	port    = flag.String("port", "", "port to bind to")
	c       = flag.String("c", "", "communication mode (stdio|tcp)")
	logfile = flag.String("logfile", "", "also log to this file (in additional to stderr) under logger/")
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {

	flag.Parse()

	if Iface == nil || *Iface == "" {
		return errors.New("-iface is required")
	}

	if port == nil || *port == "" {
		return errors.New("-port is required")
	}

	if c == nil || *c == "" {
		return errors.New("-c is required: stdio|tcp")
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(*Iface, *port))
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

	conn, err := listener.Accept()
	if err != nil {
		log.Println(err, "accepting client connection")
		return errors.Wrap(err, "accepting client connection")
	}
	var xref xrefs.Service
	options := &languageserver.Options{}
	server := languageserver.NewServer(xref, options)
	return handleClientConn(conn, server)
}

func handleClientConn(conn io.ReadWriteCloser, server languageserver.Server) error {
	defer conn.Close()

	more := true
	for more {

		// req, err := readRequest(conn)
		// if err != nil {
		// 	return err
		// }
		// writeRequest(conn, req)

		req, last, err := parseRequest(conn)
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
	ctx := context.Background()
	body := req.Body
	var result interface{}
	var err error

	switch body.Method {
	case serverInitialize:
		result, err = tcpserver.Initialize(ctx, body, server)
	case serverInitialized:
	default:
		log.Println("invalid method type")
		errors.New("invalid method type")
	}

	response, err := NewResponse(body.Id, result, err)
	if err != nil {
		return err
	}

	marshalledResponse, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	log.Print("\n")
	log.Println("sending...", string(marshalledResponse))
	log.Println("sending...", marshalledResponse)

	conn.Write(marshalledResponse)
	return nil
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
		log.Println(header)

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

	log.Println("received body... ")
	log.Println(body)
	log.Print("\n")

	return body, nil
}



// testing on netcat for stdin / stdio
func readRequest(reader io.Reader) (string, error) {
	scan := bufio.NewScanner(reader)
	var req string

	for scan.Scan() {
		headerContext := scan.Text()
		fmt.Printf("%+v\n", headerContext)

		if headerContext == "Content-Length: 5000" {
			req = headerContext
		}

		if headerContext == "" {
			return req, nil
		}
	}

	return "nothing exists", nil
}

func writeRequest(writer io.Writer, req string) error {
	fmt.Println(req, "written back to client")

	marshalledreq, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	writer.Write(marshalledreq)

	return nil
}