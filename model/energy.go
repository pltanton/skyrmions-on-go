package model

import (
	"github.com/plotnikovanton/golinal"
)

// Energy returns pair of enegry with effective energy component
func (m Model) Energy() (float64, la.Column) {
	As := m.as()

	sas := 0.
	sb := 0.
	for i, val := range m.Spins {
		sas += val.InnerProd(As[i])
		sb += val.InnerProd(m.b)
	}

	energy := 0.5*sas + sb
	energyEff := As.WithColumn(
		func(a, b la.Vector) la.Vector { return a.Add(b) },
		m.b.NewColumn(len(m.Spins)),
	)

	return energy, energyEff
}

// As in as As_n explicity
func (m Model) as() la.Column {
	res := make(la.Column, len(m.Spins))
	for i := 0; i < len(m.Spins); i++ {
		res[i] = (m.anisotropy(i).Add(m.dzMor(i)).Add(m.atomExchange(i))).Neg()
	}
	return res
}

// DzMor returns Dzyaloshinskii-Morya component
func (m Model) dzMor(idx int) la.Vector {
	res := la.NewVector(0., 0., 0.)
	for _, bound := range m.bounds[idx] {
		otherV := m.Spins[bound.id]
		d := bound.d
		res = res.Add(otherV.CrossProd(d))
	}
	return res
}

// AtomExchange returns atom exchange component
func (m Model) atomExchange(idx int) la.Vector {
	res := la.NewVector(0., 0., 0.)
	for _, bound := range m.bounds[idx] {
		otherV := m.Spins[bound.id]
		i := bound.i
		res = res.Add(otherV.ScalMul(i))
	}
	return res
}

// Anisotropy component
func (m Model) anisotropy(idx int) la.Vector {
	v := m.Spins[idx]
	return m.k.ScalMul(2. * m.k.Len()).CrossProd(m.k.DotProd(v))
}
