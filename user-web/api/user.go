package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
)


func removeTopStruct(fields map[string]string)map[string]string{
	rsp := map[string]string{}
	for field,err := range fields{
		rsp[field[strings.Index(field,".")+1:]] = err
	}

	return rsp
}



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
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context,err error){
	errs,ok := err.(validator.ValidationErrors)
	if !ok{
		c.JSON(http.StatusOK,gin.H{
			"msg":err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest,gin.H{
		"error":removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务] 失败",
			"msg", err.Error(),
		)
	}
	claims,_ := ctx.Get("claims")
	fmt.Println(claims.(*models.CustomClaims).ID)
	userSrvClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn","0")
	pnInt,_ := strconv.Atoi(pn)

	pSize := ctx.DefaultQuery("psize","10")
	pSizeInt,_ := strconv.Atoi(pSize)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			BirthDay: time.Time(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		//data["id"] = value.Id
		//data["name"]= value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)

	}

	ctx.JSON(http.StatusOK, result)

}

func PassWordLogin(c *gin.Context){
	passwordLoginForm := forms.PassWordLoginForm{}

	if err := c.ShouldBindJSON(&passwordLoginForm);err != nil{
		HandleValidatorError(c,err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId,passwordLoginForm.Captcha,true){
		c.JSON(http.StatusBadRequest,gin.H{
			"captcha":"验证码错误",
		})
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务] 失败",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(userConn)

	if rsp,err := userSrvClient.GetUserByMobile(context.Background(),&proto.MobileRequest{Mobile:passwordLoginForm.Mobile});err != nil{
		if e,ok := status.FromError(err);ok{
			switch e.Code(){
			case codes.NotFound:
				c.JSON(http.StatusBadRequest,map[string]string{
				"mobile":"用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError,map[string]string{
					"mobile":"登录失败",
				})
			}
			return
		}
	}else{
		if passRsp,passErr := userSrvClient.CheckPassWord(context.Background(),&proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		});passErr != nil{
			c.JSON(http.StatusInternalServerError,map[string]string{
				"password":"登录失败",
			})
		}else{
			if passRsp.Success{
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:             uint(rsp.Id),
					NickName:       rsp.NickName,
					AuthorityId:    uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore:time.Now().Unix(),
						ExpiresAt:time.Now().Unix()+60*60*24*30,
					},
				}
				token,err := j.CreateToken(claims)
				if err != nil{
					c.JSON(http.StatusInternalServerError,map[string]string{
						"msg":"生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK,gin.H{
					"id":rsp.Id,
					"nick_name":rsp.NickName,
					"token":token,
					"expired_at":(time.Now().Unix()+60*60*24*30)*1000,
				})
			}else{
				c.JSON(http.StatusBadRequest,map[string]string{
					"password":"登录失败",
				})
			}

		}
	}
}

func Register(c *gin.Context){
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBindJSON(&registerForm);err != nil{
		HandleValidatorError(c,err)
		return
	}

	//todo 验证码校验
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:               fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
	//})
	//value,err := rdb.Get(context.Background(),"key").Result()
	//if err == redis.Nil{
	//	c.JSON(http.StatusBadRequest,gin.H{
	//		"code":"验证码错误",
	//	})
	//	return
	//}else{
	//	if value != registerForm.Code{
	//		c.JSON(http.StatusBadRequest,gin.H{
	//			"code":"验证码错误",
	//		})
	//		return
	//	}
	//}


	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务] 失败",
			"msg", err.Error(),
		)
	}
	userSrvClient := proto.NewUserClient(userConn)
	user,err := userSrvClient.CreateUser(context.Background(),&proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil{
		zap.S().Errorf("[Register] [新建用户] 失败: %s",err.Error())
		HandleGrpcErrorToHttp(err,c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:             uint(user.Id),
		NickName:       user.NickName,
		AuthorityId:    uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore:time.Now().Unix(),
			ExpiresAt:time.Now().Unix()+60*60*24*30,
		},
	}
	token,err := j.CreateToken(claims)
	if err != nil{
		c.JSON(http.StatusInternalServerError,map[string]string{
			"msg":"生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"id":user.Id,
		"nick_name":user.NickName,
		"token":token,
		"expired_at":(time.Now().Unix()+60*60*24*30)*1000,
	})
}
