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

func (exerciseService *ExerciseServiceMock) CreateExercise(name string, increment int8) (ExerciseModel, error) {
	if exerciseService.induceError {
		return ExerciseModel{}, errors.New("Error")
	}

	return ExerciseModel{Id: 1, ExerciseName: name, Increment: increment}, nil
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

func createGinContext(json string, method string) *gin.Context {
	httpRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(httpRecorder)
	context.Request = &http.Request{Header: make(http.Header)}

	context.Request.Method = method
	context.Request.Header.Set("content-type", "application/json")

	context.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(json)))

	return context
}
