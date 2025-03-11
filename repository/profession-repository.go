package repository

import (
	"encoding/json"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/storage"
	"gopkg.in/errgo.v2/errors"
	"os"
)

type ProfessionRepository interface {
	All(language string) (error, []entity.Profession)
	FindProfessionById(ID uint64, language string) (error, entity.Profession)
	PickProfession(language string, excluded *[]int) (error, entity.Profession)
	SetProfessions(professions dto.ProfessionsSetBodyDTO)
}

type professionConnection struct {
	path string

	professions map[string]map[string][]entity.Profession
}

func NewProfessionRepository(path string) ProfessionRepository {
	return &professionConnection{
		path:        path,
		professions: map[string]map[string][]entity.Profession{},
	}
}

func (db *professionConnection) AllFromMemory(language string) (error, []entity.Profession) {
	professions := db.professions["default"][language]

	if professions == nil {
		return errors.New(storage.ErrorCardsNotFound), make([]entity.Profession, 0)
	}

	return nil, professions
}

func (db *professionConnection) All(language string) (error, []entity.Profession) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		panic(err)
	}

	var professions []entity.Profession

	err = json.Unmarshal(data, &professions)
	if err != nil {
		panic(err)
	}

	return nil, professions
}

func (db *professionConnection) FindProfessionById(ID uint64, language string) (error, entity.Profession) {
	err, professions := db.All("")

	if err != nil {
		return err, entity.Profession{}
	}

	for i := 0; i < len(professions); i++ {
		if professions[i].ID == ID {
			return nil, professions[i]
		}
	}

	return nil, entity.Profession{}
}

func (db *professionConnection) PickProfession(language string, excluded *[]int) (error, entity.Profession) {
	err, professions := db.All("")

	if err != nil {
		return err, entity.Profession{}
	}

	if excluded != nil {
		for i := 0; i < len(professions); i++ {
			if helper.Contains[int](*excluded, int(professions[i].ID)) {
				professions = append(professions[:i], professions[i+1:]...)
			}
		}
	}

	return nil, professions[helper.Random(len(professions)-1)]
}

func (db *professionConnection) SetProfessions(professions dto.ProfessionsSetBodyDTO) {
	if db.professions[professions.Type] == nil {
		db.professions[professions.Type] = make(map[string][]entity.Profession)
	}
	db.professions[professions.Type][professions.Language] = professions.Professions
}
