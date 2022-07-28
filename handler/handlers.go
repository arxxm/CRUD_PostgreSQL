package handler

import (
	"github.com/arxxm/CRUD_test/api"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	repo *api.Repository
}

var prods = []api.Product{}

func NewAPIHandler(repo *api.Repository) (*APIHandler, error) {
	var h = APIHandler{}
	h.repo = repo
	return &h, nil
}

func (h *APIHandler) InitRoutes() *mux.Router {

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

	return rtr
}
