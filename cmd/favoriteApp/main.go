package main

import (
	"favorite_service/internal/client"
	"favorite_service/internal/handlers"
	"favorite_service/internal/repositories"
	"favorite_service/internal/services"
	"favorite_service/pkg/psql"
	"fmt"
	"log"
	"os"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello Rest")

	app := fiber.New()

	err := godotenv.Load("../../.env")
	if err != nil {

		log.Fatal("Env Dosyası Yüklenemedi", err)

	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	var db = psql.Connect(host, user, password, name, port)

	itemRepository := repositories.NewFavoriteItemRepository(db)

	listRepository := repositories.NewFavoriteListRepository(db)

	itemService := services.NewFavoriItemService(itemRepository, listRepository)

	itemHandler := handlers.NewFavoriteItemHandler(itemService)

	itemHandler.SetRoutes(app)

	userClient := client.NewUserClient("http://localhost:3000/users")

	listService := services.NewFavoriteListService(listRepository, itemRepository, userClient)

	listHandler := handlers.NewFavoriteListHandler(listService)

	listHandler.SetRoutes(app)

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.ItemGetEndpoints())

	sw.AddEndpoints(handlers.ListGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":8080"))
}
