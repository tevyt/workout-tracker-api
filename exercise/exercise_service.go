package exercise

type ExerciseService interface {
	CreateExercise(name string, increment int8) (ExerciseModel, error)
}

type ExerciseServiceImpl struct {
	repository ExerciseRepository
}

func NewExerciseService(repository ExerciseRepository) *ExerciseServiceImpl {
	return &ExerciseServiceImpl{
		repository: repository,
	}
}

func (excersiseService *ExerciseServiceImpl) CreateExercise(name string, increment int8) (ExerciseModel, error) {
	return excersiseService.repository.CreateExercise(name, increment)
}
