package handler

import (
	"net/http"
	"ngc9/errhandler"
	"ngc9/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil || id <= 0 {
		c.Error(errhandler.ErrInvalidId)
		return
	}

	var p model.ProductDB
	if err := c.ShouldBindJSON(&p); err != nil {
		c.Error(errhandler.ErrBindJSON)
		return
	}

	// validate
	if false {
		c.Error(errhandler.ErrParam)
		return
	}

	err = h.Repo.UpdateProduct(id, p)
	if err != nil {
		c.Error(err)
		return
	}
	p.ID = uint(id)
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "data updated",
			"product": p,
		},
	)
}
