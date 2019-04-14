package simplex

import (
	"fmt"
	"math"
	"strconv"
)

type Simplex struct {
	variables   []string
	tableau     [][]float64
	rowsSize    int
	columnsSize int
}

func (s *Simplex) Build(problem [][]float64, numberOfRows, numberOfColumns int) {
	s.rowsSize = numberOfRows
	s.columnsSize = numberOfColumns
	s.tableau = problem
}

func (s *Simplex) PrintTableau() {
	//TODO: ADD HEADER
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
	for i := 0; i < s.rowsSize; i++ {
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
	minRatio := 9999.99
	position := 1
	for i := 1; i < s.rowsSize; i++ {
		elementValue := s.tableau[i][pivotColumn]
		if elementValue != 0 {
			ratio := s.tableau[i][s.columnsSize-1] / elementValue
			if ratio < minRatio { //Bland's rule (always choose the lower index)
				minRatio = ratio
				position = i
			}
		}
	}
	return position
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
		if i != pivotRow && s.tableau[i][pivotColumn] != 0.0 {
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

}
