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
const vscode_1 = require("vscode");
const net = require("net");
const vscode_languageclient_1 = require("vscode-languageclient");
let client;
function activate(context) {
    // The server is implemented in node
    // const serverModule = context.asAbsolutePath(
    // 	path.join('server', 'out', 'server.js')
    // );
    // The debug options for the server
    // --inspect=6009: runs the server in Node's Inspector mode so VS Code can attach to the server for debugging
    // const debugOptions = { execArgv: ['--nolazy', '--inspect=6009'] };
    // If the extension is launched in debug mode then the debug server options are used
    // Otherwise the run options are used
    console.log("starting client...");
    // The server is a started as a separate app and listens on port 5007
    const connectionInfo = {
        port: 5007
    };
    const serverOptions = () => {
        // Connect to language server via socket
        const socket = net.connect(connectionInfo);
        const result = {
            writer: socket,
            reader: socket
        };
        return Promise.resolve(result);
    };
    // const serverOptions: ServerOptions = {
    // 	run: { module: serverModule, transport: TransportKind.ipc },
    // 	debug: {
    // 		module: serverModule,
    // 		transport: TransportKind.ipc,
    // 		options: debugOptions
    // 	}
    // };
    // Options to control the language client
    const clientOptions = {
        // Register the server for plain text documents
        documentSelector: [{ scheme: 'file', language: 'plaintext' }],
        synchronize: {
            // Notify the server about file changes to '.clientrc files contained in the workspace
            fileEvents: vscode_1.workspace.createFileSystemWatcher('**/*.*')
        },
        middleware: {
            didOpen: (document, next) => {
                console.log("didOpen", document);
                return next(document);
            },
            provideCompletionItem: (document, position, context, token, next) => __awaiter(this, void 0, void 0, function* () {
                // If not in `<style>`, do not perform request forwarding
                console.log(context);
                console.log(document);
                return yield next(document, position, context, token);
            }),
        }
    };
    // Create the language client and start the client.
    client = new vscode_languageclient_1.LanguageClient('languageServerExample', 'Language Server Example', serverOptions, clientOptions);
    // Start the client. This will also launch the server
    client.start();
}
exports.activate = activate;
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
