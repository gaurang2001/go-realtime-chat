package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/gaurang2001/go-realtime-chat/shared"
)

type client struct {
	Username       string
	ServerPassword string
	ServerHost     string
}

func Client(password string, host string, username string) *client {

	/*
		An instance of the 'client' struct is created, initialized with given
		or the default data(if the user hasn't specified the data).
	*/

	if len(username) == 0 {
		username = shared.RandSeq(8)
	}
	return &client{username, password, host}
}

func getServerMessage(conn net.Conn, recMsg chan string, exitChan chan bool) {
	finalmessage := make([]byte, 256)
	_, err := io.ReadFull(conn, finalmessage)
	shared.CheckError(err)
	mess := string(finalmessage)
	msg := strings.Trim(mess, "\r\n")
	args := strings.Split(msg, "~")
	if err == io.EOF {
		fmt.Println("Server exited unexpectedly")
		exitChan <- true
	} else if strings.Compare(args[0], "2") == 0 {
		fmt.Printf("Server - %s", args[1])
		exitChan <- true
	} else {
		recMsg <- mess
	}
}

func (cli *client) listenForServerMessages(ctx context.Context, conn net.Conn, finalTermChan chan bool) {

	rMsg := make(chan string)
	exitChan := make(chan bool)
	for {
		go getServerMessage(conn, rMsg, exitChan)
		select {
		case <-ctx.Done():
			finalTermChan <- true
			return
		case mess := <-rMsg:
			msg := strings.Trim(mess, "\r\n")
			args := strings.Split(msg, "~")
			switch args[0] {
			case "0":
				fmt.Printf("\r%s[broadcast]: %s\n", args[2], args[3])
				break
			case "1":
				fmt.Printf("\r%s[private]: %s\n", args[1], args[2])
				break
			case "authenticated":
				fmt.Println("\rLogged in successfully!!")
			}
			fmt.Printf("\r>> ")
		case <-exitChan:
			finalTermChan <- true
			return
		}
	}
}

func (cli *client) getClientMessage(recMsg chan string) {
	fMsg := ""
	var ch int
	var m string
	fmt.Printf("\rChoose among the following (enter 1/2/3):\n\t1. PM\n\t2. Broadcast\n\t3. Terminate\n\r>> ")
	fmt.Scan(&ch)
	in := bufio.NewReader(os.Stdin)

	switch ch {
	case 1:
		fmt.Printf("\rName of user to send message:\n\r>> ")
		fmt.Scan(&m)
		fMsg = "1~" + m
		fmt.Printf("\r\nMessage:\n>> ")
		line, err := in.ReadString('\n')
		shared.CheckError(err)
		fMsg = fMsg + "~" + line
	case 2:
		fmt.Printf("\r\nMessage:\n>> ")
		line, err := in.ReadString('\n')
		shared.CheckError(err)
		fMsg = "0~" + "~" + cli.Username + "~" + line
	case 3:
		fmt.Printf("\r\nReason:\n>> ")
		line, err := in.ReadString('\n')
		shared.CheckError(err)
		fMsg = "2~" + "~" + line + "~"
	default:
		fmt.Printf("\nChoose right option: ")
		fMsg = "esc"
	}
	recMsg <- fMsg
}

func (cli *client) listenForClientMessages(ctx context.Context, conn net.Conn, finalTermChan chan bool, termChan chan bool) {
	rMsg := make(chan string)
	for {
		go cli.getClientMessage(rMsg)
		select {
		case <-ctx.Done():
			finalTermChan <- true
			return
		case mess := <-rMsg:
			msg := strings.Trim(mess, "\n")
			args := strings.Split(msg, "~")
			if strings.Compare(args[0], "esc") == 0 {
				continue
			}
			conn.Write([]byte(shared.Padd(msg)))
			if strings.Compare(args[0], "2") == 0 {
				finalTermChan <- true
				return
			}
		case <-termChan:
			finalTermChan <- true
			return
		}
	}
}

func (cli *client) Run(ctx context.Context, mainTermChan chan bool) {

	conn, err := net.Dial("tcp", ":"+cli.ServerHost)
	shared.CheckError(err)
	termChanClient := make(chan bool)
	termChanServer := make(chan bool)
	termChanExit := make(chan bool)
	m := "3~" + cli.ServerPassword + "~" + cli.Username + "~"
	conn.Write([]byte(shared.Padd(m)))
	go cli.listenForClientMessages(ctx, conn, termChanClient, termChanExit)
	go cli.listenForServerMessages(ctx, conn, termChanServer)
	select {
	case <-ctx.Done():
		mainTermChan <- true
		break
	case <-termChanServer:
		termChanExit <- true
		select {
		case <-termChanClient:
			break
		}
		mainTermChan <- true
		break
	case <-termChanClient:
		select {
		case <-termChanServer:
			break
		}
		mainTermChan <- true
		break
	}
}
