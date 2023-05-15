package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jaga1999/go-rest-api-postgres/models"
	"github.com/Jaga1999/go-rest-api-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create Book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book has been created"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookmodel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(bookmodel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "User deleted successfully!!!",
	})
	return nil

}

func (r *Repository) GetBookbyID(context *fiber.Ctx) error {
	id := context.Params("id")

	bookmodel := &models.Books{}
	if id == "" {
		context.Status(http.StatusBadGateway).JSON(&fiber.Map{
			"message": "Id cannot be empty",
		})
		return nil
	}
	err := r.DB.Where("id = ?", id).First(bookmodel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not fetch Book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book fetched successfully",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookmodels := &[]models.Books{}

	err := r.DB.Find(bookmodels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get Books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched successfully",
			"data":    bookmodels,
		})
	return nil
}

func (r *Repository) setupRouter(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/create", r.CreateBook)
	api.Get("/book/:id", r.GetBookbyID)
	api.Delete("/delete/:id", r.DeleteBook)
	api.Get("/books", r.GetBooks)
}

func main() {

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	r := Repository{
		DB: db,
	}

	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	app := fiber.New()
	r.setupRouter(app)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	app.Listen("0.0.0.0:" + port)
}
