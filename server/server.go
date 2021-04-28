/*
* This file contains the definition for the server struct and the functions it implements.

There are two main classes of goroutines that exist in the server process:
- The first class is a single goroutine that listens for connection requests from clients.
  Whenever it receives a connection request, it creates a new goroutine that would come under
  the second class.
- The second class contains goroutines that handles the already established connections between
  the clients. These goroutines allow the server and the client to send messages to each other.
  Whenever a message is received from a particular client, the recipient's address is determined
  using the map containing the current connections. The message is then sent to the relevant client
  through the goroutine corresponding to that client.
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

	"github.com/gaurang2001/go-realtime-chat/shared"
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
				str := "2~" + "Goodbye!\n\n" + "~\n"
				conn.Write([]byte(shared.Padd(str)))
				return
			}
		}
	}
}

func (ser *server) handleClient(ctx context.Context, conn net.Conn, m *sync.RWMutex, wg *sync.WaitGroup) {

	/*
		Spawned by Run() when a client connection is received. Performs authentication and username checking, responds
		appropriately and after that, uses listenForMessages to handle incoming messages. Should handle cancellation of context
		and messages written to the term channel in the above function
	*/

	finalmessage := make([]byte, 256)
	if _, err := io.ReadFull(conn, finalmessage); err != nil {
		shared.CheckError(err)
		return
	}
	message := string(finalmessage)
	msg := strings.Trim(message, "\r\n")
	args := strings.Split(msg, "~")
	if strings.Compare(args[0], "3") == 0 {
		fmt.Printf("\nNew User has logged in!")
		if strings.Compare(args[1], ser.password) == 0 {
			fmt.Printf("\nThe password entered is correct")
			if _, found := ser.clients[args[2]]; found == false {
				wg.Add(1)
				m.Lock()
				ser.clients[args[2]] = conn
				fmt.Printf("%s has logged in\n", args[2])
				m.Unlock()
				wg.Done()
				conn.Write([]byte(shared.Padd("\nauthenticated\n\n")))
				term := make(chan bool)
				go ser.listenForMessages(ctx, conn, args[2], term, m, wg)
				select {
				case v := <-term:
					if v == true {
						wg.Add(1)
						m.Lock()
						fmt.Printf("%s has logged out\n", args[2])
						delete(ser.clients, args[2])
						m.Unlock()
						wg.Done()
						return
					}
				case <-ctx.Done():
					wg.Add(1)
					m.Lock()
					delete(ser.clients, args[2])
					m.Unlock()
					wg.Done()
					return
				}
			} else {
				fmt.Printf("User already exists! \n")
				conn.Write([]byte(shared.Padd("2~invalid_user\n")))
			}
		} else {
			fmt.Printf(" Server Password entered is wrong!")
			conn.Write([]byte(shared.Padd("2~invalid_password~\n")))
		}
	} else {
		fmt.Printf(" Request made is not for Logging In!")
		conn.Write([]byte(shared.Padd("2~invalid_request~\n")))
	}
}

func (ser *server) listenForConnections(ctx context.Context, newConn chan net.Conn, listener *net.TCPListener, m *sync.RWMutex, wg *sync.WaitGroup) {

	// Called from Run()
	// Accept incoming connections from clients and write it to the newConn channel
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				shared.CheckError(err)
				continue
			}
			newConn <- conn
		}
	}
}

func (ser *server) Run(ctx context.Context, done chan bool) {

	newConn := make(chan net.Conn)
	service := ":" + ser.address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	shared.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	shared.CheckError(err)
	fmt.Printf("\nServer listening on Port : %s \n\n", ser.address)
	defer listener.Close()
	var m sync.RWMutex
	wg := sync.WaitGroup{}
	go ser.listenForConnections(ctx, newConn, listener, &m, &wg)
	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-newConn:
			go ser.handleClient(ctx, conn, &m, &wg)
		}
	}
	done <- true
}
