package objects

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

// Die struct representing a six-sided die
type Die struct {
	numberOfSides int
	face          int
}

// NewDie creates a new Die instance with the given number of sides and face
func NewDie(numberOfSides, face int) Die {
	return Die{
		numberOfSides: numberOfSides,
		face:          face,
	}
}

// Roll simulates rolling the die and returns the face value
func (d *Die) Roll() int {
	// Seed the random number generator to produce different results each time
	num, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(999999999999))
	rand.Seed(time.Now().UnixNano() + num.Int64())

	// Simulate rolling the die
	d.face = rand.Intn(d.numberOfSides) + 1
	d.face = rand.Intn(d.numberOfSides) + 1
	d.face = rand.Intn(d.numberOfSides) + 1

	fmt.Println("ROLLING", d.face, d.numberOfSides, num.Int64())
	d.face = 4

	return d.face
}

// Face returns the current face value of the die
func (d *Die) Face() int {
	return d.face
}
