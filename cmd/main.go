package main

import (
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/router"
)

func main() {
	if err := router.Routers().Run(":" + config.Conf.Server.Port); err != nil {
		panic(err)
		return
	}
}
