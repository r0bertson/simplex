package simplex

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/*
--------------------------------------- TODO -------------------------------------------
THIS FILE WILL BE A HELPER THAT WILL CONVERT A LP ON LINDO'S MODEL SINTAX INTO A TABLEAU
----------------------------------------------------------------------------------------
LINDO'S MODEL SINTAX:

MAX 10 X1 + 15 X2
SUBJECT TO
 X1 < 10
 X2 < 12
 X1 + 2 X1 < 16
END
----------------------------------------------------------------------------------------
WITH THE FOLLOWING RULES:
1. NOT CASE SENSITIVE
2. MUST START WITH MIN OR MAX
3. "SUBJECT TO" CAN BE REPLACED WITH ANY OF THE FOLLOWING OPTIONS (SUBJECT TO / SUCH
	THAT /	S.T. / ST )
4. VARIABLE NAMES MUST BE AT MOST 8 CHARACTERS LONG, BEGINNING WITH A LETTER. THE OTHER
	CHARACTERS CAN'T BE ANY OF THE FOLLOWING 		 : ! ) + - = < >
5. ONLY CONSTANT VALUES ON THE RIGHT HAND SIDE OF A CONSTRAINT EQUATIONS
6. ONLY VARIABLES ARE ALLOWED ON THE LEFT HAND SIDE OF A CONSTRAINT EQUATIONS
7. NAME OF CONSTRAINTS IS OPTIONAL AND FOLLOW THE SAME CONVENTIONS AS VARIABLE NAMES AND
	END WITH A CLOSING PARENTHESIS ')'

/*

type LinearProblem {

}
func ImportLP() {


}
*/
func loadFile(filename string) {
	file, err := os.Open(filename + ".ltx")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	fmt.Print(b)

}
