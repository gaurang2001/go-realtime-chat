/*
* This file contains the definition for the server struct and the functions it implements.
 */

package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gaurang2001/simple-go-service/shared"
)

type server struct {
	clients  map[string]net.Conn
	password string
	address  string
}

func Server(pass string, address string) *server {

	/*
		An instance of the 'server' struct is created, initialized with given
		or the default data(if the user hasn't specified the data).
	*/
	if len(pass) == 0 {
		pass = "1234"
	}

	if len(address) == 0 {
		address = "8080"
	}

	return &server{
		clients:  make(map[string]net.Conn),
		password: pass,
		address:  address,
	}
}

func (ser *server) listenForMessages(ctx context.Context, conn net.Conn, username string, term chan bool, m *sync.RWMutex, wg *sync.WaitGroup) {

	/**
	Method parameter description:
	1. ctx - cancellable context
	2. conn - represents the socket connection to the client
	2. username - client username
	3. term - write to this channel on receiving termination request from client
	*/

	/**
	Spawned by handleClient(). Listens for messages sent by the given client, and appropriately unicasts/broadcast/
	prints error message, etc.
	*/
	timeoutDuration := 300 * time.Second
	for {
		select {
		case <-ctx.Done():
			str := "2~" + "Context" + "~\n"
			conn.Write([]byte(shared.Padd(str)))
			return
		default:
			finalmessage := make([]byte, 256)
			conn.SetReadDeadline(time.Now().Add(timeoutDuration))
			if _, err := io.ReadFull(conn, finalmessage); err != nil {
				shared.CheckError(err)
				str := "2~" + err.Error() + "~\n"
				conn.Write([]byte(shared.Padd(str)))
				term <- true
				return
			}
			message := string(finalmessage)
			msg := strings.Trim(message, "\r\n")
			args := strings.Split(msg, "~")
			switch args[0] {
			case "0":
				wg.Add(1)
				m.RLock()
				for i, cli := range ser.clients {
					fmt.Println(i)
					if strings.Compare(i, username) != 0 {
						if _, err1 := cli.Write([]byte(message)); err1 != nil {
							str := "Unable to broadcast message to" + string(i)
							fmt.Printf("Unable to broadcast message to %s from %s", i, username)
							conn.Write([]byte(shared.Padd(str)))
						}
					}
				}
				m.RUnlock()
				wg.Done()
			case "1":
				sm := args[0] + "~" + username + "~" + args[2]
				wg.Add(1)
				m.RLock()
				v := 0
				for i, cli := range ser.clients {
					if strings.Compare(i, args[1]) == 0 {
						v = 1
						if _, err1 := cli.Write([]byte(shared.Padd(sm))); err1 != nil {
							str := "Unable to send message to" + string(i)
							fmt.Printf("Unable to send message to %s from %s", i, username)
							conn.Write([]byte(shared.Padd(str)))
							v = 0
						} else {
							str := "Message sent"
							conn.Write([]byte(shared.Padd(str)))
						}
						break
					}
				}
				m.RUnlock()
				wg.Done()
				if v == 0 {
					str := "User not found"
					conn.Write([]byte(shared.Padd(str)))
				}
			case "2":
				term <- true
				str := "2~" + "Goodbye!" + "~\n"
				conn.Write([]byte(shared.Padd(str)))
				return
			}
		}
	}
}
