package server

import (
	"fmt"
	"strconv"
	"strings"
	"bufio"
	"os"
)

func (l *Listeners) Interact(sid string) {
	commands := []string{"quit", "exit", "shell", "queue"}

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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s> ", prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		cleanInput := strings.Split(input, " ")

		switch cleanInput[0] {
		case "quit":
			return
		case "exit":
			return
		case "shell":
			fmt.Println("Dropping into shell...")
			session.Shell()
		case "queue":
			session.PrintQueue()
		case "help":
			for k := range commands {
				fmt.Println(commands[k])
			}
		}
	}
}
