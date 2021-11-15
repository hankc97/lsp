package tcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
	"fmt"
)

type method func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{}
type methodMap map[string]method

type server struct {
	RootURI string
	files   map[string]string
}

func (server *server) initialize(ctx context.Context, conn jsonrpc2.JSONRPC2, params lsp.InitializeParams) (*lsp.InitializeResult, *lsp.InitializeError) {
	server.RootURI = string(params.RootURI)
	server.files = map[string]string{}

	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKFull,
				},
			},
		},
	}, nil
}

func To(fn interface{}) func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{} {

	fmt.Println("------------------")

	val := reflect.ValueOf(fn)
	fmt.Println(val)

	in := val.Type().In(2)


	return func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{} {
		v := reflect.New(in)
		json.Unmarshal(params, v.Interface())
		ret := val.Call([]reflect.Value{
			reflect.ValueOf(ctx), reflect.ValueOf(conn), v.Elem(),
		})
		if len(ret) == 0 {
			return nil
		} else {
			if !ret[0].IsNil() {
				return ret[0].Interface()
			}
			if !ret[1].IsNil() {
				return ret[1].Interface()
			}
			panic("err")
		}
	}
}

type stReadWriteClose struct{}

func (stReadWriteClose) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stReadWriteClose) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stReadWriteClose) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	return os.Stdout.Close()
}

func StartServer() {
	server := server{}
	documents := methodMap{
		"initialize": To(server.initialize),
	}

	handler := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
		value, ok := documents[req.Method]
		if !ok {
			return nil, errors.New("not found")
		}
		response := value(ctx, conn, *req.Params)

		return response, nil
	})
	<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stReadWriteClose{}, jsonrpc2.VSCodeObjectCodec{}), handler).DisconnectNotify()
}
