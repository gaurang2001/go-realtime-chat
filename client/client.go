package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
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
			fmt.Printf("\r%s\n", msg)
			fmt.Printf(">>")
		case <-exitChan:
			finalTermChan <- true
			return
		}
	}
}

func (cli *client) getClientMessage(sc bufio.Scanner, recMsg chan string) {
	fMsg := ""
	var ch int
	var m string
	fmt.Printf("Choose among the following (enter 1/2/3):\n\t1. PM\n\t2. Broadcast\n\t3. Terminate\n\t >>")
	fmt.Scan(&ch)
	switch ch {
	case 1:
		fmt.Printf("\t Name of user to send message:\n>>")
		fmt.Scan(&m)
		fMsg = "1~" + m
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		fMsg = fMsg + "~" + m
	case 2:
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		fMsg = "0~" + "~" + cli.Username + "~" + m
	case 3:
		fmt.Printf("\t Reason:\n\t >>")
		fmt.Scan(&m)
		fMsg = "2~" + "~" + m + "~"
	default:
		fmt.Printf("Choose right option\n\n")
		fMsg = "esc"
	}
	recMsg <- fMsg
}

func (cli *client) listenForClientMessages(ctx context.Context, sc bufio.Scanner, conn net.Conn, finalTermChan chan bool, termChan chan bool) {
	rMsg := make(chan string)
	for {
		go cli.getClientMessage(sc, rMsg)
		select {
		case <-ctx.Done():
			finalTermChan <- true
			return
		case mess := <-rMsg:
			msg := strings.Trim(mess, "\r\n")
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
	sc := bufio.Scanner{}
	termChanClient := make(chan bool)
	termChanServer := make(chan bool)
	termChanExit := make(chan bool)
	m := "3~" + cli.ServerPassword + "~" + cli.Username + "~"
	conn.Write([]byte(shared.Padd(m)))
	go cli.listenForClientMessages(ctx, sc, conn, termChanClient, termChanExit)
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
