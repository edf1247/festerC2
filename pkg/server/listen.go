package server

import (
	"log"
	"net/http"
	"fmt"
	"strconv"
	"sync"
)

type Listeners struct {
	mu 		  sync.Mutex
	listeners map[int]*http.Server
}

func InitListeners() (Listeners) {
	l := Listeners{
		listeners: make(map[int]*http.Server),
	}
	return l
}	

func (l *Listeners) StartListener(lhost string, lport string) {
	listenerID := len(l.listeners)

	addr := lhost + ":" + lport

	go func() {
		listener := &http.Server{
			Addr: addr,
		}
		
		l.mu.Lock()
		l.listeners[listenerID] = listener
		l.mu.Unlock()

		listener.ListenAndServe()
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

	listener.Close()

	delete(l.listeners, intID)
}