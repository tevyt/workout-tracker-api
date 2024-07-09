package exercise

import "fmt"

type ExerciseNotFoundError struct {
	Id int64
}

func (exerciseError ExerciseNotFoundError) Error() string {
	return fmt.Sprintf("No exercise found with id: %d", exerciseError.Id)
}
