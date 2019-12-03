package repository

import (
	"database/sql"
	"log"
	"products_management/models"

	//postgres driver
	_ "github.com/lib/pq"
)

//DB ...
type DB interface {
	Execute(*models.Products) error
	GetProducts(*string) ([]models.Products, error)
	// Delete(*string) error
	// Update(*string, *bson.D) error
	// GetDetail(*bson.D, *models.Body) error

	//insert product category
	GetCategories(*string) ([]models.Category, error)
	QueryOnce(*string) (interface{}, error)
}

type dataBase struct {
	sqlDB *sql.DB
}

//GetPostgresSession for connect postgreSQL
func GetPostgresSession() DB {
	connStr := `postgres://admin:nimda@localhost:32784/products?sslmode=disable`
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return &dataBase{
		sqlDB: db,
	}
}

func (db *dataBase) GetProducts(sqlCommand *string) ([]models.Products, error) {
	rows, err := db.sqlDB.Query(*sqlCommand)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		return []models.Products{}, err
	}

	results := []models.Products{}
	for rows.Next() {
		var id, name, exp, cat string
		var amount, price int
		rows.Scan(&id, &name, &amount, &exp, &price, &cat)
		results = append(results, models.Products{
			ID:       id,
			Name:     name,
			Exp:      exp,
			Category: cat,
			Amount:   amount,
			Price:    price,
		})
	}

	return results, nil
}

func (db *dataBase) Execute(product *models.Products) error {
	stmt, err := db.sqlDB.Prepare("INSERT INTO product(id, product_name, amount, price, expire, category_id) VALUES(($1), ($2), ($3), ($4), ($5), ($6))")
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)

		return err
	}

	if _, err := stmt.Exec(product.ID, product.Name, product.Amount, product.Price, product.Exp, product.Category); err != nil {
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

func (db *dataBase) GetCategories(sqlCommand *string) ([]models.Category, error) {
	rows, err := db.sqlDB.Query(*sqlCommand)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	results := []models.Category{}
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		results = append(results, models.Category{id, name})
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
