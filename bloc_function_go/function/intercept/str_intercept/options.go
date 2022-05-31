package str_intercept

import (
	bloc_client "github.com/fBloc/bloc-client-go"
)

type strCompare int

const (
	strCompareEq strCompare = iota + 1
	strCompareNotEq
	strCompareContains
	strCompareNotContains
	strCompareStartsWith
	strCompareNotStartsWith
	strCompareEndsWith
	strCompareNotEndsWith
	strCompareEqIgnoreCase
	strCompareNotEqIgnoreCase
	strCompareContainsIgnoreCase
	strCompareNotContainsIgnoreCase
	strCompareStartsWithIgnoreCase
	strCompareNotStartsWithIgnoreCase
	strCompareEndsWithIgnoreCase
	strCompareNotEndsWithIgnoreCase
	max
)

func (sC strCompare) Value() int {
	return int(sC)
}

func (sC strCompare) IsValid() bool {
	return sC >= strCompareEq && sC < max
}

func (sC strCompare) String() string {
	switch sC {
	// eq
	case strCompareEq:
		return "equal"
	case strCompareEqIgnoreCase:
		return "equal(ignore case)"
	case strCompareNotEq:
		return "not equal"
	case strCompareNotEqIgnoreCase:
		return "not equal(ignore case)"
	// contains
	case strCompareContains:
		return "contains"
	case strCompareNotContains:
		return "not contains"
	case strCompareContainsIgnoreCase:
		return "contains(ignore case)"
	case strCompareNotContainsIgnoreCase:
		return "not contains(ignore case)"
	// start with
	case strCompareNotStartsWith:
		return "not start with"
	case strCompareStartsWith:
		return "start with"
	case strCompareStartsWithIgnoreCase:
		return "start with(ignore case)"
	case strCompareNotStartsWithIgnoreCase:
		return "not start with(ignore case)"
	// end with
	case strCompareEndsWith:
		return "end with"
	case strCompareNotEndsWith:
		return "not end with"
	case strCompareEndsWithIgnoreCase:
		return "end with(ignore case)"
	case strCompareNotEndsWithIgnoreCase:
		return "not end with(ignore case)"
	}
	return "unknown"
}

func strCompareSelections() []bloc_client.SelectOption {
	ret := make([]bloc_client.SelectOption, 0, max-1)
	for i := 1; i < int(max); i++ {
		ret = append(
			ret,
			bloc_client.SelectOption{
				Label: strCompare(i).String(),
				Value: i,
			},
		)
	}
	return ret
}

type interceptOrPass int

const (
	intercept interceptOrPass = iota + 1
	pass
	maxI
)

func (iOP interceptOrPass) Value() int {
	return int(iOP)
}

func (iOP interceptOrPass) IsValid() bool {
	return iOP >= intercept && iOP < maxI
}

func (iOP interceptOrPass) String() string {
	switch iOP {
	case intercept:
		return "intercept"
	case pass:
		return "pass"
	}
	return "unknown"
}

func interceptSelections() []bloc_client.SelectOption {
	ret := make([]bloc_client.SelectOption, 0, maxI-1)
	for i := 1; i < int(maxI); i++ {
		ret = append(
			ret,
			bloc_client.SelectOption{
				Label: interceptOrPass(i).String(),
				Value: i,
			},
		)
	}
	return ret
}

type blank int

const (
	isBlank blank = iota + 1
	isNotBlank
	maxK
)

func (b blank) Value() int {
	return int(b)
}

func (b blank) IsValid() bool {
	return b >= isBlank && b < maxK
}

func (b blank) String() string {
	switch b {
	case isBlank:
		return "is blank"
	case isNotBlank:
		return "is not blank"
	}
	return "unknown"
}

func blankSelections() []bloc_client.SelectOption {
	ret := make([]bloc_client.SelectOption, 0, maxK-1)
	for i := 1; i < int(maxK); i++ {
		ret = append(
			ret,
			bloc_client.SelectOption{
				Label: blank(i).String(),
				Value: i,
			},
		)
	}
	return ret
}
