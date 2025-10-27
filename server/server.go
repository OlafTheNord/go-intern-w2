package server

import (
	"fmt"
	"io"
	"net"
	"os"
)

func StartServ() {
	conn, err := net.Dial("tcp4", "192.168.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("server start")
	}
	defer conn.Close()

	io.Copy(os.Stdout, conn)
}
