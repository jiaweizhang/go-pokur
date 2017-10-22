package ranker

import "fmt"
import (
	"sort"
	"sync"
)

type Card struct {
	Rank int
	Suit rune
}

type Hand struct {
	Owner int64
	Cards []Card
}

type internalHand struct {
	owner    int64
	bestHand []Card
	ranking  int
}

type RankingResult struct {
	Owner    int64
	BestHand []Card
}

type chanItem struct {
	score    int
	bestHand []Card
	owner    int64
}

func ProcessShowdown(community []Card, hands ...Hand) ([][]RankingResult, error) {
	if err := verifyShowdownArguments(community, hands...); err != nil {
		return nil, err
	}

	internalScores := make([]internalHand, 0, 9)

	// channel for processed player hand and community cards
	c := make(chan chanItem)

	// used to ensure all hands are finished processing before continuing
	var wg sync.WaitGroup
	wg.Add(len(hands))

	for _, hand := range hands {
		// find the score for the Hand
		sevenCards := []Card{community[0], community[1], community[2], community[3], community[4], hand.Cards[0], hand.Cards[1]}
		go score7Chan(c, hand.Owner, sevenCards)
	}

	go func() {
		for ci := range c {
			internalScores = append(internalScores, internalHand{ci.owner, ci.bestHand, ci.score})
			wg.Done()
		}
	}()

	wg.Wait()
	close(c)

	sort.Slice(internalScores, func(i, j int) bool {
		return internalScores[i].ranking > internalScores[j].ranking
	})

	currentRanking := internalScores[0].ranking

	ordering := make([][]RankingResult, 0, 9)

	tiedForBest := make([]RankingResult, 0, 9)
	tiedForBest = append(tiedForBest, RankingResult{internalScores[0].owner, internalScores[0].bestHand})

	for i := 1; i < len(hands); i++ {
		if internalScores[i].ranking != currentRanking {
			// if we have reached a new Rank, copy the current slice into the 2D slice and create a new slice
			ordering = append(ordering, tiedForBest)
			tiedForBest = make([]RankingResult, 0, 9)
		}

		tiedForBest = append(tiedForBest, RankingResult{internalScores[i].owner, internalScores[i].bestHand})
	}

	// add last array in
	ordering = append(ordering, tiedForBest)
	return ordering, nil
}

func score7Chan(c chan chanItem, owner int64, cards []Card) {
	bestCards, max := score7(cards)
	c <- chanItem{max, bestCards, owner}
}

func score7(cards []Card) ([]Card, int) {
	max := 0
	bestCards := make([]Card, 5, 5)

	doScore := func(fiveCards []Card) {
		score := score5(fiveCards)
		if score > max {
			max = score
			bestCards = fiveCards
		}
	}

	// 1 1 1 1 1 0 0
	doScore([]Card{cards[0], cards[1], cards[2], cards[3], cards[4]})

	// 1 1 1 1 0 1 0
	doScore([]Card{cards[0], cards[1], cards[2], cards[3], cards[5]})

	// 1 1 1 1 0 0 1
	doScore([]Card{cards[0], cards[1], cards[2], cards[3], cards[6]})

	// 1 1 1 0 1 1 0
	doScore([]Card{cards[0], cards[1], cards[2], cards[4], cards[5]})

	// 1 1 1 0 1 0 1
	doScore([]Card{cards[0], cards[1], cards[2], cards[4], cards[6]})

	// 1 1 1 0 0 1 1
	doScore([]Card{cards[0], cards[1], cards[2], cards[5], cards[6]})

	// 1 1 0 1 1 1 0
	doScore([]Card{cards[0], cards[1], cards[3], cards[4], cards[5]})

	// 1 1 0 1 1 0 1
	doScore([]Card{cards[0], cards[1], cards[3], cards[4], cards[6]})

	// 1 1 0 1 0 1 1
	doScore([]Card{cards[0], cards[1], cards[3], cards[5], cards[6]})

	// 1 1 0 0 1 1 1
	doScore([]Card{cards[0], cards[1], cards[4], cards[5], cards[6]})

	// 1 0 1 1 1 1 0
	doScore([]Card{cards[0], cards[2], cards[3], cards[4], cards[5]})

	// 1 0 1 1 1 0 1
	doScore([]Card{cards[0], cards[2], cards[3], cards[4], cards[6]})

	// 1 0 1 1 0 1 1
	doScore([]Card{cards[0], cards[2], cards[3], cards[5], cards[6]})

	// 1 0 1 0 1 1 1
	doScore([]Card{cards[0], cards[2], cards[4], cards[5], cards[6]})

	// 1 0 0 1 1 1 1
	doScore([]Card{cards[0], cards[3], cards[4], cards[5], cards[6]})

	// 0 1 1 1 1 1 0
	doScore([]Card{cards[1], cards[2], cards[3], cards[4], cards[5]})

	// 0 1 1 1 1 0 1
	doScore([]Card{cards[1], cards[2], cards[3], cards[4], cards[6]})

	// 0 1 1 1 0 1 1
	doScore([]Card{cards[1], cards[2], cards[3], cards[5], cards[6]})

	// 0 1 1 0 1 1 1
	doScore([]Card{cards[1], cards[2], cards[4], cards[5], cards[6]})

	// 0 1 0 1 1 1 1
	doScore([]Card{cards[1], cards[3], cards[4], cards[5], cards[6]})

	// 0 0 1 1 1 1 1
	doScore([]Card{cards[2], cards[3], cards[4], cards[5], cards[6]})

	return bestCards, max
}

