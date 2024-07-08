package exercise

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExerciseController struct {
	service ExerciseService
}

func NewExerciseController(service ExerciseService) *ExerciseController {
	return &ExerciseController{
		service: service,
	}
}

func (exerciseController *ExerciseController) CreateExercise(ginContext *gin.Context) {

	var createExerciseRequestBody struct {
		ExerciseName string `json:"exerciseName"`
		Increment    int8   `json:"increment"`
	}

	err := ginContext.BindJSON(&createExerciseRequestBody)

	if err != nil {
		fmt.Printf("Error reading request json - %v\n", err)
		ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Unable to read request", "error": err.Error()})
		return

	}

	model, err := exerciseController.service.CreateExercise(createExerciseRequestBody.ExerciseName, createExerciseRequestBody.Increment)
	if err != nil {
		fmt.Printf("Error creating exercise - %v\n", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating exercise", "error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusCreated, model)
}
