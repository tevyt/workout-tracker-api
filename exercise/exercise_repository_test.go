package exercise

import (
	"database/sql"
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
	exercise, err := exerciseRepository.CreateExercise("Deadlift", 10)

	if err != nil {
		t.Errorf("Unexpected error - %v\n", err)
	}

	if exercise.Id != 1 {
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

func TestGetExercise(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Error("Unable to create database mock")
	}

	result := sqlmock.NewRows([]string{"id", "exercise_name", "increment"}).AddRow(1, "Deadlift", 10)
	mock.ExpectQuery("SELECT id, exercise_name, increment FROM exercise").WithArgs(1).WillReturnRows(result)
	exerciseRepository := NewExerciseRepository(db)

	exercise, err := exerciseRepository.GetExercise(1)

	if err != nil {
		t.Errorf("Unexpected error %v\n", err)
	}

	if exercise.Id != 1 {
		t.Errorf("Expected id 1 was %d\n", exercise.Id)
	}

	if exercise.ExerciseName != "Deadlift" {
		t.Errorf("Expected \"Deadlift\" was \"%s\"\n", exercise.ExerciseName)
	}

	if exercise.Increment != 10 {
		t.Errorf("Expected 10 was %d\n", exercise.Increment)
	}
}

func TestGetExerciseError(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Error("Unable to create database mock")
	}

	mock.ExpectQuery("SELECT id, exercise_name, increment FROM exercise").WithArgs(1).WillReturnError(sql.ErrNoRows)
	exerciseRepository := NewExerciseRepository(db)

	_, err = exerciseRepository.GetExercise(1)

	if err == nil {
		t.Error("Expected an error.")
	}
}

func TestSearchExercisesReturnsAnArrayOnSuccess(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Error("Unable to create database mock")
	}

	result := sqlmock.NewRows([]string{"id", "exercise_name", "increment"}).AddRow(1, "Deadlift", 10)
	mock.ExpectQuery("SELECT id, exercise_name, increment FROM exercise").WithArgs("%Deadlift%").WillReturnRows(result)

	exerciseRepository := NewExerciseRepository(db)

	exercises, err := exerciseRepository.SearchExercises("Deadlift")

	if err != nil {
		t.Error("Unexpected error")
	}

	if len(exercises) == 0 {
		t.Error("Expected entry to be returned")
	}
}

func TestSearchExercisesReturnsAnEmptyArrayIfNoRowsAreReturned(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Error("Unable to create database mock")
	}
	result := sqlmock.NewRows([]string{"id", "exercise_name", "increment"})
	mock.ExpectQuery("SELECT id, exercise_name, increment FROM exercise").WithArgs("%Deadlift%").WillReturnRows(result)

	exerciseRepository := NewExerciseRepository(db)

	exercises, err := exerciseRepository.SearchExercises("Deadlift")

	if err != nil {
		t.Error("Unexpected error")
	}
	if exercises == nil || len(exercises) != 0 {
		t.Error("Expected an empty array")
	}
}

func TestSearchExercisesPropogatesErrorsFromDB(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Error("Unable to create database mock")
	}
	mock.ExpectQuery("SELECT id, exercise_name, increment FROM exercise").WithArgs("%Deadlift%").WillReturnError(errors.New("Error"))

	exerciseRepository := NewExerciseRepository(db)

	_, err = exerciseRepository.SearchExercises("Deadlift")

	if err == nil {
		t.Error("Expected an error")
	}

}
