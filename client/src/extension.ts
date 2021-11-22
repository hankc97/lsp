/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import * as path from 'path';
import { workspace, ExtensionContext, TextDocument, Disposable } from 'vscode';
import * as net from 'net';
import {
	LanguageClient,
	LanguageClientOptions,
	ServerOptions,
	TransportKind,
    StreamInfo,
} from 'vscode-languageclient';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    console.log("starting client...");

    const executable = workspace.getConfiguration("plaintext").get;
    console.log(executable);

    context.subscriptions.push(startLanguageServerTCP(5007, ["plaintext"]));
}


function startLanguageServer(command: string, args: string[], documentSelector: string[]): Disposable {
    const serverOptions: ServerOptions = {
        command, 
        args,
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: documentSelector,
        synchronize: {
            configurationSection: "plaintext"
        },
    }

    return new LanguageClient(
        command, 
        serverOptions,
        clientOptions
    ).start()
} 

function startLanguageServerTCP(address: number, documentSelector: string[]): Disposable {
    const serverOptions: ServerOptions = () => {
        return new Promise((resolve, reject) => {
            const client = new net.Socket();
            client.connect(address, "127.0.0.1", () => {
                resolve({reader: client, writer: client});
            })
        })
    }

    const clientOptions: LanguageClientOptions = {
        documentSelector: documentSelector,
    }

    return new LanguageClient(`tcp language server (port ${address})`, serverOptions, clientOptions).start();
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}







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