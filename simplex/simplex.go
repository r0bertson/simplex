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
var maxRatio = 99999999.99

// Simplex holds the information of a Linear Problem an the state of the solution process
type Simplex struct {
	LP            *ltx.LinearProblem
	Status        string
	Variables     []string
	Base          []string
	Tableau       [][]float64
	RowsSize      int
	ColumnsSize   int
	NumIterations int
}

// BuildImportedProblem builds the Tableau based on a LinearProblem exported by the ltx-parser
func (s *Simplex) BuildImportedProblem(lp *ltx.LinearProblem) {
	s.LP = lp
	s.RowsSize = len(s.LP.Constraints) + 1
	s.ColumnsSize = s.getColumnsLength()

	s.Tableau = make([][]float64, s.RowsSize) //initializing Tableau's rows
	for i := 0; i < s.RowsSize; i++ {
		//initializing Tableau's columns
		s.Tableau[i] = make([]float64, s.ColumnsSize)
	}

	//Adding Z to base and applying a coefficient to it ( MIN Z is the same as MAX -Z)
	var coef float64
	if s.LP.ObjectiveFunction.OFType == "MAX" {
		s.Base = append(s.Base, "Z")
		coef = -1.0
	} else {
		s.Base = append(s.Base, "-Z")
		coef = 1.0
	}

	// Filling variables coefficients of the OF into Tableau and storing its name into variable array
	for i, v := range s.LP.ObjectiveFunction.Variables {
		s.Variables = append(s.Variables, v.Name)
		s.Tableau[0][i] = v.Coefficient * coef
	}

	//for naming purpouses
	slackCount := 0
	artificialCount := 0

	// Filling coefficients of constraints adding slack/artificial variables
	for i, constraint := range lp.Constraints {

		for j := 0; j < len(constraint.LH); j++ {
			col := s.findVariableColumn(constraint.LH[j].Name)
			s.Tableau[i+1][col] = constraint.LH[j].Coefficient

		}

		//place right hand side of constraint on Tableau
		s.Tableau[i+1][s.ColumnsSize-1] = constraint.RH

		switch constraint.Operator {
		case "<=":
			slackCount++
			slackName := "S" + strconv.Itoa(slackCount) //creating name
			s.Variables = append(s.Variables, slackName)
			s.Base = append(s.Base, slackName)
			s.Tableau[i+1][s.findVariableColumn(slackName)] = 1
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
			s.Variables = append(s.Variables, slackName)
			s.Variables = append(s.Variables, artificialName)
			s.Base = append(s.Base, artificialName)

			slackColumn := s.findVariableColumn(slackName)
			s.Tableau[i+1][slackColumn] = -1
			//s.Tableau[0][slackColumn] = bigM

			artificialColumn := s.findVariableColumn(artificialName)
			s.Tableau[i+1][artificialColumn] = 1
			s.Tableau[0][artificialColumn] = bigM

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
			s.Variables = append(s.Variables, artificialName)
			s.Base = append(s.Base, artificialName)

			artificialColumn := s.findVariableColumn(artificialName)
			s.Tableau[i+1][artificialColumn] = 1
			s.Tableau[0][artificialColumn] = bigM //STEP 2

			s.removeVBfromOF(i+1, artificialColumn) //STEP 3
		}

	}
}

/*BuildStandardizedProblem accepts a standardized Tableau without building it from the LP.
Useful when the Tableau is already given*/
func (s *Simplex) BuildStandardizedProblem(problem [][]float64, numberOfRows, numberOfColumns int) {
	s.RowsSize = numberOfRows
	s.ColumnsSize = numberOfColumns
	s.Tableau = problem
}

func (s *Simplex) PrintTableau() {
	//printing header
	fmt.Print(fmt.Sprintf("%s", "B.V.") + "\t")
	for i := 0; i < len(s.Variables); i++ {
		fmt.Print(fmt.Sprintf("%s        ", s.Variables[i]) + "\t")
	}
	fmt.Println()

	for i := 0; i < s.RowsSize; i++ {
		fmt.Print(fmt.Sprintf("%s", s.Base[i]) + "\t")
		for j := 0; j < s.ColumnsSize; j++ {
			element := fmt.Sprintf("%f", s.Tableau[i][j]) + "\t"
			fmt.Print(element)
		}
		fmt.Println()
	}
	fmt.Println()
}

// IsOptimal returns true if all values on last column are non-negative
func (s *Simplex) IsOptimal() bool {
	for i := 0; i < s.ColumnsSize-1; i++ {
		if s.Tableau[0][i] < 0 {
			return false
		}
	}
	s.Status = "OPTIMAL"
	return true
}

//FindPivotColumn obtain the index of the variable that will be placed at the base
func (s *Simplex) FindPivotColumn() int {
	chosenPosition := 0
	chosenAbsoluteValue := 0.0
	for i := 0; i < s.ColumnsSize-1; i++ {
		if s.Tableau[0][i] < 0 {
			abs := math.Abs(s.Tableau[0][i])
			if abs > chosenAbsoluteValue {
				chosenAbsoluteValue = abs
				chosenPosition = i
			}
		}
	}
	return chosenPosition
}

