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
