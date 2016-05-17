package main

import (
	"github.com/plotnikovanton/skyrmions_on_go/iterator"
	"github.com/plotnikovanton/skyrmions_on_go/model"
)

func main() {
	//m := model.NewBasicModel("/home/anton/workdir/gocode/src/github.com/plotnikovanton/skyrmions_on_go/simple_model.yml")
	m := model.NewManualModel("/home/anton/workdir/gocode/src/github.com/plotnikovanton/skyrmions_on_go/manual_model.yml")
	//runner := iterator.NewSimpleIterator(&m)
	runner := iterator.NewSimpleSplittedIterator(&m)

	//runner.Gp.Pipe("set term gif animate size 600,720")
	//runner.Gp.Pipe("set output 'animate.gif'")

	runner.Delta = 0.1
	runner.Energy = true
	//runner.Times = 1
	runner.Run()
}
