package main

import "subscriber/controllers"

func main() {
	controllers.ConnectAndSubscribe("_", "test-cluster", "nats://127.0.0.1:4222", "orders")
}
