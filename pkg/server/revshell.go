package server

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"os/signal"
)

func (session *Session) Shell() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("shell> ")
		inputChan := make(chan string, 1)

		go func() {
			input, _ := reader.ReadString('\n')
			inputChan <- strings.TrimSpace(input)
		}()

		select {
		case <-sigChan:
			return
		case input := <-inputChan:
			cleanInput := strings.Split(input, " ")
			if cleanInput[0] == "exit" {
				return
			}
			c := Command{input: input}
			if session.commChan != nil {
				session.respChan = make(chan string, 1)
				select {
				case <-sigChan:
					return
				case session.commChan <- c:
					//sent
				}
				select {
				case <-sigChan:
					return
				case resp := <-session.respChan:
					fmt.Printf("shell> %s", resp)
				}
			} else {
				if c.input != "" {
					session.commQueue = append(session.commQueue, c)
					fmt.Printf("\nCommand added to queue: %v\n", session.commQueue)
				}
			}	
		}
	}
}

func (session *Session) PrintQueue() {
	for len(session.commResponseQueue) > 0 {
		curr := session.commResponseQueue[0]
		session.commResponseQueue = session.commResponseQueue[1:]
		fmt.Printf("Input: %s | Output: %s\n", curr.input, curr.output)
	}
}
