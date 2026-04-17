package main

import (
	"net/http"
	"crypto/rand"
	b32 "encoding/base32"
	"net/url"
	"io"
	"fmt"
)

var lhost = ""
var lport = ""

func generateSessionId() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return b32.StdEncoding.EncodeToString(bytes)[:32], nil
}

func main() {
	serverAddr := lhost + ":" + lport
	
	sessionID, err := generateSessionId()
	if err != nil {
		return
	}

	resp, err := http.PostForm("http://" + serverAddr + "/createSession", url.Values{"sid": {sessionID}})
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", body)
}