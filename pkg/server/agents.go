package server

import (
	"os/exec"
	"fmt"
)

func CreateAgent(lhost string, lport string, platform string) {
	linkerString := fmt.Sprintf("-X 'main.lhost=%s' -X 'main.lport=%s'", lhost, lport)

	cmd := exec.Command("ls")

	if platform == "linux" {
		cmd = exec.Command("go", "build", "-ldflags", linkerString, "-o", "build/agent", "./cmd/client/linux/agents/main.go")
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}