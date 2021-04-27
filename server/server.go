/*
* This file contains the definition for the server struct and the functions it implements.
 */

package server

import (
	"net"
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
