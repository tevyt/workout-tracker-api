package exercise

import (
	"github.com/jmoiron/sqlx"
)

type ExerciseRepository interface {
	CreateExercise(name string, increment int8) (Exercise, error)
	GetExercise(id int64) (Exercise, error)
	SearchExercises(query string) ([]Exercise, error)
}

type ExerciseRepositoryImpl struct {
	databaseConnection *sqlx.DB
}

func NewExerciseRepository(databaseConnection *sqlx.DB) *ExerciseRepositoryImpl {
	return &ExerciseRepositoryImpl{
		databaseConnection: databaseConnection,
	}
}
func (excersiseRepository *ExerciseRepositoryImpl) CreateExercise(name string, increment int8) (Exercise, error) {
	createExerciseQuery := "INSERT INTO exercise (exercise_name, increment) VALUES ($1, $2) RETURNING id"
	exercise := Exercise{ExerciseName: name, Increment: increment}
	result := excersiseRepository.databaseConnection.QueryRowx(createExerciseQuery, name, increment)

	var id int64

	err := result.Scan(&id)

	if err != nil {
		return exercise, err
	}

	exercise.Id = id
	exercise.ExerciseName = name
	exercise.Increment = increment

	return exercise, nil
}

func (exerciseRepository *ExerciseRepositoryImpl) GetExercise(id int64) (Exercise, error) {
	getExerciseQuery := "SELECT id, exercise_name, increment FROM exercise WHERE id = $1"

	exercise := Exercise{}

	err := exerciseRepository.databaseConnection.Get(&exercise, getExerciseQuery, id)

	return exercise, err
}

func (exerciseRepository *ExerciseRepositoryImpl) SearchExercises(query string) ([]Exercise, error) {
	searchExerciseQuery := "SELECT id, exercise_name, increment FROM exercise WHERE exercise_name ILIKE $1"

	var results []Exercise

	err := exerciseRepository.databaseConnection.Select(&results, searchExerciseQuery, "%"+query+"%")

	if err == nil && results == nil {
		return make([]Exercise, 0), nil
	}

	return results, err
}
