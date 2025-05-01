package main

import "cengkeHelperBackGo/internal/router"

func main() {
	if err := router.Routers().Run(":" + "8080"); err != nil {
		panic(err)
		return
	}
}
