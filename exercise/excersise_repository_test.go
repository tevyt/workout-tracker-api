package exercise

import (
	"errors"
	"testing"

	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRepositoryCreateExercise(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Errorf("Error creating db connection: %v", err)
	}

	result := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO exercise").WithArgs("Deadlift", 10).WillReturnRows(result)

	exerciseRepository := NewExerciseRepository(db)
	model, err := exerciseRepository.CreateExercise("Deadlift", 10)

	if err != nil {
		t.Errorf("Unexpected error - %v\n", err)
	}

	if model.Id != 1 {
		t.Errorf("id not set on result")
	}
}

func TestRepositoryCreateExerciseFails(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Errorf("Error creating db connection: %v", err)
	}

	mock.ExpectQuery("INSERT INTO exercise").WithArgs("Deadlift", 10).WillReturnError(errors.New("Error with insert"))

	exerciseRepository := NewExerciseRepository(db)
	_, err = exerciseRepository.CreateExercise("Deadlift", 10)

	if err == nil {
		t.Error("Expected error to be propogated")
	}

}
