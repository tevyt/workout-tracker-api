package exercise

import (
	_ "github.com/lib/pq"
)

type ExerciseModel struct {
	Id           int64  `db:"id" json:"id"`
	ExerciseName string `db:"exercise_name" json:"exerciseName"`
	Increment    int8   `db:"increment" json:"increment"`
}
