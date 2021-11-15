package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
	// "bufio"
	// "lsp/tcpserver"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	listener, err := net.Listen("tcp", PORT)
	fmt.Printf("now listening on port%s...\n", PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	buffer := make([]byte, 0, 4096)
	contentLength := 0
	for {
		tmp := make([]byte, 256)
		_, err := conn.Read(tmp)
		// dataByte := bufio.NewReader(conn)
		// netData, err := dataByte.ReadString('\n')
		
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		fmt.Println("-----------------------------------")
		if strings.HasPrefix(strings.TrimSpace(string(tmp)), "Content") {
			contentLength = getContentlength(strings.TrimSpace(string(tmp)))
		} 
		fmt.Printf("%+v\n", strings.TrimSpace(string(tmp)))
		if !strings.HasPrefix(strings.TrimSpace(string(tmp)), "Content") {
			buffer = append(buffer, (tmp)...)
		}
		if len(buffer) >= contentLength {
			var request *jsonrpc2.Request
			fmt.Printf("%+v\n", string(buffer))
			err = json.Unmarshal(buffer, &request)

			if err != nil {
				// return fmt.Println(err)
			}
			fmt.Printf("%+v\n", request)
		}
	}

	// fmt.Println("--------buffer-----")

	// fmt.Printf("%+v\n", string(buffer))

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func getContentlength(s string) int {
	i, _ := strconv.Atoi(strings.TrimPrefix(s, "Content-Length: "))
	return i
}


// fmt.Println("--------------------")

// if err != nil {
// 	fmt.Printf("%+v\n", netData)
// 	tcpserver.To(netData)
// 	log.Fatal(err)
// }

// fmt.Printf("%+v\n", string(netData))


// if !strings.HasPrefix(string(netData), "Content") {
// 	err = json.Unmarshal(dataByte, &request)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// fmt.Printf("%+v\n", request)


// if strings.TrimSpace(string(netData)) == "STOP" {
// 	fmt.Println("Exiting TCP server!")
// 	return
// }

// fmt.Print("-> ", string(netData))
// t := time.Now()
// myTime := t.Format(time.RFC3339) + "\n"
// conn.Write([]byte(myTime))