"use strict";
/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode_1 = require("vscode");
const net = require("net");
const vscode_languageclient_1 = require("vscode-languageclient");
let client;
function activate(context) {
    console.log("starting client...");
    const executable = vscode_1.workspace.getConfiguration("plaintext").get;
    console.log(executable);
    context.subscriptions.push(startLanguageServerTCP(5007, ["plaintext"]));
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
    return new vscode_languageclient_1.LanguageClient(`tcp language server (port ${address})`, serverOptions, clientOptions).start();
}
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
exports.deactivate = deactivate;
// // The server is a started as a separate app and listens on port 5007
//     // const connectionInfo = {
//     //     port: 5007,
//     //     host: "127.0.0.1",
//     //     (connectListener) => {
//     //     }
//     // };
//     const serverOptions = async () => {
//         // Connect to language server via socket
//         const socket = net.connect({port: 5007}, () => { //'connect' listener
//           console.log('connected to server!');
//         });
//         const result: StreamInfo = {
//             writer: socket,
//             reader: socket
//         };
//         return await Promise.resolve(result);
//     };
// 	// Options to control the language client
// 	const clientOptions: LanguageClientOptions = {
// 		// Register the server for plain text documents
// 		documentSelector: [{ scheme: 'file', language: 'plaintext' }],
// 		synchronize: {
// 			// Notify the server about file changes to '.clientrc files contained in the workspace
// 			fileEvents: workspace.createFileSystemWatcher('**/*.*')
// 		},
// 		middleware: {
// 			didOpen: (document, next) => {
//                 console.log("didOpen", document);
//                 return next(document);
// 			},
// 			provideCompletionItem: async (document, position, context, token, next) => {
// 				// If not in `<style>`, do not perform request forwarding
// 				console.log(context);
// 				console.log(document);
// 				return await next(document, position, context, token);
// 			},
// 		}
// 	};
// 	// Create the language client and start the client.
// 	client = new LanguageClient(
// 		'languageServerExample',
// 		'Language Server Example',
// 		serverOptions,
// 		clientOptions
// 	);
// 	// Start the client. This will also launch the server
// 	client.start();
