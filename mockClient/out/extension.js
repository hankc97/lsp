"use strict";
/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const net = require("net");
const vscode_languageclient_1 = require("vscode-languageclient");
let client;
const initResponseHeader = "Content-Length: 2987\r\n" +
    "\r\n";
const initResponseBody = `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"processId":1950,"clientInfo":{"name":"vscode","version":"1.62.3"},"rootPath":"/home/hank/CodingWork/Go/github.com/github-actions/testplaintext","rootUri":"file:///home/hank/CodingWork/Go/github.com/github-actions/testplaintext","capabilities":{"workspace":{"applyEdit":true,"workspaceEdit":{"documentChanges":true,"resourceOperations":["create","rename","delete"],"failureHandling":"textOnlyTransactional"},"didChangeConfiguration":{"dynamicRegistration":true},"didChangeWatchedFiles":{"dynamicRegistration":true},"symbol":{"dynamicRegistration":true,"symbolKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]}},"executeCommand":{"dynamicRegistration":true},"configuration":true,"workspaceFolders":true},"textDocument":{"publishDiagnostics":{"relatedInformation":true,"versionSupport":false,"tagSupport":{"valueSet":[1,2]}},"synchronization":{"dynamicRegistration":true,"willSave":true,"willSaveWaitUntil":true,"didSave":true},"completion":{"dynamicRegistration":true,"contextSupport":true,"completionItem":{"snippetSupport":true,"commitCharactersSupport":true,"documentationFormat":["markdown","plaintext"],"deprecatedSupport":true,"preselectSupport":true,"tagSupport":{"valueSet":[1]}},"completionItemKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]}},"hover":{"dynamicRegistration":true,"contentFormat":["markdown","plaintext"]},"signatureHelp":{"dynamicRegistration":true,"signatureInformation":{"documentationFormat":["markdown","plaintext"],"parameterInformation":{"labelOffsetSupport":true}},"contextSupport":true},"definition":{"dynamicRegistration":true,"linkSupport":true},"references":{"dynamicRegistration":true},"documentHighlight":{"dynamicRegistration":true},"documentSymbol":{"dynamicRegistration":true,"symbolKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]},"hierarchicalDocumentSymbolSupport":true},"codeAction":{"dynamicRegistration":true,"isPreferredSupport":true,"codeActionLiteralSupport":{"codeActionKind":{"valueSet":["","quickfix","refactor","refactor.extract","refactor.inline","refactor.rewrite","source","source.organizeImports"]}}},"codeLens":{"dynamicRegistration":true},"formatting":{"dynamicRegistration":true},"rangeFormatting":{"dynamicRegistration":true},"onTypeFormatting":{"dynamicRegistration":true},"rename":{"dynamicRegistration":true,"prepareSupport":true},"documentLink":{"dynamicRegistration":true,"tooltipSupport":true},"typeDefinition":{"dynamicRegistration":true,"linkSupport":true},"implementation":{"dynamicRegistration":true,"linkSupport":true},"colorProvider":{"dynamicRegistration":true},"foldingRange":{"dynamicRegistration":true,"rangeLimit":5000,"lineFoldingOnly":true},"declaration":{"dynamicRegistration":true,"linkSupport":true},"selectionRange":{"dynamicRegistration":true}},"window":{"workDoneProgress":true}},"trace":"off","workspaceFolders":[{"uri":"file:///home/hank/CodingWork/Go/github.com/github-actions/testplaintext","name":"testplaintext"}]}}`;
// const initResponseBody = 
function activate(context) {
    return __awaiter(this, void 0, void 0, function* () {
        console.log("configuring workspace...");
        const client = new net.Socket();
        yield new Promise((resolve, reject) => {
            client.connect(5007, "127.0.0.1", () => {
                resolve({ reader: client, writer: client });
            });
        });
        const bufBody = Buffer.from(initResponseBody);
        const initResponseHeader = `Content-Length: ${bufBody.length}\r\n` + "\r\n";
        const bufHeader = Buffer.from(initResponseHeader);
        const concatBuf = Buffer.concat([bufHeader, bufBody]);
        client.write(concatBuf);
        // const buf = Buffer.from(initResponseHeader)
        // client.write(buf)
        // const buf2 = Buffer.from(initResponseBody)
        // const jsonObj = JSON.parse(initResponseBody)
        // const jsonStr = JSON.stringify(jsonObj)
        // client.write(Buffer.from(jsonStr))
        yield client.read();
        // context.subscriptions.push(startLanguageServerTCP(5007, ["plaintext"]));
    });
}
exports.activate = activate;
function startLanguageServer(command, args, documentSelector) {
    const serverOptions = {
        command,
        args,
    };
    const clientOptions = {
        documentSelector: documentSelector,
        synchronize: {
            configurationSection: "plaintext"
        },
    };
    return new vscode_languageclient_1.LanguageClient(command, serverOptions, clientOptions).start();
}
function startLanguageServerTCP(address, documentSelector) {
    const serverOptions = () => {
        return new Promise((resolve, reject) => {
            const client = new net.Socket();
            client.connect(address, "127.0.0.1", () => {
                resolve({ reader: client, writer: client });
            });
        });
    };
    const clientOptions = {
        documentSelector: documentSelector,
    };
    const client = new vscode_languageclient_1.LanguageClient(`tcp language server (port ${address})`, serverOptions, clientOptions);
    const disposable = client.start();
    return disposable;
}
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
