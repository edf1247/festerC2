package server

import (
	"strings"
	"strconv"
	"fmt"
	"os"
	"bufio"
)

func (l *Listeners) Interact(sid string) {

	splitString := strings.Split(sid, ":")

	listenerID, err := strconv.Atoi(splitString[len(splitString) - 1])
	if err != nil {
		fmt.Println(err)
		return
	}

	listener := l.listeners[listenerID]
	session := listener.activeSessions[sid]
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s> ", splitString[0])
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		cleanInput := strings.Split(input, " ")	
		
		if len(cleanInput) == 1 && cleanInput[0] == "" {
			continue
		} else if cleanInput[0] == "exit" {
			break
		} else {
			session.q = append(session.q, Command{input: cleanInput})
		}
	}
}