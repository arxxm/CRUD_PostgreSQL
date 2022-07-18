package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/arxxm/CRUD_test/api"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	// Initialize connection constants.
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "2245core"
	dbname   = "crud_pg"
)

func main() {

	psqlconn := fmt.Sprintf("host= %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	// psqlconn := fmt.Sprintf("user = %s password = %s dbname = %s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	repo := api.NewRepository(db)

	h, err := NewAPIHandler(repo)
	if err != nil {
		log.Fatal(err)
	}

	m := http.NewServeMux()
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", h.serveHome).Methods("GET")
	rtr.HandleFunc("/products/", h.productsPage).Methods("GET")
	rtr.HandleFunc("/product/edit/{id}", h.editProductById).Methods("GET")
	rtr.HandleFunc("/cmd/delete-product/{id}", h.deleteProduct)
	rtr.HandleFunc("/cmd/edit-product/{id}", h.editProduct)
	rtr.HandleFunc("/product/add", h.createProduct)
	rtr.HandleFunc("/cmd/add-product", h.addProduct)
	rtr.HandleFunc("/cmd/delete-all", h.deleteAll)
	rtr.HandleFunc("/q/product-search-by-name", h.searchByName)

	m.Handle("/", rtr)
	var srv = &http.Server{
		Addr:    ":8080",
		Handler: m,
		// TLSConfig: nil,

		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal()
	}
}

type APIHandler struct {
	repo *api.Repository
}

var prods = []api.Product{}

func (h *APIHandler) productsPage(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Products page")

	prods, err := h.repo.GetAllProducts()
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("templates/products.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "products", prods)
}

func NewAPIHandler(repo *api.Repository) (*APIHandler, error) {
	var h = APIHandler{}
	h.repo = repo

	return &h, nil
}

func (h *APIHandler) editProductById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/edit.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	P, err := h.repo.GetProduct(id)
	if err != nil {
		log.Fatal(err)
	}

	t.ExecuteTemplate(w, "edit", P)
}

func (h *APIHandler) editProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := r.FormValue("name")
	priceStr := r.FormValue("price")

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("name is: %s, price: %d, id: %d\n", name, price, id)
	err = h.repo.EditProduct(name, price, id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/products/", http.StatusSeeOther)

}

func (h *APIHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "create", nil)
}

func (h *APIHandler) searchByName(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	fmt.Printf("Имя: %s\n", name)
	p, err := h.repo.SearchByName(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Найденный товар: %s, цена: %d\n", p.Name, p.Price)

	prods = []api.Product{}
	prods = append(prods, p)
	t, err := template.ParseFiles("templates/found.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "found", prods)

}

func (h *APIHandler) addProduct(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	priceStr := r.FormValue("price")

	if name == "" || priceStr == "" {
		fmt.Fprintf(w, "Все поля должны быть заполнены")
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		log.Fatal(err)
	}

	if price <= 0 {
		fmt.Fprintf(w, "Цена не должна быть меньше или равна нулю")
	}

	err, ans := h.repo.CheckName(name)
	if err != nil {
		log.Fatal(err)
	}

	if ans {
		fmt.Fprintf(w, "Продукт с таким именем уже занесен в базу")
		fmt.Println("Продукт с таким именем уже занесен в базу")
		return
	}

	err = h.repo.AddProduct(name, price)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/products/", http.StatusSeeOther)
}

func (h *APIHandler) deleteAll(w http.ResponseWriter, r *http.Request) {
	h.repo.DeleteAll()
	http.Redirect(w, r, "/products/", http.StatusSeeOther)
}

func (h *APIHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {

	// var s []int
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	err = h.repo.DeleteProduct(id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/products/", http.StatusSeeOther)
}

func (h *APIHandler) serveHome(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/home.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "home", nil)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
