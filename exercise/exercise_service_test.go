package exercise

import (
	"errors"
	"testing"
)

type ExerciseRepositoryMock struct {
	exercises   []ExerciseModel
	induceError bool
}

func (exerciseRepository *ExerciseRepositoryMock) CreateExercise(name string, increment int8) (ExerciseModel, error) {
	if exerciseRepository.induceError {
		return ExerciseModel{}, errors.New("Error creating exercise")
	}
	if exerciseRepository.exercises == nil {
		exerciseRepository.exercises = make([]ExerciseModel, 0)
	}
	model := ExerciseModel{ExerciseName: name, Increment: increment}
	exerciseRepository.exercises = append(exerciseRepository.exercises, model)

	return model, nil
}
func TestServiceCreateExercise(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{})

	model, err := exerciseService.CreateExercise("Deadlift", 10)

	if err != nil {
		t.Errorf("Unexpected error - %v\n", err)
	}

	if model.ExerciseName != "Deadlift" {
		t.Errorf("Returned model name is \"%s\", expected \"Deadlift\"\n", model.ExerciseName)
	}

	if model.Increment != 10 {
		t.Errorf("Returned increment is %d, expected 10", model.Increment)
	}
}

func TestServiceCreateExerciseFails(t *testing.T) {
	exerciseService := NewExerciseService(&ExerciseRepositoryMock{induceError: true})

	_, err := exerciseService.CreateExercise("Deadlift", 10)

	if err == nil {
		t.Error("Expected an error.")
	}
}
