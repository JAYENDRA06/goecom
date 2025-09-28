package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/JAYENDRA06/apiproject/service/cart"
	"github.com/JAYENDRA06/apiproject/service/order"
	"github.com/JAYENDRA06/apiproject/service/product"
	"github.com/JAYENDRA06/apiproject/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	log.Println("Listening on", s.addr)

	userStore := user.NewStore(s.db)
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	productService := product.NewHandler(productStore)
	productService.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subRouter)

	return http.ListenAndServe(s.addr, router)
}
