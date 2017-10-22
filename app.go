package main

import "fmt"
import "github.com/jiaweizhang/goker/ranker"

func main() {
	fmt.Println("Goker app")

	aa := ranker.Hand{1, []ranker.Card{ranker.Card{12, 'H'}, ranker.Card{12, 'S'}}}
	kk := ranker.Hand{2, []ranker.Card{ranker.Card{11, 'H'}, ranker.Card{11, 'S'}}}
	qq := ranker.Hand{3, []ranker.Card{ranker.Card{10, 'H'}, ranker.Card{10, 'S'}}}
	aa2 := ranker.Hand{4, []ranker.Card{ranker.Card{12, 'D'}, ranker.Card{12, 'C'}}}

	community := []ranker.Card{
		ranker.Card{6, 'H'},
		ranker.Card{2, 'D'},
		ranker.Card{6, 'S'},
		ranker.Card{9, 'H'},
		ranker.Card{1, 'H'},
	}

	result, err := ranker.ProcessShowdown(community, aa, kk, qq, aa2)

	fmt.Println(result, err)
}
