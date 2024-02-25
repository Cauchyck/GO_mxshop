package api

import (
	"context"
	"fmt"
	"hello_go/mxshop/api/user_web/forms"
	"hello_go/mxshop/api/user_web/global"
	"hello_go/mxshop/api/user_web/global/response"
	"hello_go/mxshop/api/user_web/middlewares"
	"hello_go/mxshop/api/user_web/models"
	"hello_go/mxshop/api/user_web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
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
func GetUserList(ctx *gin.Context) {
	// // 从注册中心获取到用户服务的信息
	// cfg := api.DefaultConfig()
	// cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	// userSrvHost := ""
	// userSrvPort := 0
	// client, err := api.NewClient(cfg)

	// if err != nil {
	// 	panic(err)
	// }

	// data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Serbvice == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	// if err != nil {
	// 	panic(err)
	// }
	// for _, value := range data {
	// 	userSrvHost = value.Address
	// 	userSrvPort = value.Port
	// 	break
	// }

	// // ip := "127.0.0.1"
	// // port := 8888
	// userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	// if err != nil {
	// 	zap.S().Errorw("[GetUserList] connect user_servicer failed", "msg", err.Error())
	// 	return
	// }
	// userSrvClient := proto.NewUserClient(userConn)

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("user: %d", currentUser.ID)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] Get user list failed")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// zap.S().Debug("Get user list")

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		// data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			BirthDay: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		// data["id"] = value.Id
		// data["name"] = value.NickName
		// data["birthday"] = value.BirthDay
		// data["gender"] = value.Gender
		// data["mobile"] = value.Mobile

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
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

func PassWordLogin(c *gin.Context) {


	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "error",
		})
		return
	}


	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "user is not exist",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "login failed",
				})
			}
			return
		}
	} else {
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PassWordCheckInfo{
			PassWord:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "login failed",
			})
		} else {
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),
						ExpiresAt: time.Now().Unix() + 60*60*24*30,
						Issuer:    "imooc",
					},
				}

				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "generate token failed",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": time.Now().Unix() + 60*60*24*30 + 1000,
				})
			} else {
				c.JSON(http.StatusInternalServerError, map[string]string{
					"password": "login failed",
				})
			}
		}
	}
}

func Register(c *gin.Context) {

	registerForm := forms.RegisterForm{}

	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "sms error",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "sms code error",
			})
			return
		}
	}

	// zap.S().Info("[grpc.Dial]")
	// userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	// if err != nil {
	// 	zap.S().Errorw("[GetUserList] connect user_servicer failed", "msg", err.Error())
	// 	return
	// }
	// zap.S().Info("[proto.NewUserClient]")
	// userSrvClient := proto.NewUserClient(userConn)

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] build User failed: %s", err.Error())
		HandleValidatorError(c, err)
		return
	}

	// 生成token
	zap.S().Info("[middlewares.NewJWT()]")
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "imooc",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "generate token failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": time.Now().Unix() + 60*60*24*30 + 1000,
	})

}
