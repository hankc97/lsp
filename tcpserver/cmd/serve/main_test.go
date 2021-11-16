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
			input: "Content-Length: 23\r\n" +
				"Content-Type: application/vscode-jsonrpc; charset=utf-8\r\n" +
				"\r\n" +
				`{"Name": "ED"}`,
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
