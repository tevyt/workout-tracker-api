package exercise

import (
	"database/sql"
	"errors"
	"testing"
)

type ExerciseRepositoryMock struct {
	exercises   []Exercise
	induceError bool
}

func (exerciseRepository *ExerciseRepositoryMock) CreateExercise(name string, increment int8) (Exercise, error) {
	if exerciseRepository.induceError {
		return Exercise{}, errors.New("Error creating exercise")
	}
	if exerciseRepository.exercises == nil {
		exerciseRepository.exercises = make([]Exercise, 0)
	}
	exercise := Exercise{ExerciseName: name, Increment: increment}
	exerciseRepository.exercises = append(exerciseRepository.exercises, exercise)

	return exercise, nil
}

func (exerciseRepository *ExerciseRepositoryMock) GetExercise(id int64) (Exercise, error) {
	if exerciseRepository.induceError {
		return Exercise{}, errors.New("Error")
	}
	if id == 1 {
		return Exercise{Id: 1, ExerciseName: "Deadlift", Increment: 10}, nil
	}

	return Exercise{}, sql.ErrNoRows
}

func (exerciseRepository *ExerciseRepositoryMock) SearchExercises(query string) ([]Exercise, error) {
	if exerciseRepository.induceError {
		return []Exercise{}, errors.New("Error")
	}
	if query == "Deadlift" {
		return []Exercise{{Id: 1, ExerciseName: "Deadlift", Increment: 10}}, nil
	}
	return make([]Exercise, 0), nil
}
func TestServiceCreateExercise(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{})

	exercise, err := exerciseService.CreateExercise("Deadlift", 10)

	if err != nil {
		t.Errorf("Unexpected error - %v\n", err)
	}

	if exercise.ExerciseName != "Deadlift" {
		t.Errorf("Returned exercise name is \"%s\", expected \"Deadlift\"\n", exercise.ExerciseName)
	}

	if exercise.Increment != 10 {
		t.Errorf("Returned increment is %d, expected 10", exercise.Increment)
	}
}

func TestServiceCreateExerciseFails(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{induceError: true})

	_, err := exerciseService.CreateExercise("Deadlift", 10)

	if err == nil {
		t.Error("Expected an error.")
	}
}

func TestServiceGetExercise(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{})

	exercise, err := exerciseService.GetExercise(1)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if exercise.Id != 1 {
		t.Errorf("Expected id to be 1 was %d\n", exercise.Id)
	}

	if exercise.ExerciseName != "Deadlift" {
		t.Errorf("Expected name to be \"Deadlift\" was \"%s\"", exercise.ExerciseName)
	}

	if exercise.Increment != 10 {
		t.Errorf("Expected increment to be 10 was %d\n", exercise.Increment)
	}
}

func TestGetExerciseReturnsExcersiceNotFoundWhenRecordDoesNotExist(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{})

	_, err := exerciseService.GetExercise(2)

	_, isExerciseNotFoundError := err.(ExerciseNotFoundError)

	if !isExerciseNotFoundError {
		t.Error("Expected ExerciseNotFoundError")
	}
}

func TestGetExerciseReturnsErrorWhenASystemErrorOccurs(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{induceError: true})

	_, err := exerciseService.GetExercise(1)

	_, isExerciseNotFoundError := err.(ExerciseNotFoundError)

	if isExerciseNotFoundError {
		t.Error("expected error to be propogated, wrapped error was returned.")
	}
}

func TestSearchExercisesReturnsAListOfExercises(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{induceError: false})

	exercises, err := exerciseService.SearchExercises("Deadlift")

	if err != nil {
		t.Error("Unexpected error")
	}

	if len(exercises) == 0 {
		t.Error("Array should not be empty.")
	}
}

func TestSearchExerciseReturnsAnErrorIfTheRepositoryErrors(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{induceError: true})

	_, err := exerciseService.SearchExercises("Deadlift")

	if err == nil {
		t.Error("Error should have been propogated")
	}

}