func score5(c []Card) int {
	// sort by Rank
	sort.Slice(c[:], func(i, j int) bool {
		return c[i].Rank > c[j].Rank
	})

	// check flush state
	isFlush := c[0].Suit == c[1].Suit && c[0].Suit == c[2].Suit && c[0].Suit == c[3].Suit && c[0].Suit == c[4].Suit
	isStraight := c[1].Rank == c[2].Rank+1 && c[1].Rank == c[3].Rank+2 && c[1].Rank == c[4].Rank+3 && (c[0].Rank == c[1].Rank+1 || c[4].Rank == c[0].Rank-12)

	if isStraight && isFlush {
		// royal flush
		if c[0].Rank != 12 || c[1].Rank != 3 {
			return 9<<20 + c[0].Rank<<16
		}
		// 5 high straight flush
		return 9<<20 + 3<<16
	}

	if isStraight {
		// straight
		if c[0].Rank != 12 || c[1].Rank != 3 {
			return 5<<20 + c[0].Rank<<16
		}
		// 5 high straight
		return 5<<20 + 3<<16
	}

	if isFlush {
		// flush
		return 6<<20 + c[0].Rank<<16 + c[1].Rank<<12 + c[2].Rank<<8 + c[3].Rank<<4 + c[4].Rank
	}

	// figure out groupings
	// 4 1
	// 1 4
	// 3 2
	// 2 3
	// 3 1 1
	// 1 3 1
	// 1 1 3
	// 2 2 1
	// 2 1 2
	// 1 2 2
	// 2 1 1 1
	// 1 2 1 1
	// 1 1 2 1
	// 1 1 1 2
	// 1 1 1 1 1

	if c[0].Rank != c[1].Rank {
		// 1 4
		// 1 3 1
		// 1 1 3
		// 1 2 2
		// 1 2 1 1
		// 1 1 2 1
		// 1 1 1 2
		// 1 1 1 1 1
		if c[1].Rank != c[2].Rank {
			// 1 1 3
			// 1 1 2 1
			// 1 1 1 2
			// 1 1 1 1 1
			if c[2].Rank != c[3].Rank {
				// 1 1 1 2
				// 1 1 1 1 1
				if c[3].Rank != c[4].Rank {
					// 1 1 1 1 1
					// high Card
					return 1<<20 + c[0].Rank<<16 + c[1].Rank<<12 + c[2].Rank<<8 + c[3].Rank<<4 + c[4].Rank
				}
				// 1 1 1 2
				// pair
				return 2<<20 + c[3].Rank<<16 + c[0].Rank<<12 + c[1].Rank<<8 + c[2].Rank<<4
			}

			// 1 1 3
			// 1 1 2 1
			if c[3].Rank != c[4].Rank {
				// 1 1 2 1
				// pair
				return 2<<20 + c[2].Rank<<16 + c[0].Rank<<12 + c[1].Rank<<8 + c[4].Rank<<4
			}

			// 1 1 3
			// three of a kind
			return 4<<20 + c[2].Rank<<16 + c[0].Rank<<12 + c[1].Rank<<8
		}

		// 1 4
		// 1 3 1
		// 1 2 2
		// 1 2 1 1
		if c[2].Rank != c[3].Rank {
			// 1 2 2
			// 1 2 1 1
			if c[3].Rank != c[4].Rank {
				// 1 2 1 1
				// pair
				return 2<<20 + c[1].Rank<<16 + c[0].Rank<<12 + c[3].Rank<<8 + c[4].Rank<<4
			}

			// 1 2 2
			// two pair
			return 3<<20 + c[1].Rank<<16 + c[3].Rank<<12 + c[0].Rank<<8
		}

		// 1 4
		// 1 3 1
		if c[3].Rank != c[4].Rank {
			// 1 3 1
			// three of a kind
			return 4<<20 + c[1].Rank<<16 + c[0].Rank<<12 + c[4].Rank<<8
		}

		// 1 4
		// four of a kind
		return 8<<20 + c[1].Rank<<16 + c[0].Rank<<12
	}

	// 4 1
	// 3 2
	// 3 1 1
	// 2 3
	// 2 2 1
	// 2 1 2
	// 2 1 1 1
	if c[1].Rank == c[2].Rank {
		// 4 1
		// 3 2
		// 3 1 1
		if c[3].Rank == c[4].Rank {
			// 3 2
			// full house
			return 7<<20 + c[0].Rank<<16 + c[3].Rank<<12
		}
		if c[2].Rank != c[3].Rank {
			// 3 1 1
			// three of a kind
			return 4<<20 + c[0].Rank<<16 + c[3].Rank<<12 + c[4].Rank<<8
		}

		// 4 1
		// four of a kind
		return 8<<20 + c[0].Rank<<16 + c[4].Rank<<12
	}

	// 2 3
	// 2 2 1
	// 2 1 2
	// 2 1 1 1
	if c[2].Rank != c[3].Rank {
		// 2 1 2
		// 2 1 1 1
		if c[3].Rank != c[4].Rank {
			// 2 1 1 1
			// pair
			return 2<<20 + c[0].Rank<<16 + c[2].Rank<<12 + c[3].Rank<<8 + c[4].Rank<<4
		}
		// 2 1 2
		// two pair
		return 3<<20 + c[0].Rank<<16 + c[3].Rank<<12 + c[2].Rank<<8
	}
	// 2 3
	// 2 2 1
	if c[3].Rank != c[4].Rank {
		// 2 2 1
		// two pair
		return 3<<20 + c[0].Rank<<16 + c[2].Rank<<12 + c[4].Rank<<8
	}

	// 2 3
	// full house
	return 7<<20 + c[2].Rank<<16 + c[0].Rank<<12
}

