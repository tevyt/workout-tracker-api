package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"tevyt.io/workout-planner/api/db"
	"tevyt.io/workout-planner/api/exercise"
)

func main() {
	dbConnection, err := db.InitializeDatabase()

	if err != nil {
		log.Fatal("Error initializing database connection, shutting down.")
	}

	defer dbConnection.Close()

	router := gin.Default()
	//exercise
	exerciseRepository := exercise.NewExerciseRepository(dbConnection)
	exerciseService := exercise.NewExerciseService(exerciseRepository)
	exerciseController := exercise.NewExerciseController(exerciseService)

	exerciseRoutes := router.Group("/exercise")
	exerciseRoutes.POST("/", exerciseController.CreateExercise)

	router.Run()

}
