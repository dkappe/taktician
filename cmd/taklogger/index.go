package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nelhage/taktician/logs"
	"github.com/nelhage/taktician/ptn"
)

func PTNGame(g *ptn.PTN) *logs.Game {
	day := g.FindTag("Date")
	id, e := strconv.Atoi(g.FindTag("Id"))
	if day == "" || e != nil {
		return nil
	}
	size, _ := strconv.Atoi(g.FindTag("Size"))
	t, _ := time.Parse(time.RFC3339, g.FindTag("Time"))
	player1 := g.FindTag("Player1")
	player2 := g.FindTag("Player2")
	result := g.FindTag("Result")
	winner := (&ptn.Result{Result: result}).Winner().String()
	moves := countMoves(g)
	return &logs.Game{
		Day:       day,
		ID:        id,
		Timestamp: t,
		Size:      size,
		Player1:   player1,
		Player2:   player2,
		Result:    result,
		Winner:    winner,
		Moves:     moves,
	}
}

func indexPTN(repo *logs.Repository, dir string, db string) error {
	ptns, e := readPTNs(dir)
	if e != nil {
		return e
	}

	var gs []*logs.Game
	for _, g := range ptns {
		lg := PTNGame(g)
		if lg == nil {
			continue
		}
		gs = append(gs, lg)
	}
	err := repo.InsertGames(gs)
	if err != nil {
		return fmt.Errorf("insert: %v", err)
	}

	return nil
}

func countMoves(g *ptn.PTN) int {
	i := 0
	for _, o := range g.Ops {
		if _, ok := o.(*ptn.Move); ok {
			i++
		}
	}
	return i
}

func readPTNs(d string) ([]*ptn.PTN, error) {
	var out []*ptn.PTN
	e := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".ptn") {
			return nil
		}
		g, e := ptn.ParseFile(path)
		if e != nil {
			log.Printf("parse(%s): %v", path, e)
			return nil
		}
		out = append(out, g)
		return nil
	})
	return out, e
}
