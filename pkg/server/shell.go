package server

import (
	"github.com/common-nighthawk/go-figure"
	"bufio"
	"os"
	"fmt"
	"strings"
)

func Splash() {
	myFigure := figure.NewColorFigure("Fester", "colossal", "white", true)
  	myFigure.Print()
}

func HandleInput(l Listeners) {
	
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
				fmt.Println("Usage: listen <listen ip> <listen port>")
				continue
			}
			l.StartListener(cleanInput[1], cleanInput[2])
		case "ls":
			fmt.Printf("ListenerID  Addr\n")
			for k, v := range l.listeners {
				fmt.Printf("%d \t  %s\n", k, v.Addr())
			}
		}
	}
}

func Init() {
	Splash()
	l := InitListeners()
	HandleInput(l)
}