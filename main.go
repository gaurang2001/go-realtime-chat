package main

import (
	"context"
	"fmt"

	"github.com/gaurang2001/go-realtime-chat/client"
	"github.com/gaurang2001/go-realtime-chat/server"
)

func main() {
	var arguments string
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("To start the server press 1 or to use the client press 2")
	fmt.Scan(&arguments)
	if arguments == "1" {
		defer cancel()
		var pass, addr string
		fmt.Printf("Enter the Server Password (leave empty for no password) : ")
		fmt.Scanln(&pass)
		fmt.Printf("Enter the Server Port (leave empty for default port) : ")
		fmt.Scanln(&addr)
		s := server.Server(pass, addr)
		done := make(chan bool)
		go s.Run(ctx, done)
		select {
		case <-done:
		}
	} else if arguments == "2" {
		var pass, name, addr string
		fmt.Printf("Enter the Client username (leave empty for default username) : ")
		fmt.Scanln(&name)
		fmt.Printf("Enter the Server Password : ")
		fmt.Scanln(&pass)
		fmt.Printf("Enter the Server Port you want the Client to attach to : ")
		fmt.Scan(&addr)
		term := make(chan bool)
		cli := client.Client(pass, addr, name)
		go cli.Run(ctx, term)
		select {
		case <-ctx.Done():
			break
		case <-term:
			break
		}
	}
}
