package objects

type Dice struct {
	numberOfSides int
	numberOfDice  int
	dice          []Die
}

// NewDice creates a new Dice instance with the given number of sides, number of dice, and face
func NewDice(face int, numberOfDice int, numberOfSides int) Dice {
	dice := make([]Die, numberOfDice)
	for i := range dice {
		dice[i] = NewDie(numberOfSides, face)
	}
	return Dice{
		numberOfSides: numberOfSides,
		numberOfDice:  numberOfDice,
		dice:          dice,
	}
}

// Serialize returns the face values of all dice as a slice
func (d *Dice) Serialize() []int {
	faceValues := make([]int, len(d.dice))
	for i, die := range d.dice {
		faceValues[i] = die.Face()
	}
	return faceValues
}

// Roll simulates rolling the dice and returns the face values
func (d *Dice) Roll(numberOfDice int) []int {
	if numberOfDice > 0 {
		numberOfDice = d.numberOfDice
	}

	faceValues := make([]int, numberOfDice)
	for i := 0; i < numberOfDice; i++ {
		faceValues[i] = d.dice[i].Roll()
	}
	return faceValues
}

// Face returns the current face values of the dice
func (d *Dice) Face(numberOfDice int) []int {
	if numberOfDice > 0 {
		numberOfDice = d.numberOfDice
	}

	faceValues := make([]int, numberOfDice)
	for i := 0; i < numberOfDice; i++ {
		faceValues[i] = d.dice[i].Face()
	}
	return faceValues
}
