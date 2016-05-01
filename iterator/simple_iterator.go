package iterator

import (
	"log"
	"time"

	"github.com/plotnikovanton/golinal"
	"github.com/plotnikovanton/skyrmions_on_go/model"
)

// SimpleIterator iterates model by adding to current models spins state
// gradient of effective energy with energy loss with fixed delta
type SimpleIterator struct {
	Times   int
	Energy  bool
	m       *model.Model
	Delta   float64
	Gp      *model.Gnuplot
	iterNum int
}

// NewSimpleIterator produces SimpleIterator by given model
// with default params:
// Times = -1
// Energy = false
// delta = 0.01
// And gnuplot instance for plotting
func NewSimpleIterator(m *model.Model) SimpleIterator {
	ret := SimpleIterator{}
	ret.Times = -1
	ret.Energy = false
	ret.m = m
	ret.Delta = 0.01
	plot := model.NewGnuplot(m)
	ret.Gp = &plot
	return ret
}

// Run runs iterator
func (iter SimpleIterator) Run() {
	if iter.Times != -1 {
		startTime := time.Now()
		for i := 0; i < iter.Times; i++ {
			iter.iterate()
		}
		log.Printf("Total time: %s", time.Since(startTime))
	} else {
		for {
			iter.iterate()
		}
	}
}

func (iter *SimpleIterator) iterate() {
	iter.iterNum++
	// Update model state by adding effective energy
	energy, energyEff := iter.m.Energy()
	m := iter.m
	dsdt := m.Spins.WithColumn(
		func(a, b la.Vector) la.Vector {
			return a.ScalMul(m.Gamma).CrossProd(b).Neg().Sub(a.ScalMul(m.Gamma * m.Lam).CrossProd(a.CrossProd(b)))
		},
		energyEff,
	)
	newSpins := m.Spins.WithColumn(
		func(a, b la.Vector) la.Vector { return a.Add(b.ScalMul(iter.Delta)).Unit() },
		dsdt,
	)

	m.Spins = newSpins
	// Plot
	if iter.Gp != nil && iter.iterNum%15 == 0 {
		iter.Gp.PlotModel()
	}
	if iter.Energy {
		log.Println("Energy: ", energy)
	}
}
