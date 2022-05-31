package main

import (
	bloc_client "github.com/fBloc/bloc-client-go"
	"github.com/fBloc/bloc-function/bloc_function_go/function/generate/time_about"
	"github.com/fBloc/bloc-function/bloc_function_go/function/intercept/str_intercept"
)

const appName = "bloc-function"

func main() {
	client := bloc_client.NewClient(appName)

	client.GetConfigBuilder().SetRabbitConfig(
		"blocRabbit", "blocRabbitPasswd", []string{"127.0.0.1:5672"}, "",
	).SetServer(
		"127.0.0.1", 8080,
	).BuildUp()

	generateGroup := client.RegisterFunctionGroup("生成器")
	generateGroup.AddFunction("时间", "生成时间", &time_about.TimeGen{})

	interceptGroup := client.RegisterFunctionGroup("拦截器")
	interceptGroup.AddFunction("字符串", "拦截字符串", &str_intercept.StrIntercept{})

	client.Run()
}
