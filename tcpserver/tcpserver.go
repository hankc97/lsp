package tcpserver

import (
	"context"
	"lsp/tcpserver/parse"
	"go.lsp.dev/protocol"
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
)

func Initialize(ctx context.Context, body *parse.LspBody) (protocol.InitializeParams, error) {
	params := body.Params 
	initializeParamStruct := protocol.InitializeParams{}
	err := json.Unmarshal(params, &initializeParamStruct)
	if err != nil {
		return initializeParamStruct, errors.New("decoding lsp body params")
	}
	fmt.Printf("%+v\n", initializeParamStruct)
	return initializeParamStruct, nil;
}
