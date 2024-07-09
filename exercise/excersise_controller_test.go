package exercise

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type ExerciseServiceMock struct {
	induceError bool
}

func (exerciseService *ExerciseServiceMock) CreateExercise(name string, increment int8) (Exercise, error) {
	if exerciseService.induceError {
		return Exercise{}, errors.New("Error")
	}

	return Exercise{Id: 1, ExerciseName: name, Increment: increment}, nil
}

func (exerciseService *ExerciseServiceMock) GetExercise(id int64) (Exercise, error) {
	if exerciseService.induceError {
		return Exercise{}, errors.New("Error")
	}
	if id == 1 {
		return Exercise{Id: 1, ExerciseName: "Deadlift", Increment: 10}, nil
	}
	return Exercise{}, ExerciseNotFoundError{Id: id}
}

func TestCreateExerciseReturns201WhenSuccessful(t *testing.T) {
	requestJSON := "{\"exerciseName\":\"Deadlift\", \"increment\":10}"
	context := createGinContext(requestJSON, "POST")

	exerciseController := NewExerciseController(&ExerciseServiceMock{})

	exerciseController.CreateExercise(context)

	if context.Writer.Status() != http.StatusCreated {
		t.Errorf("Expected 201 status was %d", context.Writer.Status())
	}
}

func TestCreateExerciseReturns500WhenTheresAnError(t *testing.T) {
	requestJSON := "{\"exerciseName\":\"Deadlift\", \"increment\":10}"
	context := createGinContext(requestJSON, "POST")

	exerciseController := NewExerciseController(&ExerciseServiceMock{induceError: true})

	exerciseController.CreateExercise(context)

	if context.Writer.Status() != http.StatusInternalServerError {
		t.Errorf("Expected status 500 was %d\n", context.Writer.Status())
	}
}

func TestCreateExerciseReturns400IfJSONIsUnparsable(t *testing.T) {
	requestJSON := "{\"invalid\"}"
	context := createGinContext(requestJSON, "POST")
	exerciseController := NewExerciseController(&ExerciseServiceMock{induceError: false})

	exerciseController.CreateExercise(context)

	if context.Writer.Status() != 400 {
		t.Errorf("Expected status 400 was %d\n", context.Writer.Status())
	}

}

func TestGetExercise200WhenExerciseExists(t *testing.T) {
	context := createGinContext("", "GET")
	context.Params = append(context.Params, gin.Param{Key: "id", Value: "1"})
	exerciseController := NewExerciseController(&ExerciseServiceMock{induceError: false})

	exerciseController.GetExercise(context)

	if context.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status to be %d was %d", http.StatusOK, context.Writer.Status())
	}
}

func TestGetServiceReturns404WhenExerciseDoesNotExist(t *testing.T) {
	context := createGinContext("", "GET")
	context.Params = append(context.Params, gin.Param{Key: "id", Value: "2"})
	exerciseController := NewExerciseController(&ExerciseServiceMock{induceError: false})

	exerciseController.GetExercise(context)

	if context.Writer.Status() != http.StatusNotFound {
		t.Errorf("Expected status to be %d was %d", http.StatusNotFound, context.Writer.Status())
	}
}

func TestGetServiceReturns500WhenThereIsASystemError(t *testing.T) {
	context := createGinContext("", "GET")
	context.Params = append(context.Params, gin.Param{Key: "id", Value: "1"})
	exerciseController := NewExerciseController(&ExerciseServiceMock{induceError: true})

	exerciseController.GetExercise(context)

	if context.Writer.Status() != http.StatusInternalServerError {
		t.Errorf("Expected status to be %d was %d", http.StatusInternalServerError, context.Writer.Status())
	}
}

func createGinContext(json string, method string) *gin.Context {
	httpRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(httpRecorder)
	context.Request = &http.Request{Header: make(http.Header)}

	context.Request.Method = method
	context.Request.Header.Set("content-type", "application/json")

	context.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(json)))

	return context
}
