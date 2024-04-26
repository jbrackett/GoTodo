package main

import (
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *gorm.DB

func init() {
	// Initialize 'db'
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DATASOURCE_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Failed to connect to migrate database: %v", err)
	}
}

func main() {
	e := echo.New()

	e.POST("/todos", saveTodo)
	e.GET("/todos", getAllTodos)
	e.GET("/todos/:id", getTodo)
	e.PUT("/todos/:id", updateTodo)
	e.PATCH("/todos/:id", patchTodo)
	e.DELETE("/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(":1323"))
}

func saveTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}

	newId, _ := uuid.NewV7()
	todo.ID = &newId

	if todo.Completed == nil {
		defaultCompleted := false
		todo.Completed = &defaultCompleted
	}

	if todo.Message == nil {
		defaultMessage := ""
		todo.Message = &defaultMessage
	}

	db.Create(&todo)
	return c.JSON(http.StatusCreated, todo)
}

func getAllTodos(c echo.Context) error {
	var completed *bool
	completedParam, err := strconv.ParseBool(c.QueryParam("completed"))
	if err != nil {
		completed = nil
	} else {
		completed = &completedParam
	}
	var todos []Todo
	if completed == nil {
		db.Find(&todos)
	} else {
		db.Where("completed=?", completed).Find(&todos)
	}
	return c.JSON(http.StatusOK, todos)
}

func getTodo(c echo.Context) error {
	id := c.Param("id")
	var u uuid.UUID
	if err := u.Parse(id); err != nil {
		return err
	}
	todo := Todo{ID: &u}
	db.First(&todo)
	return c.JSON(http.StatusOK, todo)
}

func updateTodo(c echo.Context) error {
	id := c.Param("id")
	var u uuid.UUID
	if err := u.Parse(id); err != nil {
		return err
	}
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	if &u != todo.ID {
		return c.NoContent(http.StatusBadRequest)
	}
	db.Save(todo)
	return c.JSON(http.StatusOK, todo)
}

func patchTodo(c echo.Context) error {
	id := c.Param("id")
	var u uuid.UUID
	if err := u.Parse(id); err != nil {
		return err
	}

	existingTodo := Todo{ID: &u}
	db.First(&existingTodo)

	newTodo := new(Todo)
	if err := c.Bind(newTodo); err != nil {
		return err
	}

	if newTodo.Completed != nil {
		existingTodo.Completed = newTodo.Completed
	}
	if newTodo.Message != nil {
		existingTodo.Message = newTodo.Message
	}

	db.Save(existingTodo)
	return c.JSON(http.StatusOK, existingTodo)
}

func deleteTodo(c echo.Context) error {
	id := c.Param("id")
	var u uuid.UUID
	if err := u.Parse(id); err != nil {
		return err
	}
	db.Delete(&Todo{}, u)
	return c.NoContent(http.StatusOK)
}

type Todo struct {
	ID        *uuid.UUID `json:"id,omitempty"`
	Message   *string    `json:"message,omitempty"`
	Completed *bool      `json:"completed,omitempty"`
}
