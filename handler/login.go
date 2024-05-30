package handler

import (
	"fmt"
	"log"
	"net/http"
	"ngc9/errhandler"
	"ngc9/model"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func generateToken(u model.User) (string, error) {
	// create payload
	payload := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	}

	// define the method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("unable to get .env")
	}

	// get token string
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("unable to get token String")
	}

	return tokenString, nil
}
func (h *UserHandler) Login(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.Error(err)
		return
	}

	if u.Email == "" || u.Pwd == "" {
		c.Error(errhandler.ErrParam)
		return
	}

	newUser, err := h.Repo.Login(u)
	if err != nil {
		c.Error(err)

		return
	}

	token, err := generateToken(newUser)
	if err != nil {
		log.Println("unable to generate token:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "unable to generate token",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "login succeed",
			"token":   token,
		},
	)
}
