package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"github.com/sslab-instapay/instapay-go-server/model"
	"github.com/sslab-instapay/instapay-go-server/db"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetClientList() ([]model.Client, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("clients")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var clients []model.Client

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var client model.Client
		err := cur.Decode(&client)
		if err != nil {
			log.Println(err)
		}
		// To get the raw bson bytes use cursor.Current
		clients = append(clients, client)
	}

	return clients, nil
}

func GetClientInfo(address string) (*model.Client, error) {
	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("clients")

	cur, err := collection.Find(context.TODO(), bson.M{"address": address})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	cur.Next(context.Background())

	var info model.Client
	err = cur.Decode(&info)
	if err != nil {
		log.Println(err)
	}

	return &info, nil
}
