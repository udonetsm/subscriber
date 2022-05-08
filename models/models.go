package models

import (
	"encoding/json"
	"net/http"
	"time"
)

type Delivery struct {
	Name    string `json:"name ,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Zip     string `json:"zip,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
	Region  string `json:"region,omitempty"`
	Email   string `json:"email,omitempty"`
}

type Payment struct {
	Tarnsaction   string  `json:"tranaction,omitempty"`
	Request_Id    string  `json:"request_id,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	Provider      string  `json:"provider,omitempty"`
	Amount        int     `json:"amount,omitempty"`
	Payment_Dt    int64   `json:"payment_dt,omitempty"`
	Bank          string  `json:"bank,omitempty"`
	Delivery_cost float64 `json:"delivery_cost,omitempty"`
	Goods_total   float64 `json:"goods_total,omitempty"`
	Custom_fee    int     `json:"custom_fee,omitempty"`
}

type Items struct {
	Chrt_id      int64   `json:"chrt_id,omitempty"`
	Track_number string  `json:"track_number,omitempty"`
	Price        float64 `json:"price,omitempty"`
	Rid          string  `json:"rid,omitempty"`
	Name         string  `json:"name,omitempty"`
	Sale         float64 `json:"sale,omitempty"`
	Size         string  `json:"size,omitempty"`
	Total_price  float64 `json:"total_price,omitempty"`
	Nm_id        int64   `json:"nm_id,omitempty"`
	Brand        string  `json:"brand,omitempty"`
	Status       int     `json:"status,omitempty"`
}

type Order struct {
	Order_id           string    `json:"order_uid,omitempty"`
	Track_number       string    `json:"track_number,omitempty"`
	Entry              string    `json:"entry,omitempty"`
	Delivery           Delivery  `json:"delivery,omitempty"`
	Payment            Payment   `json:"payment,omitempty"`
	Items              []Items   `json:"items,omitempty"`
	Locale             string    `json:"locale,omitempty"`
	Internal_signature string    `json:"internal_signature,omitempty"`
	Customer_id        string    `json:"customer_id,omitempty"`
	Delivery_service   string    `json:"delivery_service,omitempty"`
	Shardkey           string    `json:"shardKey,omitempty"`
	Sm_id              int       `json:"sm_id,omitempty"`
	DateCreated        time.Time `json:"date_created,omitempty"`
	Oof_shard          string    `json:"oof_shard,omitempty"`
}

// Interface
type Render interface {
	ShowJSON(w http.ResponseWriter, r *http.Request)
}

// Object realize method
func (order *Order) ShowJSON(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&order)
}

// func which use interface
func Show(object Render, w http.ResponseWriter, r *http.Request) {
	object.ShowJSON(w, r)
}
