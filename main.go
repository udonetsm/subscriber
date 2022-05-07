package main

import (
	"fmt"
	"log"
	"subscriber/db"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	ConnectAndSubscribe("_", "test-cluster", "nats://127.0.0.1:4222", "orders")
}

func ConnectAndSubscribe(clientid, clusterid, url, sub string) {
	fmt.Print("Connecting to nats streaming server --> ")
	sc, err := stan.Connect(clusterid, clientid, stan.NatsURL(url),
		stan.SetConnectionLostHandler(func(c stan.Conn, err error) {
			log.Println("Disconnected")
		}))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("[SUCCESS]")
	Subscribe(sub, sc)
}

func Subscribe(sub string, sc stan.Conn) {
	lastTimeWrittenData := db.GetMaxValue()
	delta := time.Now().UnixNano() - lastTimeWrittenData
	_, err := sc.QueueSubscribe(sub, "que", func(msg *stan.Msg) {
		db.New(msg.Data, msg.Timestamp)
	}, stan.StartAtTimeDelta(time.Duration(delta)))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Scanln()
}
