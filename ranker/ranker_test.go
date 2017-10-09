package ranker

import (
	"fmt"
	"testing"
)

var customTests = []struct {
	cards    [7]Card // input
	expected int     // expected result
}{
	// SF Ace High
	{g7("AS", "KS", "QS", "JS", "TS", "4S", "7H"), 9<<20 + 12<<16},
	// SF King HIgh with Straight Ace High
	{g7("AH", "KS", "QS", "JS", "TS", "9S", "7H"), 9<<20 + 11<<16},
	// SF Seven High
	{g7("4S", "5S", "6S", "7S", "3S", "2S", "8H"), 9<<20 + 5<<16},
	// SF Six High
	{g7("4S", "5S", "6S", "3S", "6H", "2S", "7H"), 9<<20 + 4<<16},
	// SF Five High
	{g7("4S", "5S", "3D", "3S", "6H", "2S", "AS"), 9<<20 + 3<<16},
	// Q 4 3
	{g7("AS", "AD", "AH", "AC", "TS", "TH", "TD"), 8<<20 + 12<<16 + 8<<12},
	// Q 3 4
	{g7("AS", "AD", "AH", "TC", "TS", "TH", "TD"), 8<<20 + 8<<16 + 12<<12},
	// Q 4 2 1
	{g7("AS", "AD", "AH", "AC", "TS", "TH", "9D"), 8<<20 + 12<<16 + 8<<12},
	// Q 4 1 2
	{g7("AS", "AD", "AH", "AC", "TS", "TH", "QD"), 8<<20 + 12<<16 + 10<<12},
	// Q 2 4 1
	{g7("AS", "AD", "TD", "TC", "TS", "TH", "9D"), 8<<20 + 8<<16 + 12<<12},
	// Q 2 1 4
	{g7("AS", "AD", "TD", "TC", "TS", "TH", "JD"), 8<<20 + 8<<16 + 12<<12},
	// Q 1 4 2
	{g7("AS", "9D", "TD", "TC", "TS", "TH", "9H"), 8<<20 + 8<<16 + 12<<12},
	// Q 1 2 4
	{g7("AS", "QD", "TD", "TC", "TS", "TH", "QC"), 8<<20 + 8<<16 + 12<<12},
	// Q 4 1 1 1
	{g7("AS", "AD", "AH", "AC", "TS", "JH", "QC"), 8<<20 + 12<<16 + 10<<12},
	// Q 1 4 1 1
	{g7("KS", "KD", "KH", "KC", "AS", "JH", "QC"), 8<<20 + 11<<16 + 12<<12},
	// Q 1 1 4 1
	{g7("JS", "JD", "JH", "JC", "9S", "QH", "KC"), 8<<20 + 9<<16 + 11<<12},
	// Q 1 1 1 4
	{g7("JS", "JD", "JH", "JC", "AS", "QH", "KC"), 8<<20 + 9<<16 + 12<<12},
	// FH 3 3 1
	{g7("JS", "JD", "JH", "TC", "TS", "TH", "9C"), 7<<20 + 9<<16 + 8<<12},
	// FH 3 1 3
	{g7("JS", "JD", "JH", "4C", "4S", "4H", "9C"), 7<<20 + 9<<16 + 2<<12},
	// FH 1 3 3
	{g7("JS", "JD", "JH", "TC", "TS", "TH", "KC"), 7<<20 + 9<<16 + 8<<12},
	// FH 3 2 2
	{g7("JS", "JD", "JH", "TC", "TS", "9H", "9C"), 7<<20 + 9<<16 + 8<<12},
	// FH 2 3 2
	{g7("JS", "JD", "JH", "TC", "TS", "AH", "AC"), 7<<20 + 9<<16 + 12<<12},
	// FH 2 2 3
	{g7("JS", "JD", "9D", "TC", "TS", "9H", "9C"), 7<<20 + 7<<16 + 9<<12},
	// FH 3 2 1 1
	{g7("JS", "JD", "JC", "TC", "TS", "9H", "8C"), 7<<20 + 9<<16 + 8<<12},
	// FH 3 1 2 1
	{g7("JS", "JD", "JC", "TC", "9S", "9H", "8C"), 7<<20 + 9<<16 + 7<<12},
	// FH 3 1 1 2
	{g7("JS", "JD", "JC", "TC", "9S", "8H", "8C"), 7<<20 + 9<<16 + 6<<12},
	// FH 2 3 1 1
	{g7("JS", "JD", "TD", "TC", "TS", "9H", "8c"), 7<<20 + 8<<16 + 9<<12},
	// FH 2 1 3 1
	{g7("JS", "JD", "TC", "9C", "9S", "9H", "8C"), 7<<20 + 7<<16 + 9<<12},
	// FH 2 1 1 3
	{g7("JS", "JD", "TC", "9C", "7S", "7H", "7C"), 7<<20 + 5<<16 + 9<<12},
	// FH 1 3 2 1
	{g7("JS", "TD", "TC", "TH", "9S", "9H", "8C"), 7<<20 + 8<<16 + 7<<12},
	// FH 1 3 1 2
	{g7("JS", "TD", "TC", "TH", "8S", "7H", "7C"), 7<<20 + 8<<16 + 5<<12},
	// FH 1 2 3 1
	{g7("JS", "TD", "TC", "9C", "9S", "9H", "8C"), 7<<20 + 7<<16 + 8<<12},
	// FH 1 2 1 3
	{g7("JS", "TD", "TC", "9C", "8S", "8H", "8C"), 7<<20 + 6<<16 + 8<<12},
	// FH 1 1 3 2
	{g7("QS", "JD", "TH", "TC", "TS", "9H", "9C"), 7<<20 + 8<<16 + 7<<12},
	// FH 1 1 2 3
	{g7("QS", "JD", "TH", "TC", "9S", "9H", "9C"), 7<<20 + 7<<16 + 8<<12},
	// F A K Q J 9
	{g7("AS", "KS", "QS", "JS", "9S", "9H", "9C"), 6<<20 + 12<<16 + 11<<12 + 10<<8 + 9<<4 + 7},
	// F 10 8 6 4 2
	{g7("TS", "8S", "6S", "4S", "2S", "9H", "9C"), 6<<20 + 8<<16 + 6<<12 + 4<<8 + 2<<4 + 0},
	// F 10 8 6 4 3
	{g7("TS", "8S", "6S", "4S", "3S", "2S", "9C"), 6<<20 + 8<<16 + 6<<12 + 4<<8 + 2<<4 + 1},
	// F 10 8 7 5 4
	{g7("TS", "8S", "7S", "5S", "4S", "3S", "2S"), 6<<20 + 8<<16 + 6<<12 + 5<<8 + 3<<4 + 2},
	// F 10 8 7 5 4
	{g7("TS", "8S", "7S", "5S", "4S", "6H", "9C"), 6<<20 + 8<<16 + 6<<12 + 5<<8 + 3<<4 + 2},
	// S A
	{g7("AS", "KS", "QD", "JS", "TS", "3D", "2C"), 5<<20 + 12<<16},
	// S A with more lower straight
	{g7("AS", "KS", "QD", "JS", "TS", "9D", "8C"), 5<<20 + 12<<16},
	// S 5
	{g7("5S", "4S", "3D", "2S", "AS", "3C", "2C"), 5<<20 + 3<<16},
	// T 3 1 1 1
	{g7("TS", "TH", "TD", "9S", "8S", "6H", "4C"), 4<<20 + 8<<16 + 7<<12 + 6<<8},
	// T 1 3 1 1
	{g7("TS", "TH", "TD", "9S", "8S", "6H", "QC"), 4<<20 + 8<<16 + 10<<12 + 7<<8},
	// T 1 1 3 1
	{g7("TS", "TH", "TD", "QS", "JS", "6H", "4C"), 4<<20 + 8<<16 + 10<<12 + 9<<8},
	// T 1 1 1 3
	{g7("4S", "4H", "4D", "KS", "QS", "JH", "AC"), 4<<20 + 2<<16 + 12<<12 + 11<<8},
	// TP 2 2 2 1
	{g7("TS", "TH", "8D", "8S", "6S", "6H", "4C"), 3<<20 + 8<<16 + 6<<12 + 4<<8},
	// TP 2 2 1 2
	{g7("TS", "TH", "8D", "8S", "6S", "4C", "4C"), 3<<20 + 8<<16 + 6<<12 + 4<<8},
	// TP 2 1 2 2
	{g7("TS", "TH", "8D", "6S", "6S", "4C", "4C"), 3<<20 + 8<<16 + 4<<12 + 6<<8},
	// TP 1 2 2 2
	{g7("TS", "8H", "8D", "6S", "6S", "4C", "4C"), 3<<20 + 6<<16 + 4<<12 + 8<<8},
	// TP 2 2 1 1 1
	{g7("TS", "TH", "8D", "8S", "6S", "5H", "4C"), 3<<20 + 8<<16 + 6<<12 + 4<<8},
	// TP 2 1 2 1 1
	{g7("TS", "TH", "8D", "8S", "6S", "9H", "4C"), 3<<20 + 8<<16 + 6<<12 + 7<<8},
	// TP 2 1 1 2 1
	{g7("TS", "TH", "8D", "6D", "6S", "9H", "4C"), 3<<20 + 8<<16 + 4<<12 + 7<<8},
	// TP 2 1 1 1 2
	{g7("TS", "TH", "8D", "7S", "6S", "4H", "4C"), 3<<20 + 8<<16 + 2<<12 + 6<<8},
	// TP 1 2 2 1 1
	{g7("TS", "TH", "8D", "8S", "QS", "5H", "4C"), 3<<20 + 8<<16 + 6<<12 + 10<<8},
	// TP 1 2 1 2 1
	{g7("TS", "TH", "8D", "8S", "QS", "9H", "4C"), 3<<20 + 8<<16 + 6<<12 + 10<<8},
	// TP 1 2 1 1 2
	{g7("TS", "TH", "8D", "4S", "QS", "9H", "4C"), 3<<20 + 8<<16 + 2<<12 + 10<<8},
	// TP 1 1 2 2 1
	{g7("TS", "TH", "8D", "8S", "QS", "KH", "4C"), 3<<20 + 8<<16 + 6<<12 + 11<<8},
	// TP 1 1 2 1 2
	{g7("TS", "TH", "8D", "8S", "QS", "KH", "9C"), 3<<20 + 8<<16 + 6<<12 + 11<<8},
	// TP 1 1 1 2 2
	{g7("TS", "AH", "8D", "8S", "4S", "KH", "4C"), 3<<20 + 6<<16 + 2<<12 + 12<<8},
	// P 2 1 1 1 1 1
	{g7("TS", "TH", "8D", "7S", "6S", "5H", "3C"), 2<<20 + 8<<16 + 6<<12 + 5<<8 + 4<<4},
	// P 1 2 1 1 1 1
	{g7("TS", "TH", "JD", "7S", "6S", "5H", "3C"), 2<<20 + 8<<16 + 9<<12 + 5<<8 + 4<<4},
	// P 1 1 2 1 1 1
	{g7("TS", "TH", "JD", "KS", "6S", "5H", "3C"), 2<<20 + 8<<16 + 11<<12 + 9<<8 + 4<<4},
	// P 1 1 1 2 1 1
	{g7("TS", "TH", "JD", "KS", "AS", "5H", "3C"), 2<<20 + 8<<16 + 12<<12 + 11<<8 + 9<<4},
	// P 1 1 1 1 2 1
	{g7("8S", "8H", "JD", "KS", "AS", "TH", "3C"), 2<<20 + 6<<16 + 12<<12 + 11<<8 + 9<<4},
	// P 1 1 1 1 1 2
	{g7("6S", "6H", "JD", "KS", "AS", "TH", "9C"), 2<<20 + 4<<16 + 12<<12 + 11<<8 + 9<<4},
	// H 1 1 1 1 1 1 1
	{g7("6S", "7H", "JD", "KS", "AS", "TH", "9C"), 1<<20 + 12<<16 + 11<<12 + 9<<8 + 8<<4 + 7},
	// H 1 1 1 1 1 1 1 low
	{g7("6S", "7H", "JD", "2S", "4S", "TH", "9C"), 1<<20 + 9<<16 + 8<<12 + 7<<8 + 5<<4 + 4},
}

