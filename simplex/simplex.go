package simplex

import (
	"fmt"
)

/*STEPS
1. BUILD TABLEAU
---- BEGIN LOOP
2. CHECK IF SOLUTION IS OPTIMAL
3. SELECT PIVOT COLUMN (HIGHEST ABSOLUTE VALUE)
4. SELECT PIVOT LINE (MINIMAL VALUE ON THE ~~~~~PUT TEST NAME HERE~~~~)
5. DO THE MATH TO REDUCE THE PIVOT ELEMENT TO 1 AND THE REST OF THE PIVOT COLUMN ELEMENTS TO ZERO)
---- END LOOP
*/

type Simplex struct {
	tableau     [][]float32
	rowsSize    int
	columnsSize int
}

func (s *Simplex) Build(problem [][]float32, numberOfRows, numberOfColumns int) {
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
