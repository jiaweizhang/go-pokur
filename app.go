package main

import "fmt"
import (
	//"github.com/jiaweizhang/goker/ranker"
	"flag"
	"github.com/jiaweizhang/goker/data"
	"github.com/jiaweizhang/goker/ranker"
	"github.com/jiaweizhang/goker/server"
	"log"
	"net/http"
	"os"
)

func main() {
	fsdir := flag.String("fsdir", "", "Static file directory to serve")
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

	ranker.ProcessShowdown(community, aa, kk, qq, aa2)

	log.SetOutput(os.Stdout)

	flag.Parse()

	// serve if dir is provided
	if *fsdir != "" {
		log.Printf("Starting fileserver at %s", *fsdir)
		fs := http.FileServer(http.Dir(*fsdir))
		http.Handle("/", fs)
	}

	data.Init()
	server.Init()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("ListenAndServe error: %v", err)
	}
}
