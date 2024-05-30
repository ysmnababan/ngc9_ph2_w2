package handler

import (
	"net/http"
	"strconv"
	"ngc9/errhandler"
	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) GetProductById(ctx *gin.Context) {
	param_id := ctx.Param("id")

	id, err := strconv.Atoi(param_id)
	if err != nil || id <= 0 {
		ctx.Error(errhandler.ErrInvalidId )
		return
	}

	product, err := h.Repo.GetProductById(uint(id))
	if err != nil {
		ctx.Error(err )
		return
	}

	ctx.JSON(http.StatusOK, product)
}
