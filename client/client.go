package client

import (
	"bufio"
	"context"
	"fmt"
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

	
	if len(username) == 0 {
		username = shared.RandSeq(8)
	}
	return &client{username, password, host}
}

func getServerMessage(conn net.Conn, rcvd_msg chan string, exit_chan chan bool) {

	
	finalmessage := make([]byte, 256)
	_, err := io.ReadFull(conn, finalmessage)
	shared.CheckError(err)
	mess := string(finalmessage)
	msg := strings.Trim(mess, "\r\n")
	args := strings.Split(msg, "~")
	if err == io.EOF {
		fmt.Println("Server exited unexpectedly")
		exit_chan <- true
	} else if strings.Compare(args[0], "2") == 0 {
		fmt.Printf("Server - %s", args[1])
		exit_chan <- true
	} else {
		rcvd_msg <- mess
	}
}

func (cli *client) listenForServerMessages(ctx context.Context, conn net.Conn, final_term_chan chan bool) {

	
	recv_mess := make(chan string)
	exit_chan := make(chan bool)
	for {
		go getServerMessage(conn, recv_mess, exit_chan)
		select {
		case <-ctx.Done():
			final_term_chan <- true
			return
		case mess := <-recv_mess:
			msg := strings.Trim(mess, "\r\n")
			fmt.Printf("\r%s\n", msg)
			fmt.Printf(">>")
		case <-exit_chan:
			final_term_chan <- true
			return
		}
	}
}

func (cli *client) getClientMessage(sc bufio.Scanner, rcvd_msg chan string) {
	f_m := ""
	var ch int
	var m string
	fmt.Printf("Choose among the following (enter 1/2/3):\n\t1. PM\n\t2. Broadcast\n\t3. Terminate\n\t >>")
	fmt.Scan(&ch)
	switch ch {
	case 1:
		fmt.Printf("\t Name of user to send message:\n>>")
		fmt.Scan(&m)
		f_m = "1~" + m
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		f_m = f_m + "~" + m
	case 2:
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		f_m = "0~" + "~" + cli.Username + "~" + m
	case 3:
		fmt.Printf("\t Reason:\n\t >>")
		fmt.Scan(&m)
		f_m = "2~" + "~" + m + "~"
	default:
		fmt.Printf("Choose right option\n\n")
		f_m = "esc"
	}
	rcvd_msg <- f_m
}

func (cli *client) listenForClientMessages(ctx context.Context, sc bufio.Scanner, conn net.Conn, final_term_chan chan bool, term_chan chan bool) {
	recv_mess := make(chan string)
	for {
		go cli.getClientMessage(sc, recv_mess)
		select {
		case <-ctx.Done():
			final_term_chan <- true
			return
		case mess := <-recv_mess:
			msg := strings.Trim(mess, "\r\n")
			args := strings.Split(msg, "~")
			if strings.Compare(args[0], "esc") == 0 {
				continue
			}
			conn.Write([]byte(shared.Padd(msg)))
			if strings.Compare(args[0], "2") == 0 {
				final_term_chan <- true
				return
			}
		case <-term_chan:
			final_term_chan <- true
			return
		}
	}
}

func Client(password string, host string, username string) *client {

	
	if len(username) == 0 {
		username = shared.RandSeq(8)
	}
	return &client{username, password, host}
}

func getServerMessage(conn net.Conn, rcvd_msg chan string, exit_chan chan bool) {

	
	finalmessage := make([]byte, 256)
	_, err := io.ReadFull(conn, finalmessage)
	shared.CheckError(err)
	mess := string(finalmessage)
	msg := strings.Trim(mess, "\r\n")
	args := strings.Split(msg, "~")
	if err == io.EOF {
		fmt.Println("Server exited unexpectedly")
		exit_chan <- true
	} else if strings.Compare(args[0], "2") == 0 {
		fmt.Printf("Server - %s", args[1])
		exit_chan <- true
	} else {
		rcvd_msg <- mess
	}
}

func (cli *client) listenForServerMessages(ctx context.Context, conn net.Conn, final_term_chan chan bool) {

<<<<<<< HEAD
	
	recv_mess := make(chan string)
	exit_chan := make(chan bool)
	for {
		go getServerMessage(conn, recv_mess, exit_chan)
		select {
		case <-ctx.Done():
			final_term_chan <- true
			return
		case mess := <-recv_mess:
			msg := strings.Trim(mess, "\r\n")
			fmt.Printf("\r%s\n", msg)
			fmt.Printf(">>")
		case <-exit_chan:
			final_term_chan <- true
			return
		}
	}
}

func (cli *client) getClientMessage(sc bufio.Scanner, rcvd_msg chan string) {
	f_m := ""
	var ch int
	var m string
	fmt.Printf("Choose among the following (enter 1/2/3):\n\t1. PM\n\t2. Broadcast\n\t3. Terminate\n\t >>")
	fmt.Scan(&ch)
	switch ch {
	case 1:
		fmt.Printf("\t Name of user to send message:\n>>")
		fmt.Scan(&m)
		f_m = "1~" + m
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		f_m = f_m + "~" + m
	case 2:
		fmt.Printf("\t Message:\n\t >>")
		fmt.Scan(&m)
		f_m = "0~" + "~" + cli.Username + "~" + m
	case 3:
		fmt.Printf("\t Reason:\n\t >>")
		fmt.Scan(&m)
		f_m = "2~" + "~" + m + "~"
	default:
		fmt.Printf("Choose right option\n\n")
		f_m = "esc"
	}
	rcvd_msg <- f_m
}

func (cli *client) listenForClientMessages(ctx context.Context, sc bufio.Scanner, conn net.Conn, final_term_chan chan bool, term_chan chan bool) {
	recv_mess := make(chan string)
	for {
		go cli.getClientMessage(sc, recv_mess)
		select {
		case <-ctx.Done():
			final_term_chan <- true
			return
		case mess := <-recv_mess:
			msg := strings.Trim(mess, "\r\n")
			args := strings.Split(msg, "~")
			if strings.Compare(args[0], "esc") == 0 {
				continue
			}
			conn.Write([]byte(shared.Padd(msg)))
			if strings.Compare(args[0], "2") == 0 {
				final_term_chan <- true
				return
			}
		case <-term_chan:
			final_term_chan <- true
			return
		}
	}
}

func (cli *client) Run(ctx context.Context, main_term_chan chan bool) {

=======
>>>>>>> 6806767802173e00c287c905dd43ed512233d0c5
	conn, err := net.Dial("tcp", ":"+cli.ServerHost)
	shared.CheckError(err)
	sc := bufio.Scanner{}
	term_chan_client := make(chan bool)
	term_chan_server := make(chan bool)
	term_chan_exit := make(chan bool)
	m := "3~" + cli.ServerPassword + "~" + cli.Username + "~"
	conn.Write([]byte(shared.Padd(m)))
	go cli.listenForClientMessages(ctx, sc, conn, term_chan_client, term_chan_exit)
	go cli.listenForServerMessages(ctx, conn, term_chan_server)
	select {
	case <-ctx.Done():
		main_term_chan <- true
		break
	case <-term_chan_server:
		term_chan_exit <- true
		select {
		case <-term_chan_client:
			break
		}
		main_term_chan <- true
		break
	case <-term_chan_client:
		select {
		case <-term_chan_server:
			break
		}
		main_term_chan <- true
		break
	}
}
