import * as net from 'net';

import {Trace} from 'vscode-jsonrpc';
import { window, workspace, commands, ExtensionContext, Uri } from 'vscode';
import { LanguageClient, LanguageClientOptions, StreamInfo, Position as LSPosition, Location as LSLocation } from 'vscode-languageclient/node';

export function activate(context: ExtensionContext) {
    console.log("starting client...");
    // The server is a started as a separate app and listens on port 5007
    const connectionInfo = {
        port: 5007
    };
    const serverOptions = async () => {
        // Connect to language server via socket
        const socket = net.connect(connectionInfo);
        const result: StreamInfo = {
            writer: socket,
            reader: socket
        };
        return Promise.resolve(result);
    };
    const clientOptions: LanguageClientOptions = {
        documentSelector: [{
            scheme: 'file',
            language: 'plaintext',
        }],
        // synchronize: {
        //     fileEvents: workspace.createFileSystemWatcher('**/*.*')
        // }
    };
    
    // Create the language client and start the client.
    const lc = new LanguageClient('golang-tcpserver', serverOptions, clientOptions);
    context.subscriptions.push(lc.start());

    // const disposable2 = commands.registerCommand("mydsl.a.proxy", async () => {
    //     const activeEditor = window.activeTextEditor;
    //     if (!activeEditor || !activeEditor.document || activeEditor.document.languageId !== 'mydsl') {
    //         return;
    //     }

    //     if (activeEditor.document.uri instanceof Uri) {
    //         commands.executeCommand("mydsl.a", activeEditor.document.uri.toString());
    //     }
    // });

    // context.subscriptions.push(disposable2);

    // enable tracing (.Off, .Messages, Verbose)
    // lc.trace = Trace.Verbose;
    
    // Push the disposable lc.start() to the context's subscriptions so that the 
    // client can be deactivated on extension deactivation
}