package tcpserver

import (
	"log"
	"lsp/server/parse"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"kythe.io/kythe/go/languageserver"
)

func Initialize(body *parse.LspBody, server languageserver.Server) (*lsp.InitializeResult, error) {
	params := body.Params 
	initializeParamStruct := lsp.InitializeParams{}
	err := json.Unmarshal(params, &initializeParamStruct)
	if err != nil {
		log.Println("decoding lsp body params")
		return nil, errors.New("decoding lsp body params")
	}

	initializeResult, err := server.Initialize(initializeParamStruct)
	if err != nil {
		log.Println("decoding initialized params")
		return nil, errors.New("decoding initialized params")
	}

	return initializeResult, nil;
}
