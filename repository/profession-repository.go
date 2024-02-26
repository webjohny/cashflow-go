package repository

import (
	"encoding/json"
	"github.com/webjohny/cashflow-go/entity"
	"os"
)

type ProfessionRepository interface {
	All() []entity.Profession
	FindProfessionById(ID uint64) *entity.Profession
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

func (db *professionConnection) FindProfessionById(ID uint64) *entity.Profession {
	professions := db.All()

	for i := 0; i < len(professions); i++ {
		if professions[i].ID == ID {
			return &professions[i]
		}
	}

	return nil
}
