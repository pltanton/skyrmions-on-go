package main

import (
	"github.com/plotnikovanton/skyrmions_on_go/model"
)

func main() {
	m := model.NewBasicModel("/home/anton/workdir/gocode/src/github.com/plotnikovanton/skyrmions_on_go/simple_model.yml")
	runner := model.NewSimpleIterator(&m)
	runner.Energy = true
	//runner.Times = 1000
	runner.Run()
}
