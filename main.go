package main

import (
	s "github.com/r0bertson/simplex/simplex"
)

const bigM = 1000000.00

func main() {
	problem := [][]float64{}
	/*
		problem := [][]float64{
			{-3, -5, 0, 0, 0, 0},
			{1, 0, 1, 0, 0, 4},
			{0, 2, 0, 1, 0, 12},
			{3, 2, 0, 0, 1, 18}}
	*/
	problem = [][]float64{
		{-3, -5, 0, 0, bigM, 0},
		{1, 0, 1, 0, 0, 4},
		{0, 2, 0, 1, 0, 12},
		{3, 2, 0, 0, 1, 18}}

	problem = [][]float64{
		{-1.1*bigM + 0.4, -0.9*bigM + 0.5, 0, 0, bigM, 0, -12 * bigM},
		{0.3, 0.1, 1, 0, 0, 0, 2.7},
		{0.5, 0.5, 0, 1, 0, 0, 6},
		{0.6, 0.4, 0, 0, -1, 1, 6}}

	problem2 := [][]float64{
		{-7*bigM + 4, -4*bigM + 1, 0, bigM, 0, 0, -9 * bigM},
		{3, 1, 1, 0, 0, 0, 3},
		{4, 3, 0, -1, 1, 0, 6},
		{1, 2, 0, 0, 0, 1, 3}}

	simplex := s.Simplex{}
	simplex.Build(problem, len(problem), len(problem[0]))
	simplex.Solve()

	simplex2 := s.Simplex{}
	simplex2.Build(problem2, len(problem2), len(problem2[0]))
	simplex2.Solve()
}
