package server

import (
	"os/exec"
	"fmt"
)

func CreateAgent(lhost string, lport string) {
	linkerString := fmt.Sprintf("-X 'main.lhost=%s' -X 'main.lport=%s'", lhost, lport)

	cmd := exec.Command("go", "build", "-ldflags", linkerString, "-o", "build/agent", "./cmd/client/main.go")

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}