package main

import (
	"bufio"
	"fmt"
	"net"
)

var str = `1 2 3 4`

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(conn, str)
	result, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
