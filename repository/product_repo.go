package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"products_management/models"
	"products_management/utils"
)

//Client ...
type Client interface {
	GetAll() (*[]models.Products, error)
	Add(*models.Products) error
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

func (db *dataBase) GetAll() (*[]models.Products, error) {
	collection := db.client.Database("products_management").Collection("products")
	cur, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return &[]models.Products{}, err
	}

	results := &[]models.Products{}
	for cur.Next(context.TODO()) {
		elem := struct {
			ID       string   `bson: "id"`
			Name     string   `bson: "name"`
			Exp      string   `bson: "exp"`
			Category []string `bson: "category"`
			Amount   int      `bson: "amount"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}

		*results = append(*results, models.Products{
			ID:       elem.ID,
			Name:     elem.Name,
			Exp:      elem.Exp,
			Category: elem.Category,
			Amount:   elem.Amount,
		})
	}

	return results, nil
}

func (db *dataBase) Add(product *models.Products) error {
	collection := db.client.Database("products_management").Collection("products")
	_, err := collection.InsertOne(context.TODO(), *product)

	if err != nil {
		return err
	}

	return nil
}
