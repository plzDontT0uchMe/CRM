package models

import (
	"time"
)

type Record struct {
	ID        int       `json:"id"`
	Client    *User     `json:"client"`
	Trainer   *User     `json:"trainer,omitempty"`
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
}

type RecordForDay struct {
	TrainerId int64     `json:"trainerId"`
	Day       time.Time `json:"day"`
}

type Application struct {
	ID           int          `json:"id"`
	Client       *User        `json:"client,omitempty"`
	Subscription Subscription `json:"subscription"`
}

type Program struct {
	ID          int        `json:"id"`
	Creator     User       `json:"creator"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Video       string    `json:"video,omitempty"`
	Muscles     *[]string `json:"muscles,omitempty"`
}

type TrainersInfo struct {
	Exp          int    `json:"exp"`
	Sport        string `json:"sport"`
	Achievements string `json:"achievements"`
}

type Subscription struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Price          float32  `json:"price"`
	Description    string   `json:"description"`
	Possibilities  []string `json:"possibilities"`
	Trainer        *User    `json:"trainer,omitempty"`
	DateExpiration string   `json:"dateExpiration,omitempty"`
}

type ChangeApplication struct {
	Application *Application `json:"application"`
	IsAccepted  bool         `json:"isAccepted"`
}

type User struct {
	ID           int            `json:"id,omitempty"`
	LastActivity time.Time      `json:"lastActivity,omitempty"`
	DateCreated  time.Time      `json:"dateCreated,omitempty"`
	Name         string         `json:"name"`
	Surname      string         `json:"surname"`
	Patronymic   string         `json:"patronymic"`
	Gender       int            `json:"gender"`
	DateBorn     string         `json:"dateBorn"`
	Image        string         `json:"image"`
	Position     string         `json:"position"`
	TrainerInfo  []TrainersInfo `json:"trainerInfo,omitempty"`
	Subscription *Subscription  `json:"subscription,omitempty"`
	Application  *Application   `json:"application,omitempty"`
}
