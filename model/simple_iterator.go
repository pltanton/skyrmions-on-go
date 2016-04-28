package model

import (
	"log"
	"time"

	"github.com/plotnikovanton/golinal"
)

type SimpleIterator struct {
	Times  int
	Energy bool
	m      *Model
	delta  float64
	Gp     *Gnuplot
}

func NewSimpleIterator(m *Model) SimpleIterator {
	ret := SimpleIterator{}
	ret.Times = -1
	ret.Energy = false
	ret.m = m
	ret.delta = 0.01
	plot := NewGnuplot(m)
	ret.Gp = &plot
	return ret
}

func (iter SimpleIterator) Run() {
	if iter.Times != -1 {
		start_time := time.Now()
		for i := 0; i < iter.Times; i++ {
			iter.Iterate()
		}
		log.Printf("Total time: %s", time.Since(start_time))
	} else {
		for {
			iter.Iterate()
		}
	}
}

func (iter SimpleIterator) Iterate() {
	// Update model state by adding effective energy
	energy, energy_eff := iter.m.Energy()
	m := iter.m
	dsdt := m.spins.WithColumn(
		func(a, b la.Vector) la.Vector {
			return a.ScalMul(m.gamma).CrossProd(b).Neg().Sub(a.ScalMul(m.gamma * m.lam).CrossProd(a.CrossProd(b)))
		},
		energy_eff,
	)
	new_spins := m.spins.WithColumn(
		func(a, b la.Vector) la.Vector { return a.Add(b.ScalMul(iter.delta)).Unit() },
		dsdt,
	)

	m.spins = new_spins
	// Plot
	if iter.Gp != nil {
		iter.Gp.PlotModel()
	}
	if iter.Energy {
		log.Println("Energy: ", energy)
	}
}
