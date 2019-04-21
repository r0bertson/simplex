package simplex

import (
	"fmt"
	"math"
	"strconv"

	ltx "github.com/r0bertson/ltx-parser"
)

// EPSILON value of the error while comparing two floats
var EPSILON = 0.000000001
var bigM = 1000000.00

type Simplex struct {
	lp          *ltx.LinearProblem
	status      string
	variables   []string
	base        []int
	tableau     [][]float64
	rowsSize    int
	columnsSize int
}

func (s *Simplex) BuildImportedProblem(lp *ltx.LinearProblem) {
	s.lp = lp
	s.rowsSize = len(lp.Constraints) + 1
	s.columnsSize = s.getColumnsLength()

	s.tableau = make([][]float64, s.rowsSize) //initializing tableau's rows
	for i := 0; i < s.rowsSize; i++ {
		//initializing tableau's columns
		s.tableau[i] = make([]float64, s.columnsSize)
	}

	s.PrintTableau()
	// Filling variables coefficients of the OF into tableau and storing its name into variable array

	for i, v := range lp.ObjectiveFunction.Variables {
		s.variables = append(s.variables, v.Name)
		s.tableau[0][i] = v.Coefficient * -1.0
	}

	//for naming purpouses
	slackCount := 0
	artificialCount := 0

	// Filling coefficients of constraints adding slack/artificial variables
	for i, constraint := range lp.Constraints {

		for j := 0; j < len(constraint.LH); j++ {
			col := s.findVariableColumn(constraint.LH[j].Name)
			s.tableau[i+1][col] = constraint.LH[j].Coefficient

		}

		//place right hand side of constraint on tableau
		s.tableau[i+1][s.columnsSize-1] = constraint.RH

		switch constraint.Operator {
		case "<=":
			slackCount++
			slackName := "S" + strconv.Itoa(slackCount) //creating name
			s.variables = append(s.variables, slackName)
			s.base = append(s.base, len(s.variables)-1)
			s.tableau[i+1][s.findVariableColumn(slackName)] = 1
		case ">=":
			/* STEPS:
			1. Adds a surplus variable (negative) and artificial variable (A.V.)
			2. Place A.V. with a penalty on the OF
			3. Remove A.V. from the OF
			*/
			slackCount++
			artificialCount++
			slackName := "S" + strconv.Itoa(slackCount) //creating name
			artificialName := "A" + strconv.Itoa(artificialCount)
			// STEP 1
			s.variables = append(s.variables, slackName)
			s.variables = append(s.variables, artificialName)
			s.base = append(s.base, len(s.variables)-1)

			slackColumn := s.findVariableColumn(slackName)
			s.tableau[i+1][slackColumn] = -1
			s.tableau[0][slackColumn] = bigM // STEP 2

			artificialColumn := s.findVariableColumn(artificialName)
			s.tableau[i+1][artificialColumn] = 1
			s.tableau[0][artificialColumn] = bigM

			s.removeVBfromOF(i+1, artificialColumn) //STEP 3

		case "=":
			/* STEPS:
			1. Adds artificial variable (A.V.)
			2. Place A.V. with a penalty on the OF
			3. Remove A.V. from the OF
			*/
			artificialCount++
			artificialName := "A" + strconv.Itoa(artificialCount)
			//STEP 1
			s.variables = append(s.variables, artificialName)
			s.base = append(s.base, len(s.variables)-1)

			artificialColumn := s.findVariableColumn(artificialName)
			s.tableau[i+1][artificialColumn] = 1
			s.tableau[0][artificialColumn] = bigM //STEP 2

			s.removeVBfromOF(i+1, artificialColumn) //STEP 3
		}
		//place right hand side of constraint on tableau
		//s.tableau[i+1][s.columnsSize-1] = constraint.RH

	}
	fmt.Println(s.variables)
	s.PrintTableau()
}

func (s *Simplex) BuildStandardizedProblem(problem [][]float64, numberOfRows, numberOfColumns int) {
	s.rowsSize = numberOfRows
	s.columnsSize = numberOfColumns
	s.tableau = problem
}

