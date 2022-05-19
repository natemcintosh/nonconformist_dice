package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Dice struct {
	frequencies [4]uint8
}

func NewRandomDice() Dice {
	// Pick three random ints [0,4]. These are the "split_points"
	var split_points [3]int

	for i := 0; i < 3; i++ {
		split_points[i] = rand.Intn(5)
	}

	// Sort the numbers
	sort.Ints(split_points[:])

	// Fill in the `dice` array
	var dice [4]uint8

	dice[0] = uint8(split_points[0])
	dice[1] = uint8(split_points[1] - split_points[0])
	dice[2] = uint8(split_points[2] - split_points[1])
	dice[3] = uint8(4 - split_points[2])

	return Dice{dice}
}

// Roll will re-roll any dice with duplicate values, eg. if you have two threes, then
// those two dice will be re-rolled.
func (d *Dice) Roll() {
	// Create a new array to populate
	var new_dice [4]uint8

	// For dice that has a frequency of more than 1, reset it to zero, and then roll again,
	// increase the frequency of the rolled value by 1.
	for i := 0; i < 4; i++ {
		n_dice := int(d.frequencies[i])
		if n_dice > 1 {
			// Get n_dice random numbers and add them to the frequencies
			for j := 0; j < n_dice; j++ {
				new_dice[rand.Intn(4)] += 1
			}
		} else if n_dice == 1 {
			new_dice[i] += 1
		}

	}
	d.frequencies = new_dice
}

func (d *Dice) sum() uint8 {
	return d.frequencies[0] + d.frequencies[1] + d.frequencies[2] + d.frequencies[3]
}

type GameResult int

const (
	Lost    GameResult = -1
	NotOver GameResult = 0
	Won     GameResult = 1
)

// GameIsOver means either all dice are unique or all are the same
func (d *Dice) GameIsOver() GameResult {
	all_ones := true
	for _, v := range d.frequencies {
		if v == 4 {
			// player lost
			return Lost
		} else if v != 1 {
			all_ones = false
		}
	}
	if all_ones {
		return Won
	}
	return NotOver
}

func Play() (won bool, n_rolls int) {
	d := NewRandomDice()
	n_rolls = 1
	for {
		d.Roll()

		result := d.GameIsOver()
		switch result {
		case Lost:
			// Player lost
			return false, n_rolls
		case Won:
			// Player won
			return true, n_rolls
		case NotOver:
			n_rolls += 1
		}
	}
}

const NGAMES = 50_000_000

func main() {
	var n_won, n_lost int

	rand.Seed(time.Hour.Milliseconds())

	// Run NGAMES
	for i := 0; i < NGAMES; i++ {
		won, _ := Play()
		if won {
			n_won += 1
		} else {
			n_lost += 1
		}
	}

	fmt.Println("win/loss ratio:")
	fmt.Printf("%d/%d\n", n_won, n_lost)
	fmt.Printf("\n\n total = %v\n", n_won+n_lost)
}
