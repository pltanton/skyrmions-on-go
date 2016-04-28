package model

import (
	"log"

	"github.com/plotnikovanton/golinal"
)

type SimpleIterator struct {
	times  int
	Energy bool
	m      *Model
	delta  float64
	Plot   *Plot
}

func NewSimpleIterator(m *Model) SimpleIterator {
	ret := SimpleIterator{}
	ret.times = -1
	ret.Energy = false
	ret.m = m
	ret.delta = 0.01
	plot, _ := GetGnuplotPipe()
	ret.Plot = &plot
	return ret
}

func (iter SimpleIterator) Run() {
	iter.m.ConfigureGnuplot(iter.Plot)
	if iter.times != -1 {
		for i := 0; i < iter.times; i++ {
			iter.Iterate()
		}
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
	if iter.Plot != nil {
		m.PlotModel(iter.Plot)
	}
	if iter.Energy {
		log.Println("Energy: ", energy)
	}
}
