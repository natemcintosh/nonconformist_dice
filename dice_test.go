package main

import (
	"math/rand"
	"testing"
)

func FuzzRoll(f *testing.F) {
	// Example seeds
	f.Add(1)
	f.Add(-1)

	f.Fuzz(func(t *testing.T, seed int) {
		// Set the random seed
		rand.Seed(int64(seed))

		// Create a set of four dice
		d := NewRandomDice()

		// Make sure the values of the frequencies array sum to 4
		if d.sum() != 4 {
			t.Fatalf("Dice initialization was bad. Didn't end up with four dice")
		}

		// Roll the dice 100 times
		// Each time, check that there are 4 dice, and that the unique values
		// are not less than they were before
		for roll_count := 0; roll_count < 100; roll_count++ {
			// Get the locations of the unique values
			unique_locs := make([]int, 0, 4)
			for value, n_dice := range d.frequencies {
				if n_dice == 1 {
					unique_locs = append(unique_locs, value)
				}
			}

			d.Roll()

			if d.sum() != 4 {
				t.Errorf(
					"Didn't end up with four dice after rolling. Have %v",
					d.frequencies,
				)
			}

			// Check that the number of dice at the `unique_locs` has not decreased
			for _, unique_loc := range unique_locs {
				if d.frequencies[unique_loc] < 1 {
					t.Errorf("The number of unique dice decreased. Should never happen")
				}
			}

		}
	})
}

func Benchmark_sum(b *testing.B) {
	d := Dice{[4]uint8{1, 2, 0, 1}}
	for i := 0; i < b.N; i++ {
		d.sum()
	}
}

func BenchmarkNewRandomDice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandomDice()
	}
}

func BenchmarkRoll(b *testing.B) {
	d := Dice{[4]uint8{1, 2, 0, 1}}
	for i := 0; i < b.N; i++ {
		d.Roll()
	}
}

func BenchmarkGameIsOver(b *testing.B) {
	benchCases := []struct {
		desc string
		dice Dice
	}{
		{
			desc: "won",
			dice: Dice{[4]uint8{1, 1, 1, 1}},
		},
		{
			desc: "lost_easy",
			dice: Dice{[4]uint8{4, 0, 0, 0}},
		},
		{
			desc: "lost_hard",
			dice: Dice{[4]uint8{0, 0, 0, 4}},
		},
		{
			desc: "not_over",
			dice: Dice{[4]uint8{2, 1, 0, 1}},
		},
	}
	for _, bC := range benchCases {
		b.Run(bC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bC.dice.GameIsOver()
			}
		})
	}
}

func BenchmarkPlay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Play()
	}
}
