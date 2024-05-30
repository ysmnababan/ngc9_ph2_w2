package handler

import (
	"log"
	"net/http"
	"ngc9/model"
	"ngc9/errhandler"
	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var p model.ProductDB

	if err := ctx.ShouldBindJSON(&p); err != nil {
		log.Println(err)
		ctx.Error(errhandler.ErrBindJSON)
		return
	}

	// validate product
	if p.ID == 1 {
		ctx.Error(errhandler.ErrParam)
		return
	}

	newProducts, err := h.Repo.CreateProduct(p)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, newProducts)
}
