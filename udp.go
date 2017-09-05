package main

import (
	"fmt"
	"net"
)

func main() {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 23456,
	})

	if err != nil {
		fmt.Println("listen Faild", err)
		return
	}
	fmt.Println("listen success")
	defer socket.Close()

	for {
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("read UDP data Faild", err)
			continue
		}
		fmt.Println(read, remoteAddr)
		fmt.Printf("%s\n", data)

		senddata := []byte("hello client")
		_, err = socket.WriteToUDP(senddata, remoteAddr)
		if err != nil {
			return
			fmt.Println("send data Faild", err)
		}
	}
}
