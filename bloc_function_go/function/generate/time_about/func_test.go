package time_about

import (
	"testing"
	"time"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func TestBaseNow(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&TimeGen{},
		[][]interface{}{
			{
				0,
			},
			{
				0,
				"",
			},
			{
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept")
	}
}

func TestBaseNowWithOffset(t *testing.T) {
	client := bloc_client.NewTestClient()

	now := time.Now()
	executeOpt := client.TestRunFunction(
		&TimeGen{},
		[][]interface{}{
			{
				0,
			},
			{
				0,
				"",
			},
			{
				-1,
				-1,
				-1,
				-1,
				-1,
				-1,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept")
	}
	optDatimeTimeFieldStr, ok := executeOpt.Detail["datetime_str"].(string)
	if !ok {
		t.Fatalf("opt['datetime_str'] parse to string failed")
	}
	optDatimeTimeFromStr, err := time.Parse(time.RFC3339, optDatimeTimeFieldStr)
	if err != nil {
		t.Fatalf("parse output detail field datetime_str to time failed")
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	optTimeStampInt, ok := executeOpt.Detail["timestamp_in_second"].(int64)
	if !ok {
		t.Fatalf("opt['timestamp_in_second'] parse to int64 failed")
	}
	optDatimeTimeFromTimeStamp := time.Unix(optTimeStampInt, 0)
	if !optDatimeTimeFromTimeStamp.Equal(optDatimeTimeFromStr) {
		t.Fatalf("time should be equal, but %s != %s", optDatimeTimeFromTimeStamp, optDatimeTimeFromStr)
	}
}

func TestShouldIntercept(t *testing.T) {
	client := bloc_client.NewTestClient()

	executeOpt := client.TestRunFunction(
		&TimeGen{},
		[][]interface{}{
			{
				2, // chooseInterceptIfNotInput
			},
			{ // no valid input time
				0,
				"",
			},
			{
				-1,
				-1,
				-1,
				-1,
				-1,
				-1,
			},
		},
	)
	if executeOpt.Suc {
		t.Fatalf("should not success")
	}
	if !executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should intercept")
	}
}

func TestInputBaseTimeWithTimestampOffset(t *testing.T) {
	client := bloc_client.NewTestClient()

	now := time.Now()
	executeOpt := client.TestRunFunction(
		&TimeGen{},
		[][]interface{}{
			{
				2,
			},
			{
				now.Unix(),
				"",
			},
			{
				-1,
				-1,
				-1,
				-1,
				-1,
				-1,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept")
	}
	optDatimeTimeFieldStr, ok := executeOpt.Detail["datetime_str"].(string)
	if !ok {
		t.Fatalf("opt['datetime_str'] parse to string failed")
	}
	optDatimeTimeFromStr, err := time.Parse(time.RFC3339, optDatimeTimeFieldStr)
	if err != nil {
		t.Fatalf("parse output detail field datetime_str to time failed")
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	optTimeStampInt, ok := executeOpt.Detail["timestamp_in_second"].(int64)
	if !ok {
		t.Fatalf("opt['timestamp_in_second'] parse to int64 failed")
	}
	optDatimeTimeFromTimeStamp := time.Unix(optTimeStampInt, 0)
	if !optDatimeTimeFromTimeStamp.Equal(optDatimeTimeFromStr) {
		t.Fatalf("time should be equal, but %s != %s", optDatimeTimeFromTimeStamp, optDatimeTimeFromStr)
	}
}

func TestInputBaseTimeWithDatetimeStrOffset(t *testing.T) {
	client := bloc_client.NewTestClient()

	now := time.Now()
	executeOpt := client.TestRunFunction(
		&TimeGen{},
		[][]interface{}{
			{
				2,
			},
			{
				0,
				now.Format(time.RFC3339),
			},
			{
				-1,
				-1,
				-1,
				-1,
				-1,
				-1,
			},
		},
	)
	if !executeOpt.Suc {
		t.Fatalf("should success, but failed with %s", executeOpt.ErrorMsg)
	}
	if executeOpt.InterceptBelowFunctionRun {
		t.Fatalf("should not intercept")
	}
	optDatimeTimeFieldStr, ok := executeOpt.Detail["datetime_str"].(string)
	if !ok {
		t.Fatalf("opt['datetime_str'] parse to string failed")
	}
	optDatimeTimeFromStr, err := time.Parse(time.RFC3339, optDatimeTimeFieldStr)
	if err != nil {
		t.Fatalf("parse output detail field datetime_str to time failed")
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	if (now.Year() - optDatimeTimeFromStr.Year()) != 1 {
		t.Fatalf("year offset should be 1, but %d", now.Year()-optDatimeTimeFromStr.Year())
	}
	optTimeStampInt, ok := executeOpt.Detail["timestamp_in_second"].(int64)
	if !ok {
		t.Fatalf("opt['timestamp_in_second'] parse to int64 failed")
	}
	optDatimeTimeFromTimeStamp := time.Unix(optTimeStampInt, 0)
	if !optDatimeTimeFromTimeStamp.Equal(optDatimeTimeFromStr) {
		t.Fatalf("time should be equal, but %s != %s", optDatimeTimeFromTimeStamp, optDatimeTimeFromStr)
	}
}
