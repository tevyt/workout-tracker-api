package exercise

import "database/sql"

type ExerciseService interface {
	CreateExercise(name string, increment int8) (Exercise, error)
	GetExercise(id int64) (Exercise, error)
	SearchExercises(query string) ([]Exercise, error)
}

type ExerciseServiceImpl struct {
	repository ExerciseRepository
}

func NewExerciseService(repository ExerciseRepository) *ExerciseServiceImpl {
	return &ExerciseServiceImpl{
		repository: repository,
	}
}

func (excersiseService *ExerciseServiceImpl) CreateExercise(name string, increment int8) (Exercise, error) {
	return excersiseService.repository.CreateExercise(name, increment)
}

func (exerciseService *ExerciseServiceImpl) GetExercise(id int64) (Exercise, error) {
	exercise, err := exerciseService.repository.GetExercise(id)

	if err == sql.ErrNoRows {
		return Exercise{}, ExerciseNotFoundError{Id: id}
	}

	return exercise, err
}

func (exerciseService *ExerciseServiceImpl) SearchExercises(query string) ([]Exercise, error) {
	return exerciseService.repository.SearchExercises(query)
}
