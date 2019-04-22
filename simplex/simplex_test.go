package simplex_test

import (
	"testing"

	ltx "github.com/r0bertson/ltx-parser"
	s "github.com/r0bertson/simplex/simplex"
)

func TestSimplexSolve(t *testing.T) {
	var tests = []struct {
		lp      *ltx.LinearProblem
		OFValue float64
	}{
		// Original
		{
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MAX", Variables: []ltx.Variable{{"X1", 3}, {"X2", 5}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X1", 1}}, Operator: "<=", RH: 4.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X2", 2}}, Operator: "<=", RH: 12.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X1", 3}, {"X2", 2}}, Operator: "<=", RH: 18.0}},
			},
			OFValue: 36.0,
		},
		{
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MIN", Variables: []ltx.Variable{{"X", 0.4}, {"Y", .5}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X", 0.3}, {"Y", 0.1}}, Operator: "<=", RH: 2.7},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X", 0.5}, {"Y", 0.5}}, Operator: "=", RH: 6.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X", 0.6}, {"Y", 0.4}}, Operator: ">=", RH: 6.0}},
			},
			OFValue: -5.25,
		},
		{
			lp: &ltx.LinearProblem{
				ObjectiveFunction: ltx.OF{OFType: "MAX", Variables: []ltx.Variable{{"X", 3}, {"Y", 5}}},
				Constraints: []ltx.Constraint{{Name: "", LH: []ltx.Variable{{"X", 1}}, Operator: "<=", RH: 4.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"Y", 2}}, Operator: "<=", RH: 12.0},
					ltx.Constraint{Name: "", LH: []ltx.Variable{{"X", 3}, {"Y", +2}}, Operator: "=", RH: 18.0}},
			},
			OFValue: 36.0,
		},
	}

	for i, tt := range tests {
		lp := s.Simplex{}
		lp.BuildImportedProblem(tt.lp)
		lp.SolveQuietly()
		if result := lp.Tableau[0][lp.ColumnsSize-1]; s.NotEqual(result, tt.OFValue) {
			t.Errorf("%d: error mismatch:\n  exp=%f\n  got=%f\n\n", i, tt.OFValue, result)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
