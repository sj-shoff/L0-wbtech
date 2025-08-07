package model

import "time"

type Order struct {
	OrderUID          string    `db:"order_uid" json:"order_uid"`
	TrackNumber       string    `db:"track_number" json:"track_number"`
	Entry             string    `db:"entry" json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `db:"locale" json:"locale"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerID        string    `db:"customer_id" json:"customer_id"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service"`
	Shardkey          string    `db:"shardkey" json:"shardkey"`
	SmID              int       `db:"sm_id" json:"sm_id"`
	DateCreated       time.Time `db:"date_created" json:"date_created"`
	OofShard          string    `db:"oof_shard" json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
