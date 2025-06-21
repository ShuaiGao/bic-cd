package main

import (
	"bic-cd/cmd"
	"bic-cd/internal/util"
)

//	@title			BicCD api文档
//	@version		0.1
//	@contact.name	ZiJiaGao
//	@contact.email	boringmanman@qq.com

//	@host	bic-cd.farmergao.cn

var (
	gitTag string
)

func main() {
	util.SetServerTag(gitTag)
	cmd.Execute()
}
