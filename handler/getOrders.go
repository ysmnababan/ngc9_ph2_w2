package handler

import (
	"fmt"
	"net/http"
	"ngc9/repo"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Repo repo.ProductRepo
}

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	product, err := h.Repo.GetAllProducts()
	if err != nil {
		ctx.Error(err)
		return
	}

	fmt.Println(product)
	ctx.JSON(http.StatusOK, product)
}
