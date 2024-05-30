package main

import (
	"ngc9/config"
	"ngc9/errhandler"
	"ngc9/handler"
	"ngc9/middleware"
	"ngc9/model"
	"ngc9/repo"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectPostgre()
	Repo := &repo.PostgreRepo{DB: db}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.ProductDB{})
	db.AutoMigrate(&model.Store{})

	h := &handler.ProductHandler{Repo: Repo}
	userhandler := handler.UserHandler{Repo: Repo}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(errhandler.ErrorHandler())
	r.POST("/users/register", userhandler.Register)
	r.POST("/users/login", userhandler.Login)

	product := r.Group("/")
	product.Use(middleware.Auth())
	{
		product.GET("/products", h.GetProducts)
		product.GET("/product/:id", h.GetProductById)
		product.POST("/product", h.CreateProduct)
		product.PUT("/product/:id", h.UpdateProduct)
		product.DELETE("/product/:id", h.DeleteProduct)
	}

	r.Run(":8080")
}

// // Insert example data
// exampleStores := []*model.Store{
// 	{StoreName: "Store One", StorePwd: "password123", StoreEmail: "storeone@example.com", StoreType: "Retail"},
// 	{StoreName: "Store Two", StorePwd: "password456", StoreEmail: "storetwo@example.com", StoreType: "Wholesale"},
// 	{StoreName: "Store Three", StorePwd: "password789", StoreEmail: "storethree@example.com", StoreType: "Online"},
// }
// db.Create(exampleStores)
