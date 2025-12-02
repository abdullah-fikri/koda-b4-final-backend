package main

import (
	"backend/config"
	"backend/routes"
	

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDb()
	config.Redis()
	r := gin.Default()


	routes.Routes(r)

	r.Run(":8008")
}
