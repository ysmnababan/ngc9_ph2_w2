package errhandler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoRows        = errors.New("no rows in result set")
	ErrQuery         = errors.New("query execution failed")
	ErrScan          = errors.New("row scanning failed")
	ErrInvalidId     = errors.New("invalid id")
	ErrUserExists    = errors.New("user already exist")
	ErrRowsAffected  = errors.New("unable to get affected row")
	ErrNoAffectedRow = errors.New("rows affected is 0")
	ErrLastInsertId  = errors.New("unable to get last insert id")
	ErrNoUpdate      = errors.New("data already exists")
	ErrBindJSON      = errors.New("unable to bind json")
	ErrParam         = errors.New("error or missing parameter")
	ErrCredential    = errors.New("password or email doesn't match")
)

func parseError(err error, ctx *gin.Context) {
	fmt.Println("MIDDLEWARE")
	log.Println(err)
	status := http.StatusOK
	message := ""
	switch {
	case errors.Is(err, ErrQuery):
		fallthrough
	case errors.Is(err, ErrScan):
		fallthrough
	case errors.Is(err, ErrRowsAffected):
		fallthrough
	case errors.Is(err, ErrLastInsertId):
		fallthrough
	case errors.Is(err, ErrNoAffectedRow):
		status = http.StatusInternalServerError
		message = "Internal Server Error"
	case errors.Is(err, ErrNoRows):
		status = http.StatusNotFound
		message = "No row found"
	case errors.Is(err, ErrParam):
		status = http.StatusBadRequest
		message = "error or missing param"
	case errors.Is(err, ErrBindJSON):
		status = http.StatusBadRequest
		message = "Bad request"
	case errors.Is(err, ErrInvalidId):
		status = http.StatusBadRequest
		message = "Invalid ID"
	case errors.Is(err, ErrCredential):
		status = http.StatusBadRequest
		message = "Incorrect credential"
	case errors.Is(err, ErrUserExists):
		status = http.StatusBadRequest
		message = "User Already Exists"
	case errors.Is(err, ErrNoUpdate):
		status = http.StatusBadRequest
		message = "Data is the same"
	default:
		status = http.StatusInternalServerError
		message = "Unknown error:" + err.Error()
	}

	ctx.JSON(
		status,
		gin.H{
			"message": message,
		},
	)
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0].Err
			parseError(err, ctx)
		}
	}
}
