package api

import (
	"database/sql"
	"fmt"
)

type Bucket string

const (
	ProductsBucket Bucket = "products_bucket"
	IndexBucket    Bucket = "index_bucket"
)

type Repository struct {
	db *sql.DB
}

type Product struct {
	Id    int
	Name  string
	Price int
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddProduct(name string, price int) error {

	_, err := r.db.Exec("insert into products (prod_name, price) values ($1, $2)", name, price)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAll() error {

	_, err := r.db.Exec("DELETE FROM products WHERE prod_name <> ''")
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllProducts() ([]Product, error) {

	products := []Product{}

	rows, err := r.db.Query("select * from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Name, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) SearchByName(name string) (Product, error) {

	// var id int
	prod := Product{}
	row := r.db.QueryRow("SELECT * FROM products WHERE prod_name = $1", name)

	err := row.Scan(&prod.Id, &prod.Name, &prod.Price)
	if err != nil {
		return Product{}, err
	}

	return prod, nil
}

func (r *Repository) CheckName(name string) (error, bool) {

	rows, err := r.db.Query("select * from products where prod_name = $1", name)
	if err != nil {
		return err, false
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {

		p := Product{}
		err := rows.Scan(&p.Id, &p.Name, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	if len(products) == 0 {
		return nil, false
	} else {
		return nil, true
	}
}

func (r *Repository) GetProduct(id int) (Product, error) {

	prod := Product{}

	row := r.db.QueryRow("SELECT * FROM products where prod_id = $1", id)
	err := row.Scan(&prod.Id, &prod.Name, &prod.Price)
	if err != nil {
		return prod, err
	}
	return prod, nil
}

func (r *Repository) EditProduct(newName string, newPrice, id int) error {

	_, err := r.db.Exec("UPDATE products SET prod_name = $1, price = $2 WHERE prod_id = $3", newName, newPrice, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteProduct(id int) error {

	_, err := r.db.Exec("DELETE FROM products WHERE prod_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
