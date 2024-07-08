package exercise

import (
	"github.com/jmoiron/sqlx"
)

type ExerciseRepository interface {
	CreateExercise(name string, increment int8) (ExerciseModel, error)
}

type ExerciseRepositoryImpl struct {
	databaseConnection *sqlx.DB
}

func NewExerciseRepository(databaseConnection *sqlx.DB) *ExerciseRepositoryImpl {
	return &ExerciseRepositoryImpl{
		databaseConnection: databaseConnection,
	}
}
func (excersiseRepository *ExerciseRepositoryImpl) CreateExercise(name string, increment int8) (ExerciseModel, error) {
	createExerciseQuery := "INSERT INTO exercise (exercise_name, increment) VALUES ($1, $2) RETURNING id"
	model := ExerciseModel{ExerciseName: name, Increment: increment}
	result := excersiseRepository.databaseConnection.QueryRowx(createExerciseQuery, name, increment)

	var id int64

	err := result.Scan(&id)

	if err != nil {
		return model, err
	}

	model.Id = id
	model.ExerciseName = name
	model.Increment = increment

	return model, nil
}
