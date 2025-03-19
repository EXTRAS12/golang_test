package main

import (
	"log"

	"songapp/config"
	"songapp/handlers"
	"songapp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db := config.InitDB()

	// Создание обработчиков
	songHandler := handlers.NewSongHandler(db)

	// Настройка маршрутизации
	r := gin.Default()
	routes.SetupRoutes(r, songHandler)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
