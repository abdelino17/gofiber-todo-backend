package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abdelino17/gofiber-todo-api/database"
	"github.com/abdelino17/gofiber-todo-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	sanityCheck()

	app := fiber.New()
	app.Use(cors.New())
	initDatabase()

	setupRoutes(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func sanityCheck() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envProps := []string{
		"PORT",
		"PG_HOST",
		"PG_PORT",
		"PG_USERNAME",
		"PG_PASSWORD",
		"PG_DATABASE",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			log.Println(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}

func initDatabase() {
	pgUser := os.Getenv("PG_USERNAME")
	pgPass := os.Getenv("PG_PASSWORD")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDB := os.Getenv("PG_DATABASE")

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s", pgHost, pgUser, pgPass, pgPort, pgDB)
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connected")
	database.DBConn.AutoMigrate(&models.Todo{})
	fmt.Println("Migrated DB")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/todos", models.GetTodos)
	app.Get("/todos/:id", models.GetTodoById)
	app.Post("/todos", models.CreateTodo)
	app.Put("/todos/:id", models.UpdateTodo)
	app.Delete("/todos/:id", models.DeleteTodo)
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}