func (s *Simplex) PrintTableau() {
	//PrintingHeader
	for i := 0; i < len(s.variables); i++ {
		fmt.Print(fmt.Sprintf("%s        ", s.variables[i]) + "\t")
	}

	fmt.Println()
	for i := 0; i < s.rowsSize; i++ {
		for j := 0; j < s.columnsSize; j++ {
			element := fmt.Sprintf("%f", s.tableau[i][j]) + "\t"
			fmt.Print(element)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (s *Simplex) IsOptimal() bool {
	for i := 0; i < s.columnsSize-1; i++ {
		if s.tableau[0][i] < 0 {
			return false
		}
	}
	return true
}

//FindPivotColumn obtain the index of the chosen pivotColumn
func (s *Simplex) FindPivotColumn() int {
	chosenPosition := 0
	chosenAbsoluteValue := 0.0
	for i := 0; i < s.columnsSize-1; i++ {
		if s.tableau[0][i] < 0 {
			abs := math.Abs(s.tableau[0][i])
			if abs > chosenAbsoluteValue {
				chosenAbsoluteValue = abs
				chosenPosition = i
			}
		}
	}
	return chosenPosition
}

//FindPivotRow obtain the index of the chosen pivotRow
func (s *Simplex) FindPivotRow(pivotColumn int) int {
	minRatio := 99999999.99
	position := 1
	for i := 1; i < s.rowsSize; i++ {
		elementValue := s.tableau[i][pivotColumn]
		if elementValue > 0 || s.tableau[i][s.columnsSize-1] > 0 {
			ratio := s.tableau[i][s.columnsSize-1] / elementValue
			if ratio < minRatio { //Bland's rule (always choose the lower index)
				minRatio = ratio
				position = i
			}
		}
	}
	return position
}

func (s *Simplex) removeVBfromOF(pivotRow, pivotColumn int) {

	pivotElement := s.tableau[0][pivotColumn]
	fmt.Println(fmt.Sprintf("%f", pivotElement))

	for j := 0; j < s.columnsSize; j++ {
		s.tableau[0][j] = s.tableau[0][j] - (s.tableau[pivotRow][j] * pivotElement)
	}
}

//UpdateTableau performs algebraic operations to swap variables at the base
func (s *Simplex) UpdateTableau(pivotRow, pivotColumn int) {
	pivotElement := s.tableau[pivotRow][pivotColumn]
	pivotColumnValues := s.getPivotColumn(pivotColumn)

	//divide all elements in pivotRow by pivotElement value
	for i := 0; i < s.columnsSize; i++ {
		s.tableau[pivotRow][i] = s.tableau[pivotRow][i] / pivotElement
	}

	for i := 0; i < s.rowsSize; i++ {
		//skip pivotRow and when element is already 0
		if i != pivotRow && !equals(s.tableau[i][pivotColumn], 0.0) {
			for j := 0; j < s.columnsSize; j++ {
				s.tableau[i][j] = s.tableau[i][j] - (s.tableau[pivotRow][j] * pivotColumnValues[i])
			}
		}
	}
}

func (s *Simplex) getPivotColumn(pivotColumn int) []float64 {
	var columnValues []float64
	for i := 0; i < s.rowsSize; i++ {
		columnValues = append(columnValues, s.tableau[i][pivotColumn])
	}
	return columnValues
}

//Solve loop until it find the optimal answer or fails (WHEN IT FAILS?)
func (s *Simplex) Solve() {
	fmt.Println("Initial Tableau: ")
	s.PrintTableau()
	iteration := 0
	for !s.IsOptimal() {
		iteration++
		pivotColumn := s.FindPivotColumn()
		fmt.Println("Selected pivot column index: " + strconv.Itoa(pivotColumn))
		pivotRow := s.FindPivotRow(pivotColumn)
		fmt.Println("Selected pivot row index: " + strconv.Itoa(pivotRow))
		s.UpdateTableau(pivotRow, pivotColumn)
		fmt.Println("Tableau after " + strconv.Itoa(iteration) + " iteration(s):")
		s.PrintTableau()
	}
	s.status = "OPTIMAL"
	if !s.isFeasible() {
		s.status = "INFEASIBLE"
	}
	fmt.Println(s.status)
}

func (s *Simplex) getRanges() {
	//TODO:  b = S∗∆b + b ∗ ≥ 0

}

//

func (s *Simplex) isFeasible() bool {
	for i := 1; i < s.rowsSize; i++ {
		if s.tableau[i][s.columnsSize-1] < 0 {
			return false
		}
	}
	return true
}

func equals(a, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

func getCoefficient(name string, c ltx.Constraint) float64 {
	for _, v := range c.LH {
		if v.Name == name {
			return v.Coefficient
		}
	}
	return 0.0
}

func (s *Simplex) findVariableColumn(name string) int {
	for index, value := range s.variables {
		if value == name {
			return index
		}
	}
	return -1
}

func (s *Simplex) getColumnsLength() int {
	count := 1
	count += len(s.lp.ObjectiveFunction.Variables) //one for each variable
	for _, c := range s.lp.Constraints {
		switch c.Operator {
		case "<=", "=":
			count++ //one artificial
		case ">=":
			count += 2 //one slack, one artificial
		}
	}
	return count
}
