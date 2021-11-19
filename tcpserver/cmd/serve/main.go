package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"lsp/tcpserver/parse"
	"lsp/tcpserver"
	"context"
	"fmt"
	"kythe.io/kythe/go/languageserver"
	"kythe.io/kythe/go/services/xrefs"
)

const (
	maxContentLength = 1 << 20
)

const (
	serverInitialize string = "initialize"
	serverInitialized string = "initialized"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	iface := flag.String("iface", "127.0.0.1", "interface to bind to, defaults to localhost")
	port := flag.String("port", "5007", "port to bind to")
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

	addr := listener.Addr().(*net.TCPAddr)

	log.Printf("listening on %q", addr.String())

	conn, err := listener.Accept()
	if err != nil {
		return errors.Wrap(err, "accepting client connection")
	}
	var xref xrefs.Service
	options := &languageserver.Options{}
	server := languageserver.NewServer(xref , options)
	return handleClientConn(conn, server)
}

func handleClientConn(conn io.ReadWriteCloser, server languageserver.Server) error {
	defer conn.Close()

	more := true
	for more {
		req, last, err := parseRequest(conn)
		if err != nil {
			return errors.Wrap(err, "parsing request")
		}

		if last {
			more = false
		}

		// handle request and respond
		if err := serveReq(conn, req, server); err != nil {
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
		errors.New("invalid method type")
	}

	response, err := NewResponse(body.Id, result, err)
	if err != nil {
		return err
	}
	
	// fmt.Println("jsonrpc: ",response.jsonrpc)
	// fmt.Println("result: ",string(response.result))
	// fmt.Println("err: ",response.err)
	// fmt.Println("id: ",response.id)

	marshalledResponse, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	conn.Write(marshalledResponse)
	return nil
}

type Response struct {
	jsonrpc string
    // result is the content of the response.
    result json.RawMessage
    // err is set only if the call failed.
    err error
    // ID of the request this is a response to.
    id int
}

func NewResponse(id int, result interface{}, err error) (*Response, error) {
	r, merr := marshalInterface(result)
	resp := &Response{
		jsonrpc: "2.0",
		result: r,
		err:    err,
		id:     id,
	}
	return resp, merr
}

func marshalInterface(obj interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("failed to marshal json: %w", err)
	}
	return json.RawMessage(data), nil
}

func parseRequest(in io.Reader) (_ *parse.LspRequest, last bool, err error) {
	header, err := parseHeader(in)
	if err != nil {
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
		return nil, false, errors.Wrap(err, "decoding body")
	}

	// do something with `req`
	return &parse.LspRequest{Header: header, Body: body}, last, nil
}


func parseHeader(in io.Reader) (*parse.LspHeader, error) {
	var lsp parse.LspHeader
	scan := bufio.NewScanner(in)

	for scan.Scan() {
		header := scan.Text()
		if header == "" {
			// last header
			return &lsp, nil
		}
		name, value, err := splitOnce(header, ": ")
		if err != nil {
			return nil, errors.Wrap(err, "parsing an header entry")
		}
		switch name {
		case "Content-Length":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid Content-Length: %q", value)
			}
			lsp.ContentLength = v
		case "Content-Type":
			lsp.ContentType = value
		}
	}
	if err := scan.Err(); err != nil {
		return nil, errors.Wrap(err, "scanning header entries")
	}
	return nil, errors.New("no body contained")
}

func splitOnce(in, sep string) (prefix, suffix string, err error) {
	sepIdx := strings.Index(in, sep)
	if sepIdx < 0 {
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
		return "", errors.Wrap(err, "scanning body entries")
	}

	return body, nil
}