func TestCustom(t *testing.T) {
	for _, tt := range customTests {
		_, actual := score7(tt.cards)
		if actual != tt.expected {
			t.Errorf("Cards: %s, expected %s, actual %s", formatCards(tt.cards), formatRank(tt.expected), formatRank(actual))
		}
	}
}

func formatCards(cards [7]Card) string {
	return formatCard(cards[0]) + " " + formatCard(cards[1]) + " " + formatCard(cards[2]) + " " + formatCard(cards[3]) + " " +
		formatCard(cards[4]) + " " + formatCard(cards[5]) + " " + formatCard(cards[6])
}

func formatCard(card Card) string {
	var rank rune
	switch card.Rank {
	case 0:
		rank = '2'
	case 1:
		rank = '3'
	case 2:
		rank = '4'
	case 3:
		rank = '5'
	case 4:
		rank = '6'
	case 5:
		rank = '7'
	case 6:
		rank = '8'
	case 7:
		rank = '9'
	case 8:
		rank = 'T'
	case 9:
		rank = 'J'
	case 10:
		rank = 'Q'
	case 11:
		rank = 'K'
	case 12:
		rank = 'A'
	}

	return string(rank) + string(card.Suit)
}

func formatRank(rank int) string {
	return fmt.Sprintf("%2d %2d %2d %2d %2d %2d", rank>>20, rank>>16&15, rank>>12&15, rank>>8&15, rank>>4&15, rank&15)
}

func generate7Cards(a, b, c, d, e, f, g string) [7]Card {
	return [7]Card{
		generateCard(a),
		generateCard(b),
		generateCard(c),
		generateCard(d),
		generateCard(e),
		generateCard(f),
		generateCard(g),
	}
}

func g7(a, b, c, d, e, f, g string) [7]Card {
	return generate7Cards(a, b, c, d, e, f, g)
}

func generateCard(card string) Card {
	runified := []rune(card)

	var rank int
	switch runified[0] {
	case 'A':
		rank = 12
	case 'K':
		rank = 11
	case 'Q':
		rank = 10
	case 'J':
		rank = 9
	case 'T':
		rank = 8
	case '9':
		rank = 7
	case '8':
		rank = 6
	case '7':
		rank = 5
	case '6':
		rank = 4
	case '5':
		rank = 3
	case '4':
		rank = 2
	case '3':
		rank = 1
	case '2':
		rank = 0
	}

	return Card{rank, runified[1]}
}
