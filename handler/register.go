package handler

import (
	"net/http"
	"ngc9/model"
	"ngc9/repo"
	"ngc9/errhandler"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo repo.UserRepo
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var u model.User

	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.Error(errhandler.ErrBindJSON )
		return
	}

	// validate user
	if u.Name == "" || u.Email == "" || u.Pwd == "" {
		ctx.Error(errhandler.ErrParam )
		return
	}

	newUser, err := h.Repo.Register(u)
	if err != nil {
		ctx.Error(err  )
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"message": "new user added",
			"user":    newUser,
		},
	)
}
