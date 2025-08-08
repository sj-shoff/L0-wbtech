package model

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid"  db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" db:"shardkey"`
	SmID              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"    db:"name"`
	Phone   string `json:"phone"   db:"phone"`
	Zip     string `json:"zip"     db:"zip"`
	City    string `json:"city"    db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region"  db:"region"`
	Email   string `json:"email"   db:"email"`
}

type Payment struct {
	ID           int    `json:"id"          db:"id"`
	Transaction  string `json:"transaction" db:"transaction"`
	RequestID    string `json:"request_id"  db:"request_id"`
	Currency     string `json:"currency"    db:"currency"`
	Provider     string `json:"provider"    db:"provider"`
	Amount       int    `json:"amount"      db:"amount"`
	PaymentDt    int64  `json:"payment_dt"  db:"payment_dt"`
	Bank         string `json:"bank"        db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"  db:"goods_total"`
	CustomFee    int    `json:"custom_fee"  db:"custom_fee"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id"      db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price"       db:"price"`
	Rid         string `json:"rid"         db:"rid"`
	Name        string `json:"name"        db:"name"`
	Sale        int    `json:"sale"        db:"sale"`
	Size        string `json:"size"        db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmID        int64  `json:"nm_id"       db:"nm_id"`
	Brand       string `json:"brand"       db:"brand"`
	Status      int    `json:"status"      db:"status"`
}
