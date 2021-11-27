package tcpserver

import (
	"log"
	"lsp/server/parse"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"kythe.io/kythe/go/languageserver"
)

func Initialize(body *parse.LspBody, server languageserver.Server) (*ResultValue, error) {
	params := body.Params 
	initializeParamStruct := lsp.InitializeParams{}
	err := json.Unmarshal(params, &initializeParamStruct)
	if err != nil {
		log.Println("decoding lsp body params")
		return nil, errors.New("decoding lsp body params")
	}

	// initializeResult, err := server.Initialize(initializeParamStruct)
	if err != nil {
		log.Println("decoding initialized params")
		return nil, errors.New("decoding initialized params")
	}

	result := ResultValue {
		Capabilities: CapabilitiesValue {
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
	}

	return &result, nil;
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



// requestBody := Resp{
// 	Jsonrpc: "2.0",
// 	Id:      0,
// 	Result: ResultValue{
// 		Capabilities: CapabilitiesValue{
// 			TextDocumentSync: 2,
// 			CompletionProvider: ResolveProviderValue{
// 				ResolveProvider: true,
// 			},
// 			Workspace: WorkspaceValue{
// 				WorkspaceFolders: WorkspaceFoldersValue{
// 					Supported: true,
// 				},
// 			},
// 		},
// 	},
// }