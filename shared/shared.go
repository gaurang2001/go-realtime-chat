package shared

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandSeq(n int) string {

	// Generates a random sequence of characters of fixed length.
	// Used to generate a username when the user doesn't enter one

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CheckError(err error) {

	// Generic error checking. This function should be called
	// to check and print errors, and exit the program if an error
	// has occured.

	if err == io.EOF {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func Padd(str string) string {
	for len(str) < 256 {
		str = "\r" + str
	}
	return str
}
