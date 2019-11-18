package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"products_management/api/v1/models"
	"products_management/utils"
)

//Client ...
type Client interface {
	GetAll() models.Products
}

type dataBase struct {
	client *mongo.Client
}

//GetDbSession return DB session
func GetDbSession() Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.Checker(err)

	return &dataBase{
		client: client,
	}
}

type Result struct {
	ID   string `bson: "_id"`
	Name string `bson: "name"`
	Exp  string `bson: "expDate"`
}

func (db *dataBase) GetAll() models.Products {
	collection := db.client.Database("products_management").Collection("products")

	result := Result{}
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&result)
	utils.Checker(err)

	return models.Products{
		Name: result.Name,
		ID:   result.ID,
		Exp:  result.Exp,
	}
}
