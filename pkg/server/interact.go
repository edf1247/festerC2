package server

import {
	"strings"
}

func (l *Listeners) Interact(sid string) {

	splitString := strings.Split(sid, ":")
	listenerID := splitString[-1]

	listener := l.listeners[listenerID]
	session := listener.activeSessions[sid]
}