package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	// "bufio"
	// "lsp/tcpserver"
)

const (
	maxContentLength = 1 << 20
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
	return handleClientConn(conn)
}

func handleClientConn(conn io.ReadWriteCloser) error {
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
		if err := serveReq(conn, req); err != nil {
			return errors.Wrap(err, "serving request")
		}
	}
	return nil
}

func serveReq(resp io.Writer, req *lspRequest) error {
	// write to `resp` according to what `req` contains
	return nil
}

func parseRequest(in io.Reader) (_ *lspRequest, last bool, err error) {
	header, err := parseHeader(in)
	if err != nil {
		return nil, false, errors.Wrap(err, "parsing header")
	}

	switch header.ContentType {
	case "application/vscode-jsonrpc; charset=utf-8":
		// continue
	default:
		return nil, false, errors.Errorf("unsupported content type: %q", header.ContentType)
	}
	bodyReader := io.LimitReader(in, imin(header.ContentLength, maxContentLength))

	dump, err := ioutil.ReadAll(in)
	panic(fmt.Sprintf("%q", dump))
	panic(err)

	var body *lspBody
	err = json.NewDecoder(bodyReader).Decode(&body)
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

	return &lspRequest{Header: header, Body: body}, last, nil
}

func imin(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type lspRequest struct {
	Header *lspHeader
	Body   *lspBody
}

type lspHeader struct {
	ContentLength int64
	ContentType   string
}

func parseHeader(in io.Reader) (*lspHeader, error) {
	var lsp lspHeader
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

type lspBody struct {
	// stuff goes here
}
