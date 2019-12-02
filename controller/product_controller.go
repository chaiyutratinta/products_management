package controller

import (
	"fmt"
	"log"
	"products_management/models"
	"products_management/repository"

	"github.com/google/uuid"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() []models.Products
	AddProduct(*models.Products) error
	// DeleteProduct(*string) error
	// UpdateProduct(*string, *models.Body) error
	// GetDetailProduct(*string) (*models.Products, error)

	//insert product category
	InsertProductCategory(*string) error
	SelectAllProductCategories() ([]map[string]string, error)
	RemoveProductCategory(*string) error
	IsCategoryMatch(*string) bool
}

//ProductsUseCase ...
type productController struct {
	repository.DB
}

//NewController ...
func NewController(db repository.DB) ProductController {

	return &productController{db}
}

func (r *productController) GetAllProduct() []models.Products {
	sqlCommand := fmt.Sprintf("SELECT * FROM product")
	products, err := r.GetProducts(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return []models.Products{}
	}

	return products
}

func (r *productController) AddProduct(product *models.Products) error {
	sqlCommand := fmt.Sprintf(`
				INSERT INTO product(id, product_name, amount, price, expire, category_id)
				VALUES('%s', '%s', %d, %d, '%s', '%s')
				`, product.ID, product.Name, product.Amount, product.Price, product.Exp, product.Category)
	err := r.Execute(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

// func (r *productController) DeleteProduct(id *string) error {
// 	err := r.Delete(id)

// 	if err != nil {
// 		log.Fatal(err)

// 		return err
// 	}

// 	return nil
// }

// func (r *productController) UpdateProduct(id *string, body *models.Body) error {
// 	fields := &bson.D{}

// 	if body.Name != "" {
// 		*fields = append(*fields, bson.E{"name", body.Name})
// 	}
// 	if body.Exp != "" {
// 		*fields = append(*fields, bson.E{"exp", body.Exp})
// 	}
// 	if len(body.Category) > 0 {
// 		*fields = append(*fields, bson.E{"category", body.Category})
// 	}
// 	if body.Amount != 0 {
// 		*fields = append(*fields, bson.E{"amount", body.Amount})
// 	}

// 	update := bson.D{{"$set", *fields}}
// 	err := r.Update(id, &update)

// 	if err != nil {
// 		log.Fatal(err)

// 		return err
// 	}

// 	return nil
// }

// func (r *productController) GetDetailProduct(id *string) (*models.Products, error) {
// 	filter := bson.D{{"id", *id}}
// 	result := &models.Body{}
// 	err := r.GetDetail(&filter, result)

// 	if err != nil {
// 		log.Fatal(err)

// 		return nil, err
// 	}

// 	return &models.Products{
// 		Name:     result.Name,
// 		Exp:      result.Exp,
// 		Category: result.Category,
// 		Amount:   result.Amount,
// 	}, nil
// }

func (r *productController) InsertProductCategory(categoryName *string) error {
	id := (uuid.New()).String()
	sqlCommand := fmt.Sprintf(`INSERT INTO product_category VALUES('%s', '%s')`, id, *categoryName)
	err := r.Execute(&sqlCommand)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *productController) SelectAllProductCategories() ([]map[string]string, error) {
	sqlCommand := `SELECT * FROM product_category`
	results, err := r.GetCategories(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	categories := []map[string]string{}
	for _, elm := range results {
		elmStruct := elm.(models.Category)
		categories = append(categories, map[string]string{
			"id":   elmStruct.ID,
			"name": elmStruct.Name,
		})
	}

	return categories, nil
}

func (r *productController) RemoveProductCategory(id *string) error {
	sqlCommand := fmt.Sprintf(`DELETE FROM product_category WHERE id='%s'`, *id)
	err := r.Execute(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) IsCategoryMatch(id *string) bool {
	sqlCommand := fmt.Sprintf(`SELECT EXISTS (SELECT true FROM product_category WHERE id='%s')`, *id)
	result, err := r.QueryOnce(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return false
	}

	ismatch, _ := result.(bool)

	return ismatch
}
