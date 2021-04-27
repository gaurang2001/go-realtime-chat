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