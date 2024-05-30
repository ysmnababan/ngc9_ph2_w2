package handler

import (
	"net/http"
	"strconv"
	"ngc9/errhandler"
	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil || id <= 0 {
		c.Error(errhandler.ErrInvalidId)
		return
	}

	err = h.Repo.DeleteProduct(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "data deleted successfuly",
		},
	)
}
