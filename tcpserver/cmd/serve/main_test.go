// https://microsoft.github.io/language-server-protocol/specification.html#contentPart
package main

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/require"
)



func TestParseRequest(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *lspRequest
	}{
		{
			name: "base case",
			input: "Content-Length: 4792\r\n" +
					"\r\n" +
					`{
						"jsonrpc": "2.0",
						"id": 0,
						"method": "initialize",
						"params": {
							"processId": 1344,
							"clientInfo": {
								"name": "Visual Studio Code",
								"version": "1.62.2"
							},
							"locale": "en",
							"rootPath": "/home/hank/CodingWork/Go/github.com/github-actions/testplaintext",
							"rootUri": "file:///home/hank/CodingWork/Go/github.com/github-actions/testplaintext",
							"capabilities": {
								"workspace": {
									"applyEdit": true,
									"workspaceEdit": {
										"documentChanges": true,
										"resourceOperations": [
											"create",
											"rename",
											"delete"
										],
										"failureHandling": "textOnlyTransactional",
										"normalizesLineEndings": true,
										"changeAnnotationSupport": {
											"groupsOnLabel": true
										}
									},
									"didChangeConfiguration": {
										"dynamicRegistration": true
									},
									"didChangeWatchedFiles": {
										"dynamicRegistration": true
									},
									"symbol": {
										"dynamicRegistration": true,
										"symbolKind": {
											"valueSet": [
												1,
												2,
												3,
												4,
												5,
												6,
												7,
												8,
												9,
												10,
												11,
												12,
												13,
												14,
												15,
												16,
												17,
												18,
												19,
												20,
												21,
												22,
												23,
												24,
												25,
												26
											]
										},
										"tagSupport": {
											"valueSet": [
												1
											]
										}
									},
									"codeLens": {
										"refreshSupport": true
									},
									"executeCommand": {
										"dynamicRegistration": true
									},
									"configuration": true,
									"workspaceFolders": true,
									"semanticTokens": {
										"refreshSupport": true
									},
									"fileOperations": {
										"dynamicRegistration": true,
										"didCreate": true,
										"didRename": true,
										"didDelete": true,
										"willCreate": true,
										"willRename": true,
										"willDelete": true
									}
								},
								"textDocument": {
									"publishDiagnostics": {
										"relatedInformation": true,
										"versionSupport": false,
										"tagSupport": {
											"valueSet": [
												1,
												2
											]
										},
										"codeDescriptionSupport": true,
										"dataSupport": true
									},
									"synchronization": {
										"dynamicRegistration": true,
										"willSave": true,
										"willSaveWaitUntil": true,
										"didSave": true
									},
									"completion": {
										"dynamicRegistration": true,
										"contextSupport": true,
										"completionItem": {
											"snippetSupport": true,
											"commitCharactersSupport": true,
											"documentationFormat": [
												"markdown",
												"plaintext"
											],
											"deprecatedSupport": true,
											"preselectSupport": true,
											"tagSupport": {
												"valueSet": [
													1
												]
											},
											"insertReplaceSupport": true,
											"resolveSupport": {
												"properties": [
													"documentation",
													"detail",
													"additionalTextEdits"
												]
											},
											"insertTextModeSupport": {
												"valueSet": [
													1,
													2
												]
											}
										},
										"completionItemKind": {
											"valueSet": [
												1,
												2,
												3,
												4,
												5,
												6,
												7,
												8,
												9,
												10,
												11,
												12,
												13,
												14,
												15,
												16,
												17,
												18,
												19,
												20,
												21,
												22,
												23,
												24,
												25
											]
										}
									},
									"hover": {
										"dynamicRegistration": true,
										"contentFormat": [
											"markdown",
											"plaintext"
										]
									},
									"signatureHelp": {
										"dynamicRegistration": true,
										"signatureInformation": {
											"documentationFormat": [
												"markdown",
												"plaintext"
											],
											"parameterInformation": {
												"labelOffsetSupport": true
											},
											"activeParameterSupport": true
										},
										"contextSupport": true
									},
									"definition": {
										"dynamicRegistration": true,
										"linkSupport": true
									},
									"references": {
										"dynamicRegistration": true
									},
									"documentHighlight": {
										"dynamicRegistration": true
									},
									"documentSymbol": {
										"dynamicRegistration": true,
										"symbolKind": {
											"valueSet": [
												1,
												2,
												3,
												4,
												5,
												6,
												7,
												8,
												9,
												10,
												11,
												12,
												13,
												14,
												15,
												16,
												17,
												18,
												19,
												20,
												21,
												22,
												23,
												24,
												25,
												26
											]
										},
										"hierarchicalDocumentSymbolSupport": true,
										"tagSupport": {
											"valueSet": [
												1
											]
										},
										"labelSupport": true
									},
									"codeAction": {
										"dynamicRegistration": true,
										"isPreferredSupport": true,
										"disabledSupport": true,
										"dataSupport": true,
										"resolveSupport": {
											"properties": [
												"edit"
											]
										},
										"codeActionLiteralSupport": {
											"codeActionKind": {
												"valueSet": [
													"",
													"quickfix",
													"refactor",
													"refactor.extract",
													"refactor.inline",
													"refactor.rewrite",
													"source",
													"source.organizeImports"
												]
											}
										},
										"honorsChangeAnnotations": false
									},
									"codeLens": {
										"dynamicRegistration": true
									},
									"formatting": {
										"dynamicRegistration": true
									},
									"rangeFormatting": {
										"dynamicRegistration": true
									},
									"onTypeFormatting": {
										"dynamicRegistration": true
									},
									"rename": {
										"dynamicRegistration": true,
										"prepareSupport": true,
										"prepareSupportDefaultBehavior": 1,
										"honorsChangeAnnotations": true
									},
									"documentLink": {
										"dynamicRegistration": true,
										"tooltipSupport": true
									},
									"typeDefinition": {
										"dynamicRegistration": true,
										"linkSupport": true
									},
									"implementation": {
										"dynamicRegistration": true,
										"linkSupport": true
									},
									"colorProvider": {
										"dynamicRegistration": true
									},
									"foldingRange": {
										"dynamicRegistration": true,
										"rangeLimit": 5000,
										"lineFoldingOnly": true
									},
									"declaration": {
										"dynamicRegistration": true,
										"linkSupport": true
									},
									"selectionRange": {
										"dynamicRegistration": true
									},
									"callHierarchy": {
										"dynamicRegistration": true
									},
									"semanticTokens": {
										"dynamicRegistration": true,
										"tokenTypes": [
											"namespace",
											"type",
											"class",
											"enum",
											"interface",
											"struct",
											"typeParameter",
											"parameter",
											"variable",
											"property",
											"enumMember",
											"event",
											"function",
											"method",
											"macro",
											"keyword",
											"modifier",
											"comment",
											"string",
											"number",
											"regexp",
											"operator"
										],
										"tokenModifiers": [
											"declaration",
											"definition",
											"readonly",
											"static",
											"deprecated",
											"abstract",
											"async",
											"modification",
											"documentation",
											"defaultLibrary"
										],
										"formats": [
											"relative"
										],
										"requests": {
											"range": true,
											"full": {
												"delta": true
											}
										},
										"multilineTokenSupport": false,
										"overlappingTokenSupport": false
									},
									"linkedEditingRange": {
										"dynamicRegistration": true
									}
								},
								"window": {
									"showMessage": {
										"messageActionItem": {
											"additionalPropertiesSupport": true
										}
									},
									"showDocument": {
										"support": true
									},
									"workDoneProgress": true
								},
								"general": {
									"regularExpressions": {
										"engine": "ECMAScript",
										"version": "ES2020"
									},
									"markdown": {
										"parser": "marked",
										"version": "1.1.0"
									}
								}
							},
							"trace": "off",
							"workspaceFolders": [
								{
									"uri": "file:///home/hank/CodingWork/Go/github.com/github-actions/testplaintext",
									"name": "testplaintext"
								}
							]
						}
					}`,
			want: &lspRequest{
				Header: &lspHeader{
					ContentLength: 23,
					ContentType:   "application/vscode-jsonrpc; charset=utf-8",
				},
				Body: &lspBody{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := parseRequest(strings.NewReader(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnMarshal(t *testing.T) {
	tests := []struct {
		name string
		input string
		want string
	}{
		{	
			name: "testUnmarshal",
			input: string(`{"operation": "get", "key": "example"}`),
			want: "Jsonrpc: gext",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := unmarshal()
			// require.Equal(t, tt.want, got)
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *lspHeader
	}{
		{
			name: "base case",
			input: "Content-Length: 23\r\n" +
				"Content-Type: json\r\n" +
				"\r\n" +
				`{}`,
			want: &lspHeader{
				ContentLength: 23,
				ContentType:   "json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeader(strings.NewReader(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_splitOnce(t *testing.T) {
	tests := []struct {
		input      string
		sep        string
		wantPrefix string
		wantSuffix string
	}{
		{
			input:      "Content-Length: hello",
			sep:        ": ",
			wantPrefix: "Content-Length",
			wantSuffix: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotPrefix, gotSuffix, err := splitOnce(tt.input, tt.sep)
			require.NoError(t, err)
			require.Equal(t, tt.wantPrefix, gotPrefix)
			require.Equal(t, tt.wantSuffix, gotSuffix)
		})
	}
}
