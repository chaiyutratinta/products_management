package repository

import (
	"database/sql"
	"fmt"
	"log"
	"products_management/configs"
	"products_management/models"

	//postgres driver
	_ "github.com/lib/pq"
)

//DB ...
type DB interface {
	InsertProduct(*models.Products) error
	GetProducts(*string) ([]models.Products, error)
	Delete(*string) error
	Update(*string, *string) error
	GetDetail(*string) (*models.Products, error)

	//insert product category
	GetCategories(*string) ([]models.Category, error)
	QueryOnce(*string) (interface{}, error)
	InsertCategory(*string, *string) error
	DeleteCategory(*string) error
}

type dataBase struct {
	sqlDB *sql.DB
}

//GetPostgresSession for connect postgreSQL
func GetPostgresSession() DB {
	conf := configs.Config.Database
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.Username, conf.Password, conf.Server, conf.DatabaseName)
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

func (db *dataBase) InsertProduct(product *models.Products) error {
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

func (db *dataBase) Delete(id *string) error {
	_, err := db.sqlDB.Exec("DELETE FROM product WHERE id=($1)", *id)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) Update(id, update *string) error {
	sqlCommand := fmt.Sprintf("UPDATE product SET %s WHERE id='%s'", *update, *id)
	stmt, err := db.sqlDB.Prepare(sqlCommand)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)

		return err
	}

	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) GetDetail(id *string) (*models.Products, error) {
	row := db.sqlDB.QueryRow(`
		SELECT product.id, product.product_name, product.amount, product.expire, product.price, product_category.category_name
		FROM product
		LEFT JOIN product_category
		ON product.category_id = product_category.id
		WHERE product.id = ($1)
		LIMIT 1
	`, *id)

	var pid, name, exp, cat string
	var amount, price int

	if err := row.Scan(&pid, &name, &amount, &exp, &price, &cat); err != nil {
		return &models.Products{}, err
	}

	results := models.Products{
		ID:       pid,
		Name:     name,
		Exp:      exp,
		Category: cat,
		Amount:   amount,
		Price:    price,
	}

	return &results, nil
}

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
		results = append(results, models.Category{
			ID:   id,
			Name: name,
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

func (db *dataBase) InsertCategory(id, categoryName *string) error {
	stmt, err := db.sqlDB.Prepare("INSERT INTO product_category VALUES(($1), ($2))")
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)

		return err
	}

	if _, err := stmt.Exec(*id, *categoryName); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (db *dataBase) DeleteCategory(id *string) error {
	_, err := db.sqlDB.Exec("DELETE FROM product_category WHERE id=($1)", *id)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
