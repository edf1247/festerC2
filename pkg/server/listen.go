package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Command struct {
	input  string
	output string
}

type Session struct {
	id         		  string
	commChan   		  chan Command
	respChan   		  chan string
	rhost      		  string
	listenerID 		  int
	commQueue  		  []Command
	commResponseQueue []Command
}

type Listener struct {
	id             int
	mu             sync.Mutex
	server         *http.Server
	activeSessions map[string]*Session // map session id -> Session object
}

type Listeners struct {
	mu        sync.Mutex
	listeners map[int]*Listener
}

func (l *Listener) createSession(w http.ResponseWriter, r *http.Request) {
	sid := r.PostFormValue("sid") + ":" + strconv.Itoa(l.id)
	sessionType := r.PostFormValue("type")
	if sid == "" {
		fmt.Println("Malformed request")
		return
	}

	if sessionType == "" {
		fmt.Println("Malformed request")
	}

	rh := r.RemoteAddr

	s := Session{id: sid, rhost: rh, listenerID: l.id}

	if sessionType == "1" {
		s.commChan = make(chan Command)
	} else {
		s.commQueue = []Command{}
	}

	l.activeSessions[sid] = &s
	fmt.Printf("\nConnection recieved from %s. Session started with id %s. \n", rh, sid)
}

func InitListeners() Listeners {
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
			Addr:    addr,
			Handler: mux,
		}

		listener := Listener{server: s, activeSessions: make(map[string]*Session), id: listenerID}

		mux.HandleFunc("/createSession", listener.createSession)

		mux.HandleFunc("/beacon/getCommands/{id}", listener.GetCommands)
		mux.HandleFunc("/beacon/response/{id}", listener.BeaconResponse)

		mux.HandleFunc("/agent/input/{id}", listener.AgentSSE)
		mux.HandleFunc("/agent/response/{id}", listener.AgentResponse)

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
