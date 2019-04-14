package main

import (
	s "github.com/r0bertson/simplex/simplex"
)

func main() {
	problem := [][]float64{
		{-3, -5, 0, 0, 0, 0},
		{1, 0, 1, 0, 0, 4},
		{0, 2, 0, 1, 0, 12},
		{3, 2, 0, 0, 1, 18}}
	simplex := s.Simplex{}
	simplex.Build(problem, len(problem), len(problem[0]))
	simplex.Solve()
}
