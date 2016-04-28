package ai

import (
	"bytes"
	"log"

	"nelhage.com/tak/ptn"
	"nelhage.com/tak/tak"
)

const (
	maxEval int64 = 1 << 30
	minEval       = -maxEval
)

type MinimaxAI struct {
	depth int

	Debug bool
}

func formatpv(ms []tak.Move) string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, m := range ms {
		if i != 0 {
			out.WriteString(" ")
		}
		out.WriteString(ptn.FormatMove(&m))
	}
	out.WriteString("]")
	return out.String()
}

func (m *MinimaxAI) GetMove(p *tak.Position) tak.Move {
	ms, _ := m.Analyze(p)
	return ms[0]
}

func (m *MinimaxAI) Analyze(p *tak.Position) ([]tak.Move, int64) {
	var ms []tak.Move
	var v int64
	for i := 1; i <= m.depth; i++ {
		ms, v = m.minimax(p, i, ms, minEval-1, maxEval+1)
		if m.Debug {
			log.Printf("[minimax] depth=%d val=%d pv=%s",
				i, v, formatpv(ms))
		}
	}
	return ms, v
}

func (ai *MinimaxAI) minimax(
	p *tak.Position,
	depth int,
	pv []tak.Move,
	α, β int64) ([]tak.Move, int64) {
	over, _ := p.GameOver()
	if depth == 0 || over {
		return nil, ai.evaluate(p)
	}
	moves := p.AllMoves()
	if len(pv) > 0 {
		for i, m := range moves {
			if m.Equal(&pv[0]) {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	best := make([]tak.Move, 1, depth)
	max := minEval - 1
	for _, m := range moves {
		child, e := p.Move(&m)
		if e != nil {
			continue
		}
		ms, v := ai.minimax(child, depth-1, nil, -β, -α)
		v = -v
		if v > max {
			max = v
			best[0] = m
			best = append(best[:1], ms...)
		}
		if v > α {
			α = v
			if α > β {
				break
			}
		}
	}
	return best, max
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *MinimaxAI) evaluate(p *tak.Position) int64 {
	if over, winner := p.GameOver(); over {
		switch winner {
		case tak.NoColor:
			return 0
		case p.ToMove():
			return maxEval - int64(p.MoveNumber())
		default:
			return minEval + int64(p.MoveNumber())
		}
	}
	me, them := 0, 0
	for x := 0; x < p.Size(); x++ {
		for y := 0; y < p.Size(); y++ {
			sq := p.At(x, y)
			if len(sq) == 0 {
				continue
			}
			val := 0
			val += imin(x, p.Size()-x-1)
			val += imin(y, p.Size()-y-1)
			if sq[0].Kind() == tak.Flat {
				if sq[0].Color() == p.ToMove() {
					me += val
				} else {
					them += val
				}
			}
		}
	}
	return int64(me - them)
}

func NewMinimax(depth int) *MinimaxAI {
	return &MinimaxAI{depth: depth}
}
