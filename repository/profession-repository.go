package repository

import (
	"encoding/json"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"os"
)

type ProfessionRepository interface {
	All() []entity.Profession
	FindProfessionById(ID uint64) entity.Profession
	PickProfession(excluded *[]int) entity.Profession
}

type professionConnection struct {
	path string
}

func NewProfessionRepository(path string) ProfessionRepository {
	return &professionConnection{
		path: path,
	}
}

func (db *professionConnection) All() []entity.Profession {
	data, err := os.ReadFile(db.path)
	if err != nil {
		panic(err)
	}

	var professions []entity.Profession

	err = json.Unmarshal(data, &professions)
	if err != nil {
		panic(err)
	}

	return professions
}

func (db *professionConnection) FindProfessionById(ID uint64) entity.Profession {
	professions := db.All()

	for i := 0; i < len(professions); i++ {
		if professions[i].ID == ID {
			return professions[i]
		}
	}

	return entity.Profession{}
}

func (db *professionConnection) PickProfession(excluded *[]int) entity.Profession {
	professions := db.All()

	if excluded != nil {
		for i := 0; i < len(professions); i++ {
			if helper.Contains[int](*excluded, int(professions[i].ID)) {
				professions = append(professions[:i], professions[i+1:]...)
			}
		}
	}

	return professions[helper.Random(len(professions)-1)]
}
