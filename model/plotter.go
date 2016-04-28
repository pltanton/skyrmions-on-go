package model

import (
	"fmt"
	"io"
	"os/exec"
)

type Gnuplot struct {
	m    *Model
	pipe *chan string
}

// Returns Gnuplot object
func NewGnuplot(m *Model) Gnuplot {
	pipe, err := getGnuplotPipe()
	if err != nil {
		panic(fmt.Sprintf("Can't create gnuplot pipe: %v", err))
	}
	// Configure gnuplot
	pipe <- fmt.Sprintf("set xrange [0:%d]\nset yrange [0:%d]\nset zrange [-1:2]", m.x*2, m.y*2)
	return Gnuplot{m, &pipe}
}

// Redraws model in current state
func (gp Gnuplot) PlotModel() {
	pipe := *gp.pipe
	m := *gp.m
	pipe <- "splot \"-\" with vectors"
	for y := 0; y < m.y; y++ {
		for x := 0; x < m.x; x++ {
			cur := m.spins[y*m.x+x]
			pipe <- fmt.Sprintf("\t%d %d %d %.3f %.3f %.3f\n", x*2, y*2, 0, cur[0], cur[1], cur[2])
		}
	}
	pipe <- "EOF"
	pipe <- "pause 0.0001"
}

func gnuplotHandler(plot chan string, pipe io.WriteCloser) {
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

func getGnuplotPipe() (chan string, error) {
	var cmd *exec.Cmd

	cmd = exec.Command("gnuplot", "-persist")

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
