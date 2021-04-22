package api

import (

	"fmt"
	"github.com/gin-gonic/gin"


	"math/rand"
	"strings"
	"time"
)

//todo 发送短信验证码

func GenerateSmsCode(width int) string{
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}

	r := len(numeric)

	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder

	for i := 0;i < width;i ++{
		fmt.Fprintf(&sb,"%d",numeric[rand.Intn(r)])
	}

	return sb.String()
}

func SendSms(ctx *gin.Context){
	//sendSmsForm := forms.SendSmsForm{}
	//
	//if err := ctx.ShouldBindJSON(&sendSmsForm);err != nil{
	//	HandleValidatorError(ctx,err)
	//}
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:               fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
	//})
	//rdb.Set(context.Background(),"key","value",300*time.Second)
}
