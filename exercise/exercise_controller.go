package exercise

import (
	"fmt"
	"net/http"
	"strconv"

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

	exercise, err := exerciseController.service.CreateExercise(createExerciseRequestBody.ExerciseName, createExerciseRequestBody.Increment)
	if err != nil {
		fmt.Printf("Error creating exercise - %v\n", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating exercise", "error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusCreated, exercise)
}

func (exerciseController *ExerciseController) GetExercise(ginContext *gin.Context) {
	id, err := strconv.ParseInt(ginContext.Param("id"), 10, 64)

	if err != nil {
		ginContext.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	exercise, err := exerciseController.service.GetExercise(id)

	if err != nil {
		_, isExerciseNotFoundError := err.(ExerciseNotFoundError)
		if isExerciseNotFoundError {
			ginContext.JSON(404, gin.H{"message": "id not found", "error": err.Error()})
		} else {
			ginContext.JSON(500, gin.H{"message": "Error fetching exercise", "error": err.Error()})
		}
		return
	}

	ginContext.JSON(200, exercise)
}
