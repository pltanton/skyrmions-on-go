package model

import (
	"fmt"

	"github.com/plotnikovanton/golinal"
)

type Bound struct {
	id int
	d  la.Vector
	i  float64
}

type Model struct {
	spins  la.Column
	bounds [][]Bound
	mu     float64
	x      int
	y      int
	b      la.Vector
	k      la.Vector
}

func (m Model) String() string {
	return fmt.Sprintf(
		"Model properties is:\nSize of model: x=%v y=%v\nMU=%v\nB=%v\nK=%v",
		m.x, m.y, m.mu, m.b, m.k,
	)
}
