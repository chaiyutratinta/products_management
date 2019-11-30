package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	//postgres driver
	_ "github.com/lib/pq"
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
	Delete(*string) error
	Update(*string, *bson.D) error
	GetDetail(*bson.D, *models.Body) error

	//insert product category
	AddProtuctCatgegory(*string) error
	GetProductCategories(*string) (*[]map[string]string, error)
	DeleteProductCategory(*string) error
}

type dataBase struct {
	client *mongo.Client
	sqlDB *sql.DB
}

//GetDbSession return DB session
func GetDbSession() Client {
	client, err := mongo.
		Connect(context.TODO(), options.Client().
			ApplyURI("mongodb://postgres_db"))
	utils.Checker(err)

	return &dataBase{
		client: client,
	}
}

//GetPostgresSession for connect postgreSQL
func GetPostgresSession() Client {
	connStr := `postgres://admin:nimda@localhost:32770/products?sslmode=disable`
	db, err := sql.Open("postgres", connStr)
	
	if err != nil {
		log.Fatal(err)
	}

	return &dataBase{
		sqlDB: db,
	}
}

func (db *dataBase) GetAll() (*[]models.Products, error) {
	collection := db.client.
		Database("products_management").
		Collection("products")
	cur, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
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
	collection := db.client.
		Database("products_management").
		Collection("products")
	_, err := collection.InsertOne(context.TODO(), *product)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) Delete(id *string) error {
	collection := db.client.
		Database("products_management").
		Collection("products")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{"id", *id}})

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) Update(id *string, update *bson.D) error {
	collection := db.client.
		Database("products_management").
		Collection("products")
	filter := bson.D{{"id", *id}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) GetDetail(filter *bson.D, result *models.Body) error {
	collection := db.client.
		Database("products_management").
		Collection("products")
	err := collection.
		FindOne(context.TODO(), filter).
		Decode(&result)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) AddProtuctCatgegory(sqlCommand *string) error {
	_, err := db.sqlDB.Exec(*sqlCommand)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) GetProductCategories(sqlCommand *string) (*[]map[string]string, error) {
	rows, err := db.sqlDB.Query(*sqlCommand)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)

		return nil , err
	}
	
	results := &[]map[string]string{}
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		*results = append(*results, map[string]string{
			"id": id,
			"name": name,
		})
	}

	return results, nil

}

func (db *dataBase) DeleteProductCategory(sqlCommand *string) error {
	_, err := db.sqlDB.Exec(*sqlCommand)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}