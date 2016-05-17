package model

import (
	"bytes"
	"fmt"

	"github.com/plotnikovanton/gomath/la"
)

// Bound is used for describe bound above spins
type Bound struct {
	id int
	d  la.Vector
	i  float64
}

// Model is basic model instance to describe everything
type Model struct {
	Spins  la.Column
	bounds [][]Bound
	mu     float64
	Lam    float64
	Gamma  float64
	x      int
	y      int
	b      la.Vector
	k      la.Vector
}

// String returns string representative of model parameters
func (m Model) String() string {
	return fmt.Sprintf(
		"Model properties is:\nSize of model: x=%v y=%v\nMU=%v\nB=%v\nK=%v\nlam=%v\ngamma=%v",
		m.x, m.y, m.mu, m.b, m.k, m.Lam, m.Gamma,
	)
}

// SpinsToString returns string representative in form acceptable to gnuplot
func (m Model) SpinsToString() string {
	var buff bytes.Buffer
	for x := 0; x < m.x; x++ {
		for y := 0; y < m.y; y++ {
			s := m.Spins[x*m.y+y]
			buff.WriteString(fmt.Sprintf(
				"\t%d %d %d %.3f %.3f %.3f\n",
				x*2, y*2, 0, s[0], s[1], s[2],
			))
		}
	}
	return buff.String()
}
