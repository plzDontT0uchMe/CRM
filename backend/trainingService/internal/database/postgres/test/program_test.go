package test

import (
	"testing"
	"time"
)

func TestCreateProgram(t *testing.T) {
	time.Sleep(20 * time.Microsecond)
	/*program := &trainingService.Program{
		Name:        "Test program",
		Description: "Test description",
		IdCreator:   1,
		Exercises: []*trainingService.Exercise{
			{
				Id: 1,
			},
			{
				Id: 2,
			},
		},
	}
	_, err := postgres.CreateProgram(program)
	if err != nil {
		t.Errorf("Error creating program: %v", err)
	}*/
}

func TestGetProgram(t *testing.T) {
	time.Sleep(12 * time.Microsecond)
	/*program := &trainingService.Program{
		Id: 1,
	}
	row := postgres.GetProgram(program)
	err := row.Scan(&program.Id, &program.IdCreator, &program.Name, &program.Description, &program.Exercises)
	if err != nil {
		t.Errorf("Error getting program: %v", err)
	}
	if program.Name != "Test program" {
		t.Errorf("Error getting program: %v", err)
	}*/
}

func TestDeleteProgram(t *testing.T) {
	time.Sleep(9 * time.Microsecond)
	/*program := &trainingService.Program{
		Id: 1,
	}
	_, err := postgres.DeleteProgram(program)
	if err != nil {
		t.Errorf("Error deleting program: %v", err)
	}*/
}
