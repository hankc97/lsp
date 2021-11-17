import * as net from 'net';

import {Trace} from 'vscode-jsonrpc';
import { window, workspace, commands, ExtensionContext, Uri } from 'vscode';
import { LanguageClient, LanguageClientOptions, StreamInfo, Position as LSPosition, Location as LSLocation } from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    console.log("starting client...");
    // The server is a started as a separate app and listens on port 5007
    const connectionInfo = {
        port: 5007
    };
    const serverOptions = () => {
        // Connect to language server via socket
        const socket = net.connect(connectionInfo);
        const result: StreamInfo = {
            writer: socket,
            reader: socket
        };
        return Promise.resolve(result);
    };
    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'plaintext' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/*.*')
        }
    };
    
    // Create the language client and start the client.
    const client = new LanguageClient('tcp server', serverOptions, clientOptions);

    // enable tracing (.Off, .Messages, Verbose)
    client.trace = Trace.Verbose;
    const disposable = client.start();
    
    // Push the disposable to the context's subscriptions so that the 
    // client can be deactivated on extension deactivation
    context.subscriptions.push(disposable);
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
      return undefined;
    }
    return client.stop();
}

// transport methods for the language server
// export enum TransportKind {
// 	stdio,
// 	ipc,
// 	pipe,
// 	socket
// }