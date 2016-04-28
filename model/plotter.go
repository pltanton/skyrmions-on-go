package model

import (
	"fmt"
	"io"
	"os/exec"
)

type Plot chan string

func (m Model) ConfigureGnuplot(plot *Plot) {
	*plot <- fmt.Sprintf("set xrange [0:%d]\nset yrange [0:%d]\nset zrange [-1:2]", m.x*2, m.y*2)
}

func (m Model) PlotModel(plot *Plot) {
	*plot <- "splot \"-\" with vectors"
	for y := 0; y < m.y; y++ {
		for x := 0; x < m.x; x++ {
			cur := m.spins[y*m.x+x]
			*plot <- fmt.Sprintf("\t%d %d %d %.3f %.3f %.3f\n", x*2, y*2, 0, cur[0], cur[1], cur[2])
		}
	}
	*plot <- "EOF"
	*plot <- "pause 0.0001"
}

func gnuplotHandler(plot Plot, pipe io.WriteCloser) {
	var command string
	ok := true

	for ok {
		command, ok = <-plot

		if !ok {
			pipe.Close()
			continue
		}
		fmt.Fprintln(pipe, command)
	}
}

func GetGnuplotPipe() (Plot, error) {
	var cmd *exec.Cmd

	cmd = exec.Command("gnuplot")

	pipe, err := cmd.StdinPipe()

	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	p := make(chan string)
	go gnuplotHandler(p, pipe)

	return p, nil
}