func verifyShowdownArguments(community []Card, hands ...Hand) error {
	if len(hands) < 2 {
		return fmt.Errorf("must have at least 2 players")
	} else if len(hands) > 9 {
		return fmt.Errorf("only support up to 9 players")
	}

	if len(community) != 5 {
		return fmt.Errorf("must have 5 community cards")
	}

	// TODO(jiaweizhang) Ensure that owners are unique

	// ensure there are no duplicate cards among community cards
	var mask int64
	for _, c := range community {
		if shift, err := CalculateBit(c); err != nil {
			return err
		} else if nshift := int64(1) << uint(shift); mask&nshift != 0 {
			return fmt.Errorf("duplicate community card %v", c)
		} else {
			mask += nshift
		}
	}

	// ensure there are no duplicate cards among hands
	for _, h := range hands {
		if len(h.Cards) != 2 {
			return fmt.Errorf("hands must have 2 cards each %v", h)
		}
		if shift, err := CalculateBit(h.Cards[0]); err != nil {
			return err
		} else if nshift := int64(1) << uint(shift); mask&nshift != 0 {
			return fmt.Errorf("duplicate card %v for owner %v", h.Cards[0], h.Owner)
		} else {
			mask += nshift
		}

		if shift, err := CalculateBit(h.Cards[1]); err != nil {
			return err
		} else if nshift := int64(1) << uint(shift); mask&nshift != 0 {
			return fmt.Errorf("duplicate card %v for owner %v", h.Cards[1], h.Owner)
		} else {
			mask += nshift
		}
	}

	return nil
}
