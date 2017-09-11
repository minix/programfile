package main

import (
	//	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	socket, err := net.DialUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 23456})
	checkErr(err)

	defer socket.Close()

	data, err := ioutil.ReadAll(os.Stdin)

	_, err = socket.Write(data)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
