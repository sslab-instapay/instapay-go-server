package repository

import (
	"context"
	//"math/big"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/sslab-instapay/instapay-go-server/model"
	"github.com/sslab-instapay/instapay-go-server/db"
	"log"
)

func GetChannelList() ([]model.Channel, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("channels")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var channels []model.Channel

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var channel model.Channel
		err := cur.Decode(&channel)
		if err != nil {
			return nil, err
		}
		// To get the raw bson bytes use cursor.Current
		channels = append(channels, channel)
	}

	return channels, nil
}

func GetChannelById(channelId int) (model.Channel, error) {

	database, err := db.GetDatabase()
	if err != nil {
		return model.Channel{}, err
	}

	filter := bson.M{
		"cid": channelId,
	}

	collection := database.Collection("channels")

	channel := model.Channel{}
	singleRecord := collection.FindOne(context.TODO(), filter)
	if err := singleRecord.Decode(&channel); err != nil {
		log.Println(err)
	}
	return channel, nil
}

func PutChannelData(channelID int, from string, to string, deposit int) (*mongo.InsertOneResult, error) {

	database, err := db.GetDatabase()

	if err != nil {
		return nil, err
	}

	collection := database.Collection("channels")

	channel := model.Channel{ChannelId: channelID, From: from, To: to, Deposit: deposit}

	res, err := collection.InsertOne(context.TODO(), channel)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DropChannelData(channelID int) (*mongo.DeleteResult, error) {

	database, err := db.GetDatabase()

	if err != nil {
		return nil, err
	}

	collection := database.Collection("channels")

	res, err := collection.DeleteOne(context.TODO(), bson.M{"cid": channelID})
	if err != nil {
		return nil, err
	}

	return res, nil
}
