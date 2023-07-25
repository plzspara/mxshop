package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"mxshop-api/form"
	"mxshop-api/global"
	"mxshop-api/proto"
	"mxshop-api/utils"
	"net/http"
	"strconv"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}

		}
		return
	}
}

func GetConnClient() (*grpc.ClientConn, error) {
	fmt.Println(global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorf("Error grpc Dial error: %v", err.Error())
		return nil, err
	}
	return conn, nil
}

func GetUserList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "10")
	page, err := strconv.Atoi(pn)
	pageSize, err := strconv.Atoi(pSize)
	if err != nil {
		zap.S().Errorw("strconv Atoi err: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	conn, err := GetConnClient()
	if err != nil {
		zap.S().Errorf("GetConnClient error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	defer conn.Close()
	client := proto.NewUserClient(conn)
	list, err := client.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(page),
		PSize: uint32(pageSize),
	})
	if err != nil {
		zap.S().Errorw("GetUserList error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	v := make([]global.UserResponse, 0, len(list.UserInfo))
	for _, e := range list.UserInfo {
		t := global.UserResponse{
			Id:       e.Id,
			NickName: e.Nickname,
			Birthday: global.JsonTime(time.Unix(int64(e.Birthday), 0)),
			Gender:   e.Gender,
			Mobile:   e.Mobile,
		}
		v = append(v, t)
	}
	ctx.JSON(http.StatusOK, v)
}

func PasswordLogin(ctx *gin.Context) {
	loginForm := form.PasswordLoginForm{}
	err := ctx.ShouldBind(&loginForm)
	if err != nil {
		zap.S().Errorw("LoginForm error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "账号密码不符合要求",
		})
		return
	}
	conn, err := GetConnClient()
	if err != nil {
		zap.S().Errorf("GetConnClient error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	defer conn.Close()
	client := proto.NewUserClient(conn)
	userInfoResponse, err := client.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
				return
			}
		}
	}
	rsp, err := client.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:          loginForm.Password,
		EncryptedPassword: userInfoResponse.Password,
	})
	if err != nil {
		zap.S().Errorw("CheckPassword error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登录失败",
		})
		return
	}
	if rsp.Success {
		claims := utils.CustomClaims{
			Id:          uint(userInfoResponse.Id),
			Nickname:    userInfoResponse.Nickname,
			AuthorityId: uint(userInfoResponse.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
				Issuer:    "plz",
			},
		}
		token, err := utils.NewJwt(claims, []byte(global.ServerConfig.JwtInfo.Key))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":        userInfoResponse.Id,
			"nickname":  userInfoResponse.Nickname,
			"token":     token,
			"expiredAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
			"msg":       "登录成功",
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
		})
	}
}
