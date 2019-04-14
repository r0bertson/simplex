package main

import (
	"fmt"

	s "github.com/r0bertson/simplex/simplex"
)

func main() {
	problem := [][]float32{
		{-3, -5, 0, 0, 0, 0},
		//{3, 5, 0, 0, 0, 0}, test isOptimal()
		{1, 0, 1, 0, 0, 4},
		{0, 2, 1, 0, 1, 12},
		{3, 2, 0, 0, 0, 18}}
	simplex := s.Simplex{}
	simplex.Build(problem, len(problem), len(problem[0]))
	simplex.PrintTableau()
	fmt.Println(simplex.IsOptimal())

}
