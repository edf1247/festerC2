package server

import (
	"log"
	"net/http"
	"fmt"
	"strconv"
	"sync"
)

type Command struct {
	input string
	output string
}

type Session struct {
	id string
	q []Command
	rhost string
}

type Listener struct {
	mu 				sync.Mutex
	server 			*http.Server
	activeSessions	map[string]*Session // map session id -> Session object
}

type Listeners struct {
	mu 		  sync.Mutex
	listeners map[int]*Listener
}
 
func (l *Listener) createSession(w http.ResponseWriter, r *http.Request) {
	sid := r.PostFormValue("sid")
	if sid == "" {
		fmt.Println("Malformed request")
		return
	}

	rh := r.RemoteAddr

	s := Session{id: sid, q: make([]Command, 0), rhost: rh}
	l.activeSessions[sid] = &s
	fmt.Printf("\nConnection recieved from %s. Session started with id %s. \n", rh, sid)
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
		s := &http.Server{
			Addr: addr,
		}

		listener := Listener{server: s, activeSessions: make(map[string]*Session)}
		http.HandleFunc("/createSession", listener.createSession)

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