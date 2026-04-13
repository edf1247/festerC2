package main

import (
	"fmt"
)

var lhost = "127.0.0.1"
var lport = "8080"

func main() {
	lAddr := lhost + ":" + lport
	fmt.Println(lAddr)
}
