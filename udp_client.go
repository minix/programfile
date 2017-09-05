package main

import (
	"fmt"
	"net"
)

func main() {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 23456,
	})

	if err != nil {
		fmt.Println("connect Faild", err)
		return
	}

	defer socket.Close()

	senddata := []byte("hello Minix")

	_, err = socket.Write(senddata)

	if err != nil {
		fmt.Println("send data Faild!", err)
		return
	}

	data := make([]byte, 4096)
	read, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("read data Faild!", err)
		return
	}

	fmt.Println(read, remoteAddr)
	fmt.Printf("%s\n", data)
}