/*GetMinRatioRow obtain the index of the row which element most limits the increase of a variable
by applying the minimum ratio test */
func (s *Simplex) GetMinRatioRow(pivotColumn int) int {
	minRatio := maxRatio
	position := -1
	for i := 1; i < s.RowsSize; i++ {
		elementValue := s.Tableau[i][pivotColumn]
		if elementValue > 0 && s.Tableau[i][s.ColumnsSize-1] > 0 {
			ratio := s.Tableau[i][s.ColumnsSize-1] / elementValue
			if ratio < minRatio { //Bland's rule (always choose the lower index)
				minRatio = ratio
				position = i
			}
		}
	}
	return position
}

// removeVBfromOF removes a basic variable from the objective function
func (s *Simplex) removeVBfromOF(pivotRow, pivotColumn int) {

	pivotElement := s.Tableau[0][pivotColumn]
	for j := 0; j < s.ColumnsSize; j++ {
		s.Tableau[0][j] = s.Tableau[0][j] - (s.Tableau[pivotRow][j] * pivotElement)
	}
}

//UpdateTableau performs algebraic operations to swap variables at the base
func (s *Simplex) UpdateTableau(pivotRow, pivotColumn int) {
	pivotElement := s.Tableau[pivotRow][pivotColumn]
	pivotColumnValues := s.getColumn(pivotColumn)

	//divide all elements in pivotRow by pivotElement value
	for i := 0; i < s.ColumnsSize; i++ {
		s.Tableau[pivotRow][i] = s.Tableau[pivotRow][i] / pivotElement
	}

	for i := 0; i < s.RowsSize; i++ {
		//skip pivotRow and when element is already 0
		if i != pivotRow && !IsEqual(s.Tableau[i][pivotColumn], 0.0) {
			for j := 0; j < s.ColumnsSize; j++ {
				s.Tableau[i][j] = s.Tableau[i][j] - (s.Tableau[pivotRow][j] * pivotColumnValues[i])
			}
		}
	}
}

// getColumn stores all values of a specific Tableau's column on a slice and returns it
func (s *Simplex) getColumn(pivotColumn int) []float64 {
	var columnValues []float64
	for i := 0; i < s.RowsSize; i++ {
		columnValues = append(columnValues, s.Tableau[i][pivotColumn])
	}
	return columnValues
}

//Solve loop untils it find the optimal answer or fails, printing every interaction on console
func (s *Simplex) Solve() {
	fmt.Println("Initial Tableau: ")
	s.PrintTableau()
	iteration := 0
	for !s.IsOptimal() {
		iteration++
		pivotColumn := s.FindPivotColumn()
		fmt.Println("Selected pivot column index: " + strconv.Itoa(pivotColumn))
		pivotRow := s.GetMinRatioRow(pivotColumn)
		if pivotRow == -1 {
			s.Status = "UNBOUND"
			fmt.Println("Solution is unbound due to negative values on variable: " + s.Variables[pivotColumn])
			break
		}
		fmt.Println("Selected pivot row index: " + strconv.Itoa(pivotRow))

		//swap variables on the base and update Basic Variables column (pivotColumn enters, pivotRow leaves)
		s.UpdateTableau(pivotRow, pivotColumn)      //swap variables
		s.Base[pivotRow] = s.Variables[pivotColumn] //update base

		fmt.Println("Tableau after " + strconv.Itoa(iteration) + " iteration(s):")
		s.PrintTableau()
	}

	if !s.isFeasible() && s.Status != "UNBOUND" {
		s.Status = "INFEASIBLE"
	}

	fmt.Println(s.Status)
}

//SolveQuietly loop untils it find the optimal answer or fails without printing any result
func (s *Simplex) SolveQuietly() {
	for !s.IsOptimal() {
		s.NumIterations++
		pivotColumn := s.FindPivotColumn()
		pivotRow := s.GetMinRatioRow(pivotColumn)
		if pivotRow == -1 {
			s.Status = "UNBOUND"
			fmt.Println("Solution is unbound due to negative values on variable: " + s.Variables[pivotColumn])
			break
		}
		//swap variables on the base and update Basic Variables column (pivotColumn enters, pivotRow leaves)
		s.UpdateTableau(pivotRow, pivotColumn)      //swap variables
		s.Base[pivotRow] = s.Variables[pivotColumn] //update base

	}
	if !s.isFeasible() && s.Status != "UNBOUND" {
		s.Status = "INFEASIBLE"
	}

}

// isFeasible returns true if all coefficients on the last column of Tableau are non-negative
func (s *Simplex) isFeasible() bool {
	for i := 1; i < s.RowsSize; i++ {
		if s.Tableau[i][s.ColumnsSize-1] < 0 {
			return false
		}
	}
	return true
}

