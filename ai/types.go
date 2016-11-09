package ai

import (
	"taktician/tak"
	"golang.org/x/net/context"
)

type TakPlayer interface {
	GetMove(ctx context.Context, p *tak.Position) tak.Move
}
