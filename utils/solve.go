package utils

import (
	"errors"
	"fmt"
	"strconv"
)

// SolveGaussian solves the system of linear equations via The Gaussian Elimination.
func SolveGaussian(eqM [][]Rational, printTriangularForm bool) (res [][]Rational, err error) {
	if len(eqM) > len(eqM[0])-1 {
		err = errors.New("the number of equations can not be greater than the number of variables")
		return
	}

	dl, i, j := containsDuplicatesLines(eqM)
	if dl {
		err = fmt.Errorf("provided matrix contains duplicate lines (%d and %d)", i+1, j+1)
		return
	}

	for i := 0; i < len(eqM)-1; i++ {
		eqM = sortMatrix(eqM, i)

		var varC Rational
		for k := i; k < len(eqM); k++ {
			if k == i {
				varC = eqM[k][i]
			} else {
				multipliedLine := make([]Rational, len(eqM[i]))
				for z, zv := range eqM[i] {
					multipliedLine[z] = zv.Multiply(eqM[k][i].Divide(varC)).MultiplyByNum(-1)
				}
				newLine := make([]Rational, len(eqM[k]))
				for z, zv := range eqM[k] {
					newLine[z] = zv.Add(multipliedLine[z])
				}
				eqM[k] = newLine
			}
		}
	}

	// Removing empty lines and inverting the matrix.
	var resultEqM [][]Rational
	for i := len(eqM) - 1; i >= 0; i-- {
		if !RationalsAreNull(eqM[i]) {
			resultEqM = append(resultEqM, eqM[i])
		}
	}

	getFirstNonZeroIndex := func(sl []Rational) (index int) {
		for i, v := range sl {
			if v.GetNumerator() != 0 {
				index = i
				return
			}
		}
		return
	}

	// Back substitution.
	for z := 0; z < len(resultEqM)-1; z++ {
		var processIndex int
		var firstLine []Rational
		for i := z; i < len(resultEqM); i++ {
			v := resultEqM[i]
			if i == z {
				processIndex = getFirstNonZeroIndex(v)
				firstLine = v
			} else {
				mult := v[processIndex].Divide(firstLine[processIndex]).MultiplyByNum(-1)
				for j, jv := range v {
					resultEqM[i][j] = firstLine[j].Multiply(mult).Add(jv)
				}
			}
		}
	}

	if printTriangularForm {
		for i := len(resultEqM) - 1; i >= 0; i-- {
			var str string
			for _, jv := range resultEqM[i] {
				str += strconv.FormatFloat(jv.Float64(), 'f', 2, 64) + ","
			}
			str = str[:len(str)-1]
			fmt.Println(str)
		}
	}

	// Calculating variables.
	res = make([][]Rational, len(eqM[0])-1)
	if getFirstNonZeroIndex(resultEqM[0]) == len(resultEqM[0])-2 {
		// All the variables have been found.
		for i, iv := range resultEqM {
			index := len(res) - 1 - i
			res[index] = append(res[index], iv[len(iv)-1].Divide(iv[len(resultEqM)-1-i]))
		}
	} else {
		// Some variables remained unknown.
		var unknownStart, unknownEnd int
		for i, iv := range resultEqM {
			fnz := getFirstNonZeroIndex(iv)
			var firstRes []Rational
			firstRes = append(firstRes, iv[len(iv)-1].Divide(iv[fnz]))
			if i == 0 {
				unknownStart = fnz + 1
				unknownEnd = len(iv) - 2
				for j := unknownEnd; j >= unknownStart; j-- {
					res[j] = []Rational{New(0, 0)}
					firstRes = append(firstRes, iv[j].Divide(iv[fnz]))
				}
			} else {
				for j := unknownEnd; j >= unknownStart; j-- {
					firstRes = append(firstRes, iv[j].Divide(iv[fnz]))
				}
			}
			res[fnz] = firstRes
		}
	}

	return
}

func sortMatrix(m [][]Rational, initRow int) (m2 [][]Rational) {
	indexed := make(map[int]bool)

	for i := 0; i < initRow; i++ {
		m2 = append(m2, m[i])
		indexed[i] = true
	}

	greaterThanMax := func(rr1, rr2 []Rational) (greater bool) {
		for i := 0; i < len(rr1); i++ {
			if rr1[i].GetModule().GreaterThan(rr2[i].GetModule()) {
				greater = true
				return
			} else if rr1[i].GetModule().LessThan(rr2[i].GetModule()) {
				return
			}
		}
		return
	}

	type maxStruct struct {
		index   int
		element []Rational
	}

	for i := initRow; i < len(m); i++ {
		max := maxStruct{-1, make([]Rational, len(m[i]))}
		var firstNotIndexed int
		for k, kv := range m {
			if !indexed[k] {
				firstNotIndexed = k
				if greaterThanMax(kv, max.element) {
					max.index = k
					max.element = kv
				}
			}
		}
		if max.index != -1 {
			m2 = append(m2, max.element)
			indexed[max.index] = true
		} else {
			m2 = append(m2, m[firstNotIndexed])
			indexed[firstNotIndexed] = true
		}
	}

	return
}

func containsDuplicatesLines(eqM [][]Rational) (contains bool, l1, l2 int) {
	for i := 0; i < len(eqM); i++ {
		for j := i + 1; j < len(eqM); j++ {
			var equalElements int
			for k := 0; k < len(eqM[i]); k++ {
				if eqM[i][k] == eqM[j][k] {
					equalElements++
				} else {
					break
				}
			}
			if equalElements == len(eqM[i]) {
				contains = true
				l1 = i
				l2 = j
				return
			}
		}
	}
	return
}
