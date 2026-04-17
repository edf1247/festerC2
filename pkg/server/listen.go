package server

import (
	"log"
	"net/http"
	"fmt"
	"strconv"
	"sync"
	"strings"
)

type Command struct {
	input []string
	output []string
}

type Session struct {
	id string
	q []Command
	rhost string
	listenerID int
}

type Listener struct {
	id 				int
	mu 				sync.Mutex
	server 			*http.Server
	activeSessions	map[string]*Session // map session id -> Session object
}

type Listeners struct {
	mu 		  sync.Mutex
	listeners map[int]*Listener
}

func (l *Listener) createSession(w http.ResponseWriter, r *http.Request) {
	sid := r.PostFormValue("sid") + ":" + strconv.Itoa(l.id)
	if sid == "" {
		fmt.Println("Malformed request")
		return
	}

	rh := r.RemoteAddr

	s := Session{id: sid, q: []Command{}, rhost: rh, listenerID: l.id}
	l.activeSessions[sid] = &s
	fmt.Printf("\nConnection recieved from %s. Session started with id %s. \n", rh, sid)
}

func (l *Listener) getCommands(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	session := l.activeSessions[sid]

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	commands := ""
	for i := range session.q {
		commands = commands + "," + strings.Join(session.q[i].input, " ")
	}

	w.Write([]byte(commands))
}

func InitListeners() (Listeners) {
	l := Listeners{
		listeners: make(map[int]*Listener),
	}
	return l
}	

func (l *Listeners) StartListener(lhost string, lport string) {
	listenerID := len(l.listeners)

	addr := lhost + ":" + lport

	go func() {
		mux := http.NewServeMux()
		
		s := &http.Server{
			Addr: addr,
			Handler: mux,
		}

		listener := Listener{server: s, activeSessions: make(map[string]*Session), id: listenerID}

		mux.HandleFunc("/createSession", listener.createSession)
		mux.HandleFunc("/getCommands/{id}", listener.getCommands)

		l.mu.Lock()
		l.listeners[listenerID] = &listener
		l.mu.Unlock()

		listener.server.ListenAndServe()
	}()

	fmt.Printf("Listener %d started on %s\n", listenerID, addr)
}

func (l *Listeners) KillListener(listenerID string) {
	intID, err := strconv.Atoi(listenerID)
	if err != nil {
		log.Println(err)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	listener, ok := l.listeners[intID]

	if !ok {
		fmt.Printf("Listener with id %d does not exist", intID)
		return
	}

	listener.server.Close()

	delete(l.listeners, intID)
}