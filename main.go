package main

import (
	"github.com/akarshgangwar/examn_server/database"
	"github.com/akarshgangwar/examn_server/handlers"
	"github.com/gin-gonic/gin"
	"fmt"
)

func main()  {
	DbConnection, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB ",DbConnection)
	router := gin.Default()
	router.POST("auth/login",func(c *gin.Context) {
		handlers.LoginHandler(c, DbConnection)
	})
	router.Run("localhost:5000")

}
