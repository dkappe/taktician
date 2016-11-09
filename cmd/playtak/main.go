package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"taktician/ai"
	"taktician/ai/mcts"
	"taktician/cli"
	"taktician/ptn"
	"taktician/tak"
)

var (
	white = flag.String("white", "human", "white player")
	black = flag.String("black", "human", "white player")
	size  = flag.Int("size", 5, "game size")
	debug = flag.Int("debug", 0, "debug level")
	limit = flag.Duration("limit", time.Minute, "ai time limit")
	out   = flag.String("out", "", "write ptn to file")

	unicode = flag.Bool("unicode", false, "render board with utf8 glyphs")
)

func glyphs() *cli.Glyphs {
	if *unicode {
		return &cli.UnicodeGlyphs
	}
	return &cli.DefaultGlyphs
}

type aiWrapper struct {
	p ai.TakPlayer
}

func (a *aiWrapper) GetMove(p *tak.Position) tak.Move {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(*limit))
	defer cancel()
	return a.p.GetMove(ctx, p)
}

func parsePlayer(in *bufio.Reader, s string) cli.Player {
	if s == "human" {
		return cli.NewCLIPlayer(os.Stdout, in)
	}
	if s == "rand" {
		return &aiWrapper{ai.NewRandom(0)}
	}
	if strings.HasPrefix(s, "rand") {
		var seed int64
		if len(s) > len("rand") {
			i, err := strconv.Atoi(s[len("rand:"):])
			if err != nil {
				log.Fatal(err)
			}
			seed = int64(i)
		}
		return &aiWrapper{ai.NewRandom(seed)}
	}
	if strings.HasPrefix(s, "minimax") {
		var depth = 3
		if len(s) > len("minimax") {
			i, err := strconv.Atoi(s[len("minimax:"):])
			if err != nil {
				log.Fatal(err)
			}
			depth = i
		}
		p := ai.NewMinimax(ai.MinimaxConfig{
			Size:  *size,
			Depth: depth,
			Debug: *debug,
		})
		return &aiWrapper{p}
	}
	if strings.HasPrefix(s, "newmax") {
		var depth = 3
		if len(s) > len("newmax") {
			i, err := strconv.Atoi(s[len("newmax:"):])
			if err != nil {
				log.Fatal(err)
			}
			depth = i
		}
		p := ai.NewNewmax(ai.NewmaxConfig{
			Size:  *size,
			Depth: depth,
			Debug: *debug,
		})
		return &aiWrapper{p}
	}
	if strings.HasPrefix(s, "mcts") {
		var limit = 30 * time.Second
		if len(s) > len("mcts") {
			var err error
			limit, err = time.ParseDuration(s[len("mcts:"):])
			if err != nil {
				log.Fatal(err)
			}
		}
		p := mcts.NewMonteCarlo(mcts.MCTSConfig{
			Limit: limit,
			Debug: *debug,
			Size:  *size,
		})
		return &aiWrapper{p}
	}
	log.Fatalf("unparseable player: %s", s)
	return nil
}

func main() {
	flag.Parse()
	in := bufio.NewReader(os.Stdin)
	st := &cli.CLI{
		Config: tak.Config{Size: *size},
		Out:    os.Stdout,
		White:  parsePlayer(in, *white),
		Black:  parsePlayer(in, *black),
		Glyphs: glyphs(),
	}
	pos := st.Play()
	if *out != "" {
		p := &ptn.PTN{}
		var result string
		var winner tak.Color
		_, winner = pos.GameOver()
		if winner == tak.White {
			result = "1-0"
		} else if winner == tak.Black {
			result = "0-1"
		} else {
			result = "1/2-1/2"
		}

		p.Tags = []ptn.Tag{
			{Name: "Size", Value: strconv.Itoa(*size)},
			{Name: "Player1", Value: *white},
			{Name: "Player2", Value: *black},
			{Name: "Result", Value: result},
		}
		p.AddMoves(st.Moves())
		appendGame(*out, []byte(p.Render()))
	}
}

func appendGame(filename string, game []byte) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
         panic(err)
    }
    defer file.Close()
    if _, err = file.Write(game); err != nil {
      panic(err)
    }
    if _, err = file.WriteString("\n"); err != nil {
      panic(err)
    }
}
