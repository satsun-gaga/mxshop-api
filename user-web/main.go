package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	myvalidator "mxshop-api/user-web/validator"
)

func main(){



	initialize.InitLogger()

	initialize.InitConfig()
	Router := initialize.Routers()

	_ = initialize.InitTrans("zh")

	if v,ok := binding.Validator.Engine().(*validator.Validate);ok{
		_  = v.RegisterValidation("mobile",myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile",global.Trans,func(ut ut.Translator)error{
			return ut.Add("mobile","{0} 非法的手机号码",true)},func(ut ut.Translator,fe validator.FieldError)string{
				t,_ := ut.T("mobile",fe.Field())

				return t
		})

	}

	zap.S().Infof("启动，端口：%d",global.ServerConfig.Port)

	if err := Router.Run(fmt.Sprintf(":%d",global.ServerConfig.Port));err != nil{
		zap.S().Panic("启动失败：",err.Error())
	}

}
