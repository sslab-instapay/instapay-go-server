package model

import (
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"math/big"
)

type Client struct {
	PublicKeyAddress string          `bson:"address"`
	IP               string          `bson:"ip"`
	Port             int             `bson:"port"`
}

type Channel struct {
	ChannelId       int                `bson:"cid"`
	From            string             `bson:"from"`
	To              string             `bson:"to"`
	Deposit         int                `bson:"deposit"`
}

type Payment struct {
	PaymentNumber   int                `bson:"pn"`
	From            string             `bson:"from"`
	To              string             `bson:"to"`
	Amount          int                `bson:"amount"`
	Participants    []string           `bson:"participants"`
	AddrsSentAgr    []string           `bson:"sentagr"`
	AddrsSentUpt    []string           `bson:"sentupt"`
	Status          string             `bson:"status"`
}
