package iterator

import (
	"log"
	"time"

	"github.com/plotnikovanton/gomath/la"
	"github.com/plotnikovanton/skyrmions_on_go/model"
)

// SimpleSplittedIterator works like SimpleIterator, but usig vectors p and
// q instead of effective energy vector
type SimpleSplittedIterator struct {
	Times   int
	Energy  bool
	m       *model.Model
	Delta   float64
	Gp      *model.Gnuplot
	iterNum int
}

// NewSimpleSplittedIterator produces SimpleSplittedIterator by given model
// with default params:
// Times = -1
// Energy = false
// delta = 0.01
// And gnuplot instance for plotting
func NewSimpleSplittedIterator(m *model.Model) SimpleSplittedIterator {
	ret := SimpleSplittedIterator{}
	ret.Times = -1
	ret.Energy = false
	ret.m = m
	ret.Delta = 0.001
	plot := model.NewGnuplot(m)
	ret.Gp = &plot
	return ret
}

// Run runs iterator
func (iter SimpleSplittedIterator) Run() {
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

func (iter *SimpleSplittedIterator) iterate() {
	iter.iterNum++
	m := iter.m
	energy, p, q := iter.m.EnergySplitted()

	funcAdd := func(a, b la.Vector) la.Vector {
		return a.Add(b.ScalMul(iter.Delta)).Unit()
	}
	newSpins := m.Spins.WithColumn(funcAdd, p).WithColumn(funcAdd, q)

	m.Spins = newSpins
	// Plot
	if iter.Gp != nil && iter.iterNum%15 == 0 {
		iter.Gp.PlotModel()
	}
	if iter.Energy {
		log.Println("Energy: ", energy)
	}
}
