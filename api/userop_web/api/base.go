package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "Argument error",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "other error",
				})
			}
		}
	}

	return
}

func HandleValidatorError(c *gin.Context, err error) {
	zap.S().Info("[HandleValidatorError")
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("err.Error(): %v", err.Error()),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errs.Error(),
	})
}
