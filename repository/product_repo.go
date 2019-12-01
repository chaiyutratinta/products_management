package repository

import (
	"database/sql"
	"log"

	//postgres driver
	_ "github.com/lib/pq"
)

//DB ...
type DB interface {
	// GetAll() (*[]models.Products, error)
	Execute(*string) error
	// Delete(*string) error
	// Update(*string, *bson.D) error
	// GetDetail(*bson.D, *models.Body) error

	//insert product category
	QueryAll(*string) (*[]map[string]string, error)
	QueryOnce(*string) (interface{}, error)
}

type dataBase struct {
	sqlDB *sql.DB
}

//GetPostgresSession for connect postgreSQL
func GetPostgresSession() DB {
	connStr := `postgres://admin:nimda@localhost:32773/products?sslmode=disable`
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return &dataBase{
		sqlDB: db,
	}
}

// func (db *dataBase) GetAll() (*[]models.Products, error) {
// 	collection := db.client.
// 		Database("products_management").
// 		Collection("products")
// 	cur, err := collection.Find(context.TODO(), bson.D{{}})

// 	if err != nil {
// 		log.Fatal(err)
// 		return &[]models.Products{}, err
// 	}

// 	results := &[]models.Products{}
// 	for cur.Next(context.TODO()) {
// 		elem := struct {
// 			ID       string   `bson: "id"`
// 			Name     string   `bson: "name"`
// 			Exp      string   `bson: "exp"`
// 			Category []string `bson: "category"`
// 			Amount   int      `bson: "amount"`
// 		}{}
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		*results = append(*results, models.Products{
// 			ID:       elem.ID,
// 			Name:     elem.Name,
// 			Exp:      elem.Exp,
// 			Category: elem.Category,
// 			Amount:   elem.Amount,
// 		})
// 	}

// 	return results, nil
// }

func (db *dataBase) Execute(sqlCommand *string) error {
	_, err := db.sqlDB.Exec(*sqlCommand)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

// func (db *dataBase) Delete(id *string) error {
// 	collection := db.client.
// 		Database("products_management").
// 		Collection("products")
// 	_, err := collection.DeleteOne(context.TODO(), bson.D{{"id", *id}})

// 	if err != nil {
// 		log.Fatal(err)

// 		return err
// 	}

// 	return nil
// }

// func (db *dataBase) Update(id *string, update *bson.D) error {
// 	collection := db.client.
// 		Database("products_management").
// 		Collection("products")
// 	filter := bson.D{{"id", *id}}

// 	_, err := collection.UpdateOne(context.TODO(), filter, update)

// 	if err != nil {
// 		log.Fatal(err)

// 		return err
// 	}

// 	return nil
// }

// func (db *dataBase) GetDetail(filter *bson.D, result *models.Body) error {
// 	collection := db.client.
// 		Database("products_management").
// 		Collection("products")
// 	err := collection.
// 		FindOne(context.TODO(), filter).
// 		Decode(&result)

// 	if err != nil {
// 		log.Fatal(err)

// 		return err
// 	}

// 	return nil
// }

func (db *dataBase) QueryAll(sqlCommand *string) (*[]map[string]string, error) {
	rows, err := db.sqlDB.Query(*sqlCommand)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	results := &[]map[string]string{}
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		*results = append(*results, map[string]string{
			"id":   id,
			"name": name,
		})
	}

	return results, nil

}

func (db *dataBase) QueryOnce(sqlCommand *string) (interface{}, error) {
	rows, err := db.sqlDB.Query(*sqlCommand)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	var result interface{}
	for rows.Next() {
		rows.Scan(&result)
	}

	return result, nil
}
