package objects

import (
	"testing"
)

func TestRoll(t *testing.T) {
	dice := NewDice(6, 2, 1)
	rolledValues := dice.Roll(2)

	for _, value := range rolledValues {
		if value < 1 || value > 6 {
			t.Errorf("Invalid rolled value: %d", value)
		}
	}
}

func TestFace(t *testing.T) {
	dice := NewDice(6, 2, 1)
	faceValues := dice.Face(2)

	for _, value := range faceValues {
		if value < 1 || value > 6 {
			t.Errorf("Invalid face value: %d", value)
		}
	}
}

func TestSerialize(t *testing.T) {
	dice := NewDice(6, 2, 1)
	faceValues := dice.Serialize()

	for _, value := range faceValues {
		if value < 1 || value > 6 {
			t.Errorf("Invalid serialized value: %d", value)
		}
	}
}
