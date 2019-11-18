package repo

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"product_management/utils"
	"product_management/api/v1/models"
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
	ctx, _ 		:= context.WithTimeout(context.Background(), 10*time.Second)
	client, err 	:= mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	utils.Checker(err)

	return &dataBase{
		client: client,
	}
}

func (client *dataBase) GetAll() models.Products {

	return models.Products {
		Name	: "Dildo", 
		ID	: "id: wer23sdfq", 
		Exp	: "none",
	}
}

