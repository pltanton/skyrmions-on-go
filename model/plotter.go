package model

import (
	"fmt"
	"io"
	"os/exec"
)

// Gnuplot is object to manipulates gnuplot subprogramm
type Gnuplot struct {
	m    *Model
	pipe *chan string
}

// NewGnuplot createss gnuplot object
func NewGnuplot(m *Model) Gnuplot {
	pipe, err := getGnuplotPipe()
	if err != nil {
		panic(fmt.Sprintf("Can't create gnuplot pipe: %v", err))
	}
	// Configure gnuplot
	z := (m.x + m.y) / 4
	pipe <- fmt.Sprintf("set xrange [0:%d]\nset yrange [0:%d]\nset zrange [%d:%d]", m.x*2, m.y*2, -z, z)
	pipe <- "set view map"
	return Gnuplot{m, &pipe}
}

// PlotModel redraws model in current state
func (gp Gnuplot) PlotModel() {
	pipe := *gp.pipe
	m := *gp.m
	pipe <- "splot \"-\" notitle with vectors palette"
	pipe <- m.SpinsToString()
	pipe <- "EOF"
	pipe <- "pause 0.0001"
}

// Pipe pipes command to gnuplot
func (gp Gnuplot) Pipe(command string) {
	*gp.pipe <- command
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
