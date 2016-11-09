package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"taktician/ptn"
	"taktician/tak"
)

func NewCLIPlayer(out io.Writer, in *bufio.Reader) Player {
	return &cliPlayer{out, in}
}

type cliPlayer struct {
	out io.Writer
	in  *bufio.Reader
}

func (c *cliPlayer) GetMove(p *tak.Position) tak.Move {
	for {
		fmt.Fprintf(c.out, "%s> ", p.ToMove())
		line, err := c.in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		line = strings.TrimRight(line, "\r\n")
		m, err := ptn.ParseMove(line)
		if err != nil {
			fmt.Fprintln(c.out, "parse error: ", err)
			continue
		}
		return m
	}
}
