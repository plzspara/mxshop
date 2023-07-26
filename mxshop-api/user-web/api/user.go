package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
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

func GetUserList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "10")
	page, err := strconv.Atoi(pn)
	pageSize, err := strconv.Atoi(pSize)
	if err != nil {
		zap.S().Errorw("strconv Atoi", "err", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	conn, err := utils.GetConnClient()
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
		zap.S().Errorw("GetUserList ", "error", err.Error())
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
		zap.S().Errorw("LoginForm ", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "账号密码不符合要求",
		})
		return
	}

	if !store.Verify(loginForm.CaptchaId, loginForm.Captcha, false) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	conn, err := utils.GetConnClient()
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
		zap.S().Errorw("GetUserByMobile  ", "error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	rsp, err := client.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:          loginForm.Password,
		EncryptedPassword: userInfoResponse.Password,
	})
	if err != nil {
		zap.S().Errorw("CheckPassword ", "error", err.Error())
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

func MobileLogin(ctx *gin.Context) {
	loginForm := form.MobileLoginForm{}
	err := ctx.ShouldBind(&loginForm)
	if err != nil {
		zap.S().Errorw("LoginForm", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "请输入正确的号码",
		})
		return
	}
	client := utils.GetRedisClient()
	defer client.Close()
	result, err := client.Get(context.Background(), "mobile:"+loginForm.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "号码已过期",
			})
			return
		}
		zap.S().Errorw("LoginForm", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	if result != loginForm.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "号码错误",
		})
		return
	}
	connClient, err := utils.GetConnClient()
	if err != nil {
		zap.S().Errorw("LoginForm ", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	defer connClient.Close()
	userClient := proto.NewUserClient(connClient)
	userInfoResponse, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: loginForm.Mobile})
	if err != nil {
		zap.S().Errorw("GetUserByMobile ", "error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
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
}

func Register(ctx *gin.Context) {
	var userRegister form.UserRegister
	err := ctx.ShouldBind(&userRegister)
	if err != nil {
		zap.S().Errorw("LoginForm", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "请输入正确的号码",
		})
		return
	}
	client := utils.GetRedisClient()
	defer client.Close()
	result, err := client.Get(context.Background(), "mobile:"+userRegister.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "号码已过期",
			})
			return
		}
		zap.S().Errorw("LoginForm", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	if result != userRegister.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "号码错误",
		})
		return
	}
	connClient, err := utils.GetConnClient()
	if err != nil {
		zap.S().Errorw("LoginForm ", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	defer connClient.Close()
	userClient := proto.NewUserClient(connClient)
	_, err = userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: userRegister.Password,
		Mobile:   userRegister.Mobile,
	})
	if err != nil {
		zap.S().Errorw("CreateUser ", "error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误,: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64,
	})
}

func SendSms(c *gin.Context) {
	ctx := context.Background()
	var request form.GetSmsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请输入正确的手机号",
		})
		return
	}
	rdb := utils.GetRedisClient()
	defer rdb.Close()
	result, err := rdb.Do(ctx, "TTL", "mobile:"+request.Mobile).Result()
	if err != nil {
		if err != redis.Nil {
			zap.S().Errorw("GetRedisClient.Get ", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部错误",
			})
			return
		}
	}
	if result.(int64) != -2 && 600-result.(int64) < 60 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请稍后重试",
		})
		return
	}

	err = rdb.Set(ctx, "mobile:"+request.Mobile, utils.GetMsmCode(), time.Minute*10).Err()
	if err != nil {
		zap.S().Errorw("RedisClient.set ", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
