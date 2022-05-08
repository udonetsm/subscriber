package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"subscriber/cache"
	"subscriber/db"
	"subscriber/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

var service_cache *cache.Cache

func init() {
	cache := cache.New()
	keys, values := db.GetAll()
	for i := 0; i < len(keys); i++ {
		cache.Set(keys[i], values[i])
	}
	service_cache = cache
}

func main() {
	ConnectAndSubscribe("_", "test-cluster", "nats://127.0.0.1:4222", "orders")
}

func ConnectAndSubscribe(clientid, clusterid, url, sub string) {
	fmt.Print("Connecting to nats streaming server --> ")
	time.Sleep(2 * time.Second) //quasi DDoS protecting :-D
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
	delta := db.GetDelta()
	subscriber, err := sc.QueueSubscribe(sub, "que", func(msg *stan.Msg) {
		order_id, err := db.Set(msg.Data, msg.Timestamp)
		if err == nil {
			service_cache.Set(order_id, string(msg.Data))
		}
	}, stan.StartAtTimeDelta(time.Duration(delta)))
	if err != nil {
		log.Println(err)
		return
	}
	HTTPServing()
	defer subscriber.Unsubscribe()
}

func HTTPServing() {
	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", GetById).Methods(http.MethodGet)
	http.ListenAndServe(":8585", router)
}

func GetById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	cashed_order, ok := service_cache.Get(id)
	order := models.Order{}
	if ok {
		err := json.Unmarshal([]byte(cashed_order), &order)
		if err == nil {
			w.Header().Add("Content-Type", "application/json")
			models.Show(&order, w, r) //using interface from models
			return
		}
	}
	w.Write([]byte("No data"))
}
