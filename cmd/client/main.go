package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/decoesp/tamborete/internal/resp"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	respParser := resp.NewParser(bufio.NewReader(conn))

	for {
		fmt.Print("tamborete> ")
		cmd, _ := reader.ReadString('\n')
		if strings.TrimSpace(cmd) == "exit" {
			break
		}

		fmt.Fprintf(conn, cmd)
		response, err := respParser.Parse()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(response)
	}
}
