package routes

import (
	"github.com/gorilla/mux"
	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/handler"
	"github.com/jithesh-wq/oolio-food-cart/service"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

func CreateRoutes(memoryStore *store.MemoryStore, db db.DbOperations) *mux.Router {

	// Create a new router
	r := mux.NewRouter()

	// Create services and handlers
	generateTableSessionService := service.CreateTableSessionService(db, memoryStore)
	generateTableSessionHandler := handler.CreateApiHandler(generateTableSessionService, memoryStore)
	r.HandleFunc("/generate-table-session", generateTableSessionHandler.HandleRequest).Methods("GET")

	removeTableSessionService := service.CreateRemoveTableSessionService(db, memoryStore)
	removeTableSessionHandler := handler.CreateApiHandler(removeTableSessionService, memoryStore)
	r.HandleFunc("/checkout-table-session", removeTableSessionHandler.HandleRequest).Methods("GET")

	getProductsService := service.CreateGetProductsService(db, memoryStore)
	getProductsHandler := handler.CreateApiHandler(getProductsService, memoryStore)
	r.HandleFunc("/products", getProductsHandler.HandleRequest).Methods("GET")

	addProductsService := service.CreateStoreProducts(db, memoryStore)
	addProductsHandler := handler.CreateApiHandler(addProductsService, memoryStore)
	r.HandleFunc("/add-products", addProductsHandler.HandleRequest).Methods("POST")

	placeOrderService := service.CreateOrderItems(db, memoryStore)
	placeOrderHandler := handler.CreateApiHandler(placeOrderService, memoryStore)
	r.HandleFunc("/order", placeOrderHandler.HandleRequest).Methods("POST")

	getOrdersService := service.CreateViewOrders(db, memoryStore)
	getOrdersHandler := handler.CreateApiHandler(getOrdersService, memoryStore)
	r.HandleFunc("/orders", getOrdersHandler.HandleRequest).Methods("GET")

	return r
}
