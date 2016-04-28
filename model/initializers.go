package model

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/plotnikovanton/golinal"
	"gopkg.in/yaml.v2"
)

// Creates simple 2D grid model form yaml file
// All spins would be oriented randomly at start
// Crystalic field would be formed as 4 neigbours system
func NewBasicModel(cfg_path string) Model {
	log.Println("Initialising simple random model")

	// Format of yaml file
	type Config struct {
		SIZE struct {
			X int
			Y int
		}
		I     float64
		D_LEN float64
		MU    float64
		B     []float64
		K     []float64
	}

	cfg_array, err := ioutil.ReadFile(cfg_path)
	if err != nil {
		panic(fmt.Sprintf("Can not ReadFile: %v", err))
	}
	log.Println("Reading config from file: ", cfg_path)

	cfg := Config{}
	err = yaml.Unmarshal(cfg_array, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Can not unmarshal yaml cfg file: %v", err))
	}

	model := Model{}
	model.mu = cfg.MU
	model.b = cfg.B
	model.k = cfg.K
	model.x = cfg.SIZE.X
	model.y = cfg.SIZE.Y

	// Spins would be enumerated as
	//		id = pos_y * x + pos_x
	// where pos_x and pos_y is a atom position in grid
	// and x is width of grid
	SpinIdx := func(pos_x, pos_y int) int {
		return pos_y*cfg.SIZE.X + pos_x
	}

	// Initialize spins
	spins := make(la.Column, cfg.SIZE.X*cfg.SIZE.Y)
	for i := 0; i < len(spins); i++ {
		spins[i] = la.NewRandomVector(3)
	}
	model.spins = spins

	// Initialize bounds bentween athoms
	bounds := make([][]Bound, len(spins))
	for pos_x := 0; pos_x < cfg.SIZE.X; pos_x++ {
		for pos_y := 0; pos_y < cfg.SIZE.Y; pos_y++ {
			self_idx := SpinIdx(pos_x, pos_y)

			// Neighbours indexes
			top_idx := SpinIdx((pos_x+1)%cfg.SIZE.X, pos_y)
			bot_idx := SpinIdx((cfg.SIZE.X+pos_x-1)%cfg.SIZE.X, pos_y)
			rgh_idx := SpinIdx(pos_x, (pos_y+1)%cfg.SIZE.Y)
			lft_idx := SpinIdx(pos_x, (pos_y+cfg.SIZE.Y-1)%cfg.SIZE.Y)

			d_vert := la.NewVector(0, cfg.D_LEN, 0)
			d_hor := la.NewVector(cfg.D_LEN, 0, 0)

			self_bounds := make([]Bound, 4)
			self_bounds[0] = Bound{top_idx, d_vert, cfg.I}
			self_bounds[1] = Bound{bot_idx, d_vert.Neg(), cfg.I}
			self_bounds[2] = Bound{rgh_idx, d_hor, cfg.I}
			self_bounds[3] = Bound{lft_idx, d_hor.Neg(), cfg.I}

			bounds[self_idx] = self_bounds
		}
	}
	model.bounds = bounds

	log.Printf("Generated: \n%v", model)
	return model
}
