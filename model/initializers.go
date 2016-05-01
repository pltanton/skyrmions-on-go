package model

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/plotnikovanton/golinal"
	"gopkg.in/yaml.v2"
)

// Spins would be enumerated as
//		id = pos_y * x + pos_x
// where pos_x and pos_y is a atom position in grid
// and x is width of grid
func spinIdx(posX, posY, x int) int {
	return posY*x + posX
}

// Reads configuration file by given path into given structure
func readConfig(cfgPath string, out interface{}) {
	cfgArray, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("Can not ReadFile: %v", err))
	}
	log.Println("Reading configuration from file: ", cfgPath)

	err = yaml.Unmarshal(cfgArray, out)
	if err != nil {
		panic(fmt.Sprintf("Can not unmarshal yaml cfg file: %v", err))
	}
}

// NewBasicModel creates simple 2D grid model form yaml file
// All spins would be oriented randomly at start
// Crystal field would be formed as 4 neighbors system
func NewBasicModel(cfgPath string) Model {
	log.Println("Initializing simple random model")

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

	cfg := new(Config)
	readConfig(cfgPath, cfg)

	model := Model{}
	model.mu = cfg.MU
	model.b = cfg.B
	model.k = cfg.K
	model.x = cfg.SIZE.X
	model.y = cfg.SIZE.Y
	model.Lam = cfg.LAM
	model.Gamma = cfg.GAMMA

	// Initialize spins
	spins := make(la.Column, cfg.SIZE.X*cfg.SIZE.Y)
	for i := 0; i < len(spins); i++ {
		spins[i] = la.NewRandomVector(3).Unit()
	}
	model.Spins = spins

	model.bounds = fourBoundGrid(cfg.SIZE.X, cfg.SIZE.Y, cfg.DLEN, cfg.I)

	log.Printf("Generated: \n%v", model)
	return model
}

func NewManualModel(cfgPath string) Model {
	log.Println("Initializing manual configured spins model")

	// Format of yaml file
	type Config struct {
		SIZE struct {
			X int
			Y int
		}
		I       float64
		DLEN    float64
		MU      float64
		LAM     float64
		GAMMA   float64
		B       []float64
		K       []float64
		DEFAULT []float64
		MANUAL  []struct {
			X    int
			Y    int
			VECT []float64
		}
	}

	cfg := new(Config)
	readConfig(cfgPath, cfg)

	model := Model{}
	model.mu = cfg.MU
	model.b = cfg.B
	model.k = cfg.K
	model.x = cfg.SIZE.X
	model.y = cfg.SIZE.Y
	model.Lam = cfg.LAM
	model.Gamma = cfg.GAMMA
	model.Spins = la.Vector(cfg.DEFAULT).Unit().NewColumn(model.x * model.y)
	model.bounds = fourBoundGrid(cfg.SIZE.X, cfg.SIZE.Y, cfg.DLEN, cfg.I)

	for _, val := range cfg.MANUAL {
		id := spinIdx(val.X, val.Y, model.x)
		model.Spins[id] = la.Vector(val.VECT).Unit()
	}

	return model
}

// fourBoundGrid returns matrix of simple 4 neighbors bounds
func fourBoundGrid(x, y int, dLen, i float64) (bounds [][]Bound) {
	length := x * y

	SpinIdx := func(posX, posY int) int { return spinIdx(posX, posY, x) }

	bounds = make([][]Bound, length)
	for posX := 0; posX < x; posX++ {
		for posY := 0; posY < y; posY++ {
			selfIdx := SpinIdx(posX, posY)

			// Neighbours indexes
			topIdx := SpinIdx((posX+1)%x, posY)
			botIdx := SpinIdx((x+posX-1)%x, posY)
			rghIdx := SpinIdx(posX, (posY+1)%y)
			lftIdx := SpinIdx(posX, (posY+y-1)%y)

			dVert := la.NewVector(0, dLen, 0)
			dHor := la.NewVector(dLen, 0, 0)

			selfBounds := make([]Bound, 4)
			selfBounds[0] = Bound{topIdx, dVert, i}
			selfBounds[1] = Bound{botIdx, dVert.Neg(), i}
			selfBounds[2] = Bound{rghIdx, dHor, i}
			selfBounds[3] = Bound{lftIdx, dHor.Neg(), i}

			bounds[selfIdx] = selfBounds
		}
	}
	return
}
