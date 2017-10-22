package ranker

import "fmt"

func CalculateBit(c Card) (int, error) {
	if c.Rank > 12 || c.Rank < 0 {
		return -1, fmt.Errorf("card rank is invalid for card %v", c)
	}

	if c.Suit != 'S' && c.Suit != 'H' && c.Suit != 'C' && c.Suit != 'D' {
		return -1, fmt.Errorf("card suit is invalid for card %v", c)
	}

	multiplier := 0
	if c.Suit == 'H' {
		multiplier = 1
	} else if c.Suit == 'C' {
		multiplier = 2
	} else if c.Suit == 'D' {
		multiplier = 3
	}

	return multiplier*13 + c.Rank, nil
}
