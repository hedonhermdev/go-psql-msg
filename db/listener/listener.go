package listener

import (
	"log"
	"time"

	"github.com/lib/pq"
)

type JSONListener struct {
	listener *pq.Listener
}

func reportProblem(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("ERROR: %q", err)
	}
}

func NewJSONListener(conn_string string) *JSONListener {
	return &JSONListener{pq.NewListener(conn_string, 5*time.Second, 100*time.Second, reportProblem)}
}

func (l *JSONListener) Listen(pg_channel string, eventChan chan Event, errorChan chan error) error {
	err := l.listener.Listen(pg_channel)
	if err != nil {
		return err
	}
	for {
		select {
		case n := <-l.listener.Notify:
			event, err := ParseEvent(n)
			if err != nil {
				errorChan <- err
				continue
			}
			eventChan <- *event
		case <-time.After(90 * time.Second):
			log.Printf("INFO: No new notifications in 90 seconds, checking connection")
			err := l.listener.Ping()
			if err != nil {
				return err
			}
		}
	}
}
