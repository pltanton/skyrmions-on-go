package model

import (
	"github.com/plotnikovanton/golinal"
)

// Returns enegry and effective energy
func (m Model) Energy() (float64, la.Column) {
	As := m.As()

	sas := 0.
	sb := 0.
	for i, val := range m.spins {
		sas += val.InnerProd(As[i])
		sb += val.InnerProd(m.b)
	}

	energy := 0.5*sas + sb
	energy_eff := As.WithColumn(
		func(a, b la.Vector) la.Vector { return a.Add(b) },
		m.b.NewColumn(len(m.spins)),
	)

	return energy, energy_eff
}

// A explicity
func (m Model) As() la.Column {
	res := make(la.Column, len(m.spins))
	for i := 0; i < len(m.spins); i++ {
		res[i] = (m.Anisotropy(i).Add(m.DzMor(i)).Add(m.AtomExchange(i))).Neg()
	}
	return res
}

// Dzyaloshinskii-Morya component
func (m Model) DzMor(idx int) la.Vector {
	res := la.NewVector(0., 0., 0.)
	for _, bound := range m.bounds[idx] {
		other_v := m.spins[bound.id]
		d := bound.d
		res = res.Add(other_v.CrossProd(d))
	}
	return res
}

// Atom exchange component
func (m Model) AtomExchange(idx int) la.Vector {
	res := la.NewVector(0., 0., 0.)
	for _, bound := range m.bounds[idx] {
		other_v := m.spins[bound.id]
		i := bound.i
		res = res.Add(other_v.ScalMul(i))
	}
	return res
}

// Anisotropy component
func (m Model) Anisotropy(idx int) la.Vector {
	v := m.spins[idx]
	return m.k.ScalMul(2. * m.k.Len()).CrossProd(m.k.DotProd(v))
}
