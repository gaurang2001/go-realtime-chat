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
	fmt.Scan(&arguments)
	if arguments == "1" {
		defer cancel()
		var v int
		var pass, addr string
		fmt.Printf("Wanna enter Server Password:?(0/1) ")
		fmt.Scan(&v)
		if v == 1 {
			fmt.Printf("Enter the server Password:\t")
			fmt.Scan(&pass)
		}
		fmt.Printf("Wanna enter Port:?(0/1) ")
		fmt.Scan(&v)
		if v == 1 {
			fmt.Printf("Enter the Port you want the server to attach to: ")
			fmt.Scan(&addr)
		}
		s := server.Server(pass, addr)
		done := make(chan bool)
		go s.Run(ctx, done)
		select {
		case <-done:
		}
	} else if arguments == "2" {
		var pass, name, addr string
		var v int
		fmt.Printf("Wanna enter Username:?(0/1) ")
		fmt.Scan(&v)
		if v == 1 {
			fmt.Printf("Enter username: ")
			fmt.Scan(&name)
		}
		fmt.Printf("Enter the server Password: ")
		fmt.Scan(&pass)
		fmt.Printf("Enter the Port you want the server to attach to: ")
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
