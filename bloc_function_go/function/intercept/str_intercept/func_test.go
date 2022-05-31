package str_intercept

import (
	"strings"
	"testing"

	bloc_client "github.com/fBloc/bloc-client-go"
)

var client = bloc_client.NewTestClient()

func TestNotAllowBlank(t *testing.T) {
	executeOpt := client.TestRunFunction(
		&StrIntercept{},
		[][]interface{}{
			{
				"",
			},
			{
				isBlank.Value(),   // is blank
				intercept.Value(), // intercept
			},
			{
				0,
				"",
				0,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if !executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should intercept")
	}

	executeOpt = client.TestRunFunction(
		&StrIntercept{},
		[][]interface{}{
			{
				"",
			},
			{
				isNotBlank.Value(), // is not blank
				pass.Value(),       // pass
			},
			{
				0,
				"",
				0,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if !executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should intercept")
	}
}

func TestAllowBlank(t *testing.T) {
	executeOpt := client.TestRunFunction(
		&StrIntercept{},
		[][]interface{}{
			{
				"",
			},
			{
				isBlank.Value(), // is blank
				pass.Value(),    // pass
			},
			{
				0,
				"",
				0,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should pass")
	}

	executeOpt = client.TestRunFunction(
		&StrIntercept{},
		[][]interface{}{
			{
				"",
			},
			{
				isNotBlank.Value(), // is not blank
				intercept.Value(),  // intercept
			},
			{
				0,
				"",
				0,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should pass")
	}
}

type compareData struct {
	compare         strCompare
	targetData      string
	interceptOrPass interceptOrPass
}

func commandRunCompare(t *testing.T, tag, data string, dataMapIntercept map[compareData]bool) {
	for k, v := range dataMapIntercept {
		executeOpt := client.TestRunFunction(
			&StrIntercept{},
			[][]interface{}{
				{
					data,
				},
				{
					0,
					0,
				},
				{
					k.compare.Value(),
					k.targetData,
					k.interceptOrPass.Value(),
				},
				{
					k.compare.Value(),
					k.targetData,
					k.interceptOrPass.Value(),
				},
			},
		)
		if !executeOpt.Suc {
			t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
		}
		if executeOpt.InterceptBelowFunctionRun != v {
			t.Fatalf(
				"%s - %s %s %s %s. expect intercept == %t",
				tag, k.interceptOrPass.String(), data, k.compare.String(), k.targetData, v)
		}
	}
}

func TestCompareEq(t *testing.T) {
	data := "xxXx"
	dataLower := strings.ToLower(data)

	dataMapIntercept := map[compareData]bool{
		{
			compare:         strCompareEq,
			targetData:      data,
			interceptOrPass: pass,
		}: false,
		{
			compare:         strCompareEq,
			targetData:      dataLower,
			interceptOrPass: pass,
		}: true,
		{
			compare:         strCompareEq,
			targetData:      data + "miss",
			interceptOrPass: intercept,
		}: false,
		{
			compare:         strCompareEq,
			targetData:      data + "miss",
			interceptOrPass: pass,
		}: true,
		{
			compare:         strCompareEq,
			targetData:      data,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestCompareEq", data, dataMapIntercept)
}

func TestCompareNotEq(t *testing.T) {
	data := "xxXx"
	dataLower := strings.ToLower(data)

	compareOp := strCompareNotEq
	dataMapIntercept := map[compareData]bool{
		{ // pass not eq
			compare:         compareOp,
			targetData:      data,
			interceptOrPass: pass,
		}: true,
		{ // pass not eq
			compare:         compareOp,
			targetData:      dataLower,
			interceptOrPass: pass,
		}: false,
		{ // intercept not eq
			compare:         compareOp,
			targetData:      data + "miss",
			interceptOrPass: intercept,
		}: true,
		{ // pass not eq
			compare:         compareOp,
			targetData:      data + "miss",
			interceptOrPass: pass,
		}: false,
		{ // intercept not eq
			compare:         compareOp,
			targetData:      data,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestCompareNotEq", data, dataMapIntercept)
}

func TestEqIgnoreCase(t *testing.T) {
	data := "xxXx"
	dataLower := strings.ToLower(data)

	compareOp := strCompareEqIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      dataLower,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      dataLower + "miss",
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      dataLower + "miss",
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      dataLower,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestEqIgnoreCase", data, dataMapIntercept)
}

func TestNotEqIgnoreCase(t *testing.T) {
	data := "xxXx"
	dataLower := strings.ToLower(data)

	compareOp := strCompareNotEqIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      dataLower,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      dataLower + "miss",
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      dataLower + "miss",
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      dataLower,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestNotEqIgnoreCase", data, dataMapIntercept)
}

func TestContains(t *testing.T) {
	data := "xxXx"
	subData := "xX"

	compareOp := strCompareContains
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestContains", data, dataMapIntercept)
}

func TestNotContains(t *testing.T) {
	data := "xxXx"
	subData := "xX"

	compareOp := strCompareNotContains
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestNotContains", data, dataMapIntercept)
}

func TestStartsWith(t *testing.T) {
	data := "xxXx"
	startsWithHit := "x"
	startsWithMiss := "X"

	compareOp := strCompareStartsWith
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestStartsWith", data, dataMapIntercept)
}

func TestNotStartsWith(t *testing.T) {
	data := "xxXx"
	startsWithHit := "x"
	startsWithMiss := "X"

	compareOp := strCompareNotStartsWith
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestNotStartsWith", data, dataMapIntercept)
}

func TestEndsWith(t *testing.T) {
	data := "xxXx"
	endsWithHit := "x"
	endsWithMiss := "X"

	compareOp := strCompareEndsWith
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestEndsWith", data, dataMapIntercept)
}

func TestNotEndsWith(t *testing.T) {
	data := "xxXx"
	endsWithHit := "x"
	endsWithMiss := "X"

	compareOp := strCompareNotEndsWith
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestNotEndsWith", data, dataMapIntercept)
}

func TestContainsIgnoreCase(t *testing.T) {
	data := "xxXx"
	subData := "xxxx"

	compareOp := strCompareContainsIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestContainsIgnoreCase", data, dataMapIntercept)
}

func TestNotContainsIgnoreCase(t *testing.T) {
	data := "xxXx"
	subData := "xxxx"

	compareOp := strCompareNotContainsIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      subData + "miss",
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      subData,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestNotContainsIgnoreCase", data, dataMapIntercept)
}

func TestStartsWithIgnoreCase(t *testing.T) {
	data := "xxXx"
	startsWithHit := "X"
	startsWithMiss := "o"

	compareOp := strCompareStartsWithIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestStartsWithIgnoreCase", data, dataMapIntercept)
}

func TestNotStartsWithIgnoreCase(t *testing.T) {
	data := "xxXx"
	startsWithHit := "X"
	startsWithMiss := "0"

	compareOp := strCompareNotStartsWithIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithHit,
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      startsWithMiss,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestNotStartsWithIgnoreCase", data, dataMapIntercept)
}

func TestEndsWithIgnoreCase(t *testing.T) {
	data := "xxXx"
	endsWithHit := "X"
	endsWithMiss := "op"

	compareOp := strCompareEndsWithIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: intercept,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: intercept,
		}: false,
	}
	commandRunCompare(t, "TestEndsWithIgnoreCase", data, dataMapIntercept)
}

func TestNotEndsWithIgnoreCase(t *testing.T) {
	data := "xxXx"
	endsWithHit := "xX"
	endsWithMiss := "OP"

	compareOp := strCompareNotEndsWithIgnoreCase
	dataMapIntercept := map[compareData]bool{
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: pass,
		}: true,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: pass,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithHit,
			interceptOrPass: intercept,
		}: false,
		{
			compare:         compareOp,
			targetData:      endsWithMiss,
			interceptOrPass: intercept,
		}: true,
	}
	commandRunCompare(t, "TestNotEndsWithIgnoreCase", data, dataMapIntercept)
}
