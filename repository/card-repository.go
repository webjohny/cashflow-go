package repository

import (
	"encoding/json"
	"github.com/webjohny/cashflow-go/entity"
	"os"
)

type CardRepository interface {
	All() []entity.Card
	FindCardById(ID string) *entity.Card
}

type cardConnection struct {
	path string
}

func NewCardRepository(path string) CardRepository {
	return &cardConnection{
		path: path,
	}
}

func (db *cardConnection) All() []entity.Card {
	data, err := os.ReadFile(db.path)
	if err != nil {
		panic(err)
	}

	var cards []entity.Card

	err = json.Unmarshal(data, &cards)
	if err != nil {
		panic(err)
	}

	return cards
}

func (db *cardConnection) FindCardById(ID string) *entity.Card {
	cards := db.All()

	for i := 0; i < len(cards); i++ {
		if cards[i].ID == ID {
			return &cards[i]
		}
	}

	return nil
}
