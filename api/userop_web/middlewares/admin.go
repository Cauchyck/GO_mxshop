package middlewares

import (
	"hello_go/mxshop/api/userop_web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdminAuth() gin.HandlerFunc{
	return func(ctx *gin.Context){
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "No Authority",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
 	}
}