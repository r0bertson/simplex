package main

import (
	"fmt"
	"strings"

	ltx "github.com/r0bertson/ltx-parser"
	s "github.com/r0bertson/simplex/simplex"
)

func main() {
	// problem := [][]float64{}
	// /*
	// 	problem := [][]float64{
	// 		{-3, -5, 0, 0, 0, 0},
	// 		{1, 0, 1, 0, 0, 4},
	// 		{0, 2, 0, 1, 0, 12},
	// 		{3, 2, 0, 0, 1, 18}}
	// */
	// //WITH EQUALITY (USING BIG M)
	// problem = [][]float64{
	// 	{-3, -5, 0, 0, s.bigM, 0},
	// 	{1, 0, 1, 0, 0, 4},
	// 	{0, 2, 0, 1, 0, 12},
	// 	{3, 2, 0, 0, 1, 18}}

	// //WITH EQUALITY  AND LESS THEN CONSTRAINT
	// problem = [][]float64{
	// 	{-1.1*bigM + 0.4, -0.9*s.bigM + 0.5, 0, 0, s.bigM, 0, -12 * s.bigM},
	// 	{0.3, 0.1, 1, 0, 0, 0, 2.7},
	// 	{0.5, 0.5, 0, 1, 0, 0, 6},
	// 	{0.6, 0.4, 0, 0, -1, 1, 6}}

	// //EQUALITY, LESS THEN AND MINIMIZATION (MAX : -Z)
	// problem2 := [][]float64{
	// 	{-7*s.bigM + 4, -4*s.bigM + 1, 0, s.bigM, 0, 0, -9 * s.bigM},
	// 	{3, 1, 1, 0, 0, 0, 3},
	// 	{4, 3, 0, -1, 1, 0, 6},
	// 	{1, 2, 0, 0, 0, 1, 3}}

	// /*problem2 = [][]float64{
	// {-1, -1, 0, 0, 0, 0},
	// {1, 0, 1, 0, 0, 6},
	// {0, 1, 0, 1, 0, 6},
	// {1, 1, 0, 0, -1, 11}}*/

	// simplex := s.Simplex{}
	// simplex.Build(problem, len(problem), len(problem[0]))
	// simplex.Solve()

	// simplex2 := s.Simplex{}
	// simplex2.Build(problem2, len(problem2), len(problem2[0]))
	// simplex2.Solve()
	problem := `MIN 0.4 X1 + 0.5 X2
				SUBJECT TO
					0.3 X1 + 0.1 X2 <= 2.7
					0.5 X1 + 0.5 X2 = 6
					0.6 X1 + 0.4 X2 >= 6
				END`
	/*
		problem2:= `MAX 3 X1 + 5 X2
		SUBJECT TO
			X1 <= 4
			2 X2 <= 12
			3 X1 + 2 X2 = 18
		END`
	*/
	parser := ltx.NewParser(strings.NewReader(problem))
	lp, _ := parser.Parse()

	fmt.Println(lp.Constraints)
	simplex := s.Simplex{}
	simplex.BuildImportedProblem(lp)
	fmt.Println(simplex)
	simplex.Solve()
	//simplex.Build(problem, len(problem), len(problem[0]))
	//simplex.Solve()

}
