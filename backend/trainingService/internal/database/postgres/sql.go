package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

func GetExercises() (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT exercise.id, exercise.name, exercise.description, exercise.image, STRING_AGG(muscles.muscle, ',') AS muscles FROM exercise LEFT JOIN muscles ON exercise.id = muscles.id_exercise GROUP BY exercise.id, exercise.name, exercise.description, exercise.image")
}

func GetExerciseById(id int64) pgx.Row {
	return GetDB().QueryRow(context.Background(), "SELECT exercise.id, exercise.name, exercise.description, exercise.image, STRING_AGG(muscles.muscle, ',') AS muscles FROM exercise LEFT JOIN muscles ON exercise.id = muscles.id_exercise WHERE exercise.id = $1 GROUP BY exercise.id, exercise.name, exercise.description, exercise.image", id)
}

func CreateProgram(id int64, name string, description string, exercises []int64) (pgconn.CommandTag, error) {
	exercisesStr := strings.Trim(strings.Replace(fmt.Sprint(exercises), " ", ",", -1), "[]")

	query := fmt.Sprintf("WITH inserted_program AS ( INSERT INTO programs (id_creator, name, description) VALUES (%d, '%s', '%s') RETURNING id ), inserted_active AS ( INSERT INTO active (id_program, id_user) SELECT id, %d FROM inserted_program RETURNING id_program ) INSERT INTO programs_exercises (id_program, id_exercise) SELECT id_program, unnest(ARRAY[%s]) FROM inserted_active", id, name, description, id, exercisesStr)

	return GetDB().Exec(context.Background(), query)
}

func GetProgramsByUserId(id int64) (pgx.Rows, error) {
	return GetDB().Query(context.Background(), "SELECT programs.id AS program_id, programs.id_creator, programs.name, programs.description, STRING_AGG(programs_exercises.id_exercise::text, ',') AS exercises FROM active JOIN programs ON active.id_program = programs.id JOIN programs_exercises ON programs.id = programs_exercises.id_program WHERE active.id_user = $1 GROUP BY programs.id, programs.id_creator, programs.name, programs.description", id)
}