// IsEqual returns true if the difference between two floats is lower then EPSILON
func IsEqual(a, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

// NotEqual returns true if the difference between two floats is higher then EPSILON
func NotEqual(a, b float64) bool {
	return !IsEqual(a, b)
}

// getCoefficient returns the coefficient of a variable if it is present on a constraint. If not, returns 0.0
func getCoefficient(name string, c ltx.Constraint) float64 {
	for _, v := range c.LH {
		if v.Name == name {
			return v.Coefficient
		}
	}
	return 0.0
}

// findVariableColumn search for the index of a variable on variables slice, which is the same as Tableau's
func (s *Simplex) findVariableColumn(name string) int {
	for index, value := range s.Variables {
		if value == name {
			return index
		}
	}
	return -1
}

// getColumnsLength returns the expected number of columns on the Tableau
func (s *Simplex) getColumnsLength() int {
	count := 1                                     // one for RH
	count += len(s.LP.ObjectiveFunction.Variables) //one for each variable
	for _, c := range s.LP.Constraints {
		switch c.Operator {
		case "<=", "=":
			count++ //one artificial
		case ">=":
			count += 2 //one slack, one artificial
		}
	}
	return count
}

//getOFCoefficient returns the coefficient that will multiplies the Z row on Tableau
func (s *Simplex) getOFCoefficient() float64 {
	if s.LP.ObjectiveFunction.OFType == "MAX" {
		return -1.0
	}
	return 1.0
}

// GetRanges()
func (s *Simplex) GetRanges() {
	/*  b = S∗∆b + b∗ ≥ 0
	MEANING: S* : Non-basic variables of final tableau
			 b* : RH of final tableau
	*/
	if s.Status != "OPTIMAL" {
		return
	}
	nonBasicVariables := s.getSlackVariables()

	var S [][]float64
	for _, nbv := range nonBasicVariables {
		S = append(S, s.getColumn(nbv)[1:])
	}

	var bAsterisk [][]float64
	bAsterisk = append(bAsterisk, s.getColumn(s.ColumnsSize - 1)[1:]) //ignoring Z row
	bAsterisk = transpose(bAsterisk)                                  //get column returns an array 1xm, transpose will return a matrix mx1

	var deltas [][]float64
	var operators [][]string
	for j := 0; j < len(S); j++ {
		var deltaB []float64
		var deltaBOperator []string
		for i := 0; i < len(S[0]); i++ {

			value := -bAsterisk[i][0] * 1 / S[j][i]
			if S[j][i] < 0.0 {
				deltaBOperator = append(deltaBOperator, "<=")
			} else {
				deltaBOperator = append(deltaBOperator, ">=")
			}
			deltaB = append(deltaB, value)
		}
		deltas = append(deltas, deltaB)
		operators = append(operators, deltaBOperator)
	}

	header := "ROW" + "\t" + "A.INC" + "\t" + "A.DEC"
	fmt.Println(header)
	for i := 0; i < len(deltas); i++ {
		var allowableIncrease, allowableDecrease float64
		allowableIncrease = math.Inf(1)
		allowableDecrease = math.Inf(-1)
		for j := 0; j < len(deltas[0]); j++ {
			temp := deltas[i][j]
			if operators[i][j] == "<=" {
				if temp < allowableIncrease {
					allowableIncrease = temp
				}
			} else {
				if temp > allowableDecrease {
					allowableDecrease = temp
				}
			}

		}
		line := fmt.Sprintf("%d", i+1) + "\t" + fmt.Sprintf("%f", allowableIncrease) + "\t" + fmt.Sprintf("%f", -1*allowableDecrease)
		fmt.Println(line)
	}
}

func multiplySlices(A, B [][]float64) [][]float64 {
	//check if number of columns on A is the same as the number of rows on B
	if len(A[0]) != len(B) {
		return nil
	}

	out := make([][]float64, len(A))
	for i := 0; i < len(A); i++ {
		out[i] = make([]float64, len(B[0]))
		for j := 0; j < len(B[0]); j++ {
			for k := 0; k < len(B); k++ {
				out[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return out
}

func addSlices(A, B [][]float64) [][]float64 {
	//check if number of rows and columns on A is the same as on B
	if len(A[0]) != len(B[0]) || len(A) != len(B) {
		return nil
	}

	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0]); j++ {
			A[i][j] += B[i][j]
		}
	}
	return A

}

func (s *Simplex) getSlackVariables() []int {
	var slackVariables []int
	for i, v := range s.Variables {
		isDecisionVariable := false
		//checks if variables are present in OF
		for j := 0; j < len(s.LP.ObjectiveFunction.Variables) && !isDecisionVariable; j++ {
			if s.LP.ObjectiveFunction.Variables[j].Name == v {
				isDecisionVariable = true
			}
		}
		if !isDecisionVariable {
			slackVariables = append(slackVariables, i)
		}
	}
	return slackVariables
}

func transpose(slice [][]float64) [][]float64 {
	row := len(slice[0])
	col := len(slice)
	result := make([][]float64, row)
	for i := range result {
		result[i] = make([]float64, col)
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
