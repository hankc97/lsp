package tcpserver

import (
	"context"
	"lsp/tcpserver/parse"
	"encoding/json"
	"github.com/pkg/errors"
	"kythe.io/kythe/go/languageserver"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
)

func Initialize(ctx context.Context, body *parse.LspBody, server languageserver.Server) (*lsp.InitializeResult, error) {
	params := body.Params 
	initializeParamStruct := lsp.InitializeParams{}
	err := json.Unmarshal(params, &initializeParamStruct)
	if err != nil {
		return nil, errors.New("decoding lsp body params")
	}

	initializeResult, err := server.Initialize(initializeParamStruct)
	if err != nil {
		return nil, errors.New("decoding initialized params")
	}

	return initializeResult, nil;
}
