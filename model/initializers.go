package model

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/plotnikovanton/golinal"
	"gopkg.in/yaml.v2"
)

// NewBasicModel creates simple 2D grid model form yaml file
// All spins would be oriented randomly at start
// Crystalic field would be formed as 4 neigbours system
func NewBasicModel(cfgPath string) Model {
	log.Println("Initialising simple random model")

	// Format of yaml file
	type Config struct {
		SIZE struct {
			X int
			Y int
		}
		I     float64
		DLEN  float64
		MU    float64
		LAM   float64
		GAMMA float64
		B     []float64
		K     []float64
	}

	cfgArray, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("Can not ReadFile: %v", err))
	}
	log.Println("Reading config from file: ", cfgPath)

	cfg := Config{}
	err = yaml.Unmarshal(cfgArray, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Can not unmarshal yaml cfg file: %v", err))
	}

	model := Model{}
	model.mu = cfg.MU
	model.b = cfg.B
	model.k = cfg.K
	model.x = cfg.SIZE.X
	model.y = cfg.SIZE.Y
	model.Lam = cfg.LAM
	model.Gamma = cfg.GAMMA

	// Spins would be enumerated as
	//		id = pos_y * x + pos_x
	// where pos_x and pos_y is a atom position in grid
	// and x is width of grid
	SpinIdx := func(posX, posY int) int {
		return posY*cfg.SIZE.X + posX
	}

	// Initialize spins
	spins := make(la.Column, cfg.SIZE.X*cfg.SIZE.Y)
	for i := 0; i < len(spins); i++ {
		spins[i] = la.NewRandomVector(3).Unit()
	}
	model.Spins = spins

	// Initialize bounds bentween athoms
	bounds := make([][]Bound, len(spins))
	for posX := 0; posX < cfg.SIZE.X; posX++ {
		for posY := 0; posY < cfg.SIZE.Y; posY++ {
			selfIdx := SpinIdx(posX, posY)

			// Neighbours indexes
			topIdx := SpinIdx((posX+1)%cfg.SIZE.X, posY)
			botIdx := SpinIdx((cfg.SIZE.X+posX-1)%cfg.SIZE.X, posY)
			rghIdx := SpinIdx(posX, (posY+1)%cfg.SIZE.Y)
			lftIdx := SpinIdx(posX, (posY+cfg.SIZE.Y-1)%cfg.SIZE.Y)

			dVert := la.NewVector(0, cfg.DLEN, 0)
			dHor := la.NewVector(cfg.DLEN, 0, 0)

			selfBounds := make([]Bound, 4)
			selfBounds[0] = Bound{topIdx, dVert, cfg.I}
			selfBounds[1] = Bound{botIdx, dVert.Neg(), cfg.I}
			selfBounds[2] = Bound{rghIdx, dHor, cfg.I}
			selfBounds[3] = Bound{lftIdx, dHor.Neg(), cfg.I}

			bounds[selfIdx] = selfBounds
		}
	}
	model.bounds = bounds

	log.Printf("Generated: \n%v", model)
	return model
}
