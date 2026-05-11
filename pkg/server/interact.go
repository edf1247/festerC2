package server

import (
	"fmt"
	"strconv"
	"strings"
	"bufio"
	"os"
	"os/signal"
)

func (l *Listeners) Interact(sid string) {

	splitString := strings.Split(sid, ":")

	listenerID, err := strconv.Atoi(splitString[len(splitString)-1])
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, ok := l.listeners[listenerID]
	if !ok {
		fmt.Printf("Listener with id %d does not exist.\n", listenerID)
		return
	}

	session, ok := listener.activeSessions[sid]
	if !ok {
		fmt.Printf("Session with id %s does not exist.\n", splitString[0])
		return
	}

	prompt := splitString[0]
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	reader := bufio.NewReader(os.Stdin)

	for {
		for len(session.commResponseQueue) > 0 {
			curr := session.commResponseQueue[0]
			session.commResponseQueue = session.commResponseQueue[1:]
			fmt.Printf("Input: %s | Output: %s\n", curr.input, curr.output)
		}

		fmt.Printf("%s> ", prompt)

		inputChan := make(chan string, 1)

		go func() {
			input, _ := reader.ReadString('\n')
			inputChan <- strings.TrimSpace(input)
		}()
		
		select {
		case <-sigChan:
			fmt.Println("Quitting Session")
			return
		case input := <-inputChan:
			cleanInput := strings.Split(input, " ")

			if cleanInput[0] == "quit" {
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
					fmt.Printf("%s> %s", prompt, resp)
				}
			} else {
				if c.input != " " {
					session.commQueue = append(session.commQueue, c)
					fmt.Printf("Command added to queue: [%v]\n", session.commQueue)
				}
			}	
		}
	}
}
