package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/sslab-instapay/instapay-go-server/model"
	"github.com/sslab-instapay/instapay-go-server/db"
	"log"
)

func PutPaymentData(pn int, from string, to string, amount int, p []string) (*mongo.InsertOneResult, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("payments")

	payment := model.Payment{
		PaymentNumber: pn,
		From:          from,
		To:            to,
		Amount:        amount,
		Participants:  p,
		AddrsSentAgr:  nil,
		AddrsSentUpt:  nil,
	}

	res, err := collection.InsertOne(context.TODO(), payment)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetPaymentData(pn int) (*model.Payment, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("payments")

	cur, err := collection.Find(context.TODO(), bson.M{"pn": int(pn)})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	cur.Next(context.Background())

	var pm model.Payment
	err = cur.Decode(&pm)
	if err != nil {
		log.Fatal(err)
	}

	return &pm, nil

}
func UpdatePaymentAddrsSentAgr(pn int, address string) (*mongo.UpdateResult, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("payments")

	cur, err := collection.Find(context.TODO(), bson.M{"pn": pn})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	cur.Next(context.Background())

	var pm model.Payment
	err = cur.Decode(&pm)
	if err != nil {
		log.Fatal(err)
	}

	pm.AddrsSentAgr = append(pm.AddrsSentAgr, address)

	res, err := collection.UpdateOne(context.TODO(), bson.M{"pn": pn}, bson.M{"$set": bson.M{"sentagr": pm.AddrsSentAgr}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdatePaymentAddrsSentUpt(pn int, address string) (*mongo.UpdateResult, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("payments")

	cur, err := collection.Find(context.TODO(), bson.M{"pn": pn})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	cur.Next(context.Background())

	var pm model.Payment
	err = cur.Decode(&pm)
	if err != nil {
		log.Fatal(err)
	}

	pm.AddrsSentUpt = append(pm.AddrsSentUpt, address)

	res, err := collection.UpdateOne(context.TODO(), bson.M{"pn": pn}, bson.M{"$set": bson.M{"sentupt": pm.AddrsSentUpt}})
	if err != nil {
		return nil, err
	}

	return res, nil
}
