package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type databaseNotification struct {
	Action string     `json:"action"`
	Table  string     `json:"table"`
	Data   notifyData `json:"data"`
}

type notifyData struct {
	Message Message `json:"message"`
}

func ListenMessageNotify(l *pq.Listener) {
	var (
		err        error
		notifyData databaseNotification
	)

	select {
	case notify := <-l.Notify:
		if err = json.Unmarshal([]byte(notify.Extra), &notifyData); err != nil {
			log.Printf("Error from ListenMessageNotification: %s", err.Error())
			return
		}

		notifyData.Data.Message.CreatedAt = notifyData.Data.Message.CreatedAt.UTC()

		b, _ := json.MarshalIndent(notifyData, "", "\t")
		fmt.Println("################################\n", string(b))

		return
	case <-time.After(2 * time.Minute):
		log.Println("Received no events from message_channel for 2 minutes, checking connection...")
		go func() {
			if err = l.Ping(); err != nil {
				log.Printf("Error from ping to message_channel: %s", err.Error())
				return
			} else {
				log.Println("Connection with message_channel is ok")
			}
		}()
		return
	}
}

func ListenDatabase(channel, db string, listenFunc func(*pq.Listener)) {
	var err error
	var listener *pq.Listener

	if func() error {
		if _, err = sql.Open("postgres", db); err != nil {
			return err
		}

		reportProblem := func(ev pq.ListenerEventType, err error) {
			if err != nil {
				log.Fatalln(err.Error())
			}
		}

		listener = pq.NewListener(db, 10*time.Second, time.Minute, reportProblem)
		if err = listener.Listen(channel); err != nil {
			return err
		}

		return nil
	}() != nil {
		log.Fatalf("Error from ListenDatabase: %s", err.Error())
	}

	log.Printf("Started monitoring database on %s channel...", channel)

	for {
		listenFunc(listener)
	}
}
