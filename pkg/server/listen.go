package server

import (
	"net"
	"log"
	"net/http"
)

type Listeners struct {
	listeners map[int]net.Listener
}

func InitListeners() (Listeners) {
	l := Listeners{
		make(map[int]net.Listener),
	}
	return l
}	

func (l *Listeners) StartListener(lhost string, lport string) {
	listenerId := len(l.listeners)

	addr := lhost + ":" + lport

	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		l.listeners[listenerId] = listener
		http.Serve(listener, nil)
	}()

	log.Printf("Listener %d started on %s\n", listenerId, addr)
}