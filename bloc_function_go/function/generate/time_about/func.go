package time_about

import (
	"context"
	"fmt"
	"time"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	var _ bloc_client.BlocFunctionNodeInterface = &TimeGen{}
}

type TimeGen struct {
}

func (tG *TimeGen) AllProgressMilestones() []string {
	return []string{}
}

func (tG *TimeGen) IptConfig() bloc_client.Ipts {
	return bloc_client.Ipts{
		{
			Key:     "choose",
			Display: "choose base time",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "what time to use",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   chooseSelections(),
					AllowMulti:      false,
				},
			},
		},
		{
			Key:     "input_time",
			Display: "input time",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "timestamp in second",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
				{
					Hint:            "string datetime. require match RFC3339 standard. e.g: '2006-01-02T15:04:05Z07:00'",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
				},
			},
		},
		{
			Key:     "time_offset",
			Display: "time offset",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "year",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
				{
					Hint:            "month",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
				{
					Hint:            "day",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
				{
					Hint:            "hour",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
				{
					Hint:            "minute",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
				{
					Hint:            "second",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.InputFormControl,
					AllowMulti:      false,
					DefaultValue:    0,
				},
			},
		},
	}
}

// OptConfig define the opt config of this function
func (tG *TimeGen) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{
		{
			Key:         "year",
			Description: "year of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "month",
			Description: "month of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "day",
			Description: "day of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "hour",
			Description: "hour of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "minute",
			Description: "minute of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "second",
			Description: "second of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "weekday",
			Description: "weekday of opt time, Sunday ... Saturday",
			ValueType:   bloc_client.StringValueType,
		},
		{
			Key:         "timestamp_in_second",
			Description: "timestamp in second of opt time",
			ValueType:   bloc_client.IntValueType,
		},
		{
			Key:         "datetime_str",
			Description: "datetime str of opt time",
			ValueType:   bloc_client.StringValueType,
		},
	}
}

// Run function's actual execute logic
func (tG *TimeGen) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	// parse params
	chooseIptValue, err := ipts.GetIntValue(0, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get msg content from ipt failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
		}
		return
	}
	chooseValue := choose(chooseIptValue)
	logger.Infof("get choose value: %d, label: %s", chooseIptValue, chooseValue.String())

	var baseTime time.Time
	if chooseValue == chooseNowAsBase {
		baseTime = time.Now()
	} else {
		inputTimeStampValue, _ := ipts.GetIntValue(1, 0)
		inputDatetimeStrValue, _ := ipts.GetStringValue(1, 1)
		if inputTimeStampValue <= 0 && inputDatetimeStrValue == "" {
			if chooseValue == chooseInterceptIfNotInput {
				blocOptChan <- &bloc_client.FunctionRunOpt{
					InterceptBelowFunctionRun: true,
					ErrorMsg:                  "no input time base value",
				}
				return
			}
			baseTime = time.Now()
		} else {
			if inputTimeStampValue > 0 {
				baseTime = time.Unix(int64(inputTimeStampValue), 0)
			} else if inputDatetimeStrValue != "" {
				tmp, err := time.ParseInLocation(time.RFC3339, inputDatetimeStrValue, time.Local)
				if err != nil {
					blocOptChan <- &bloc_client.FunctionRunOpt{
						InterceptBelowFunctionRun: true,
						ErrorMsg:                  fmt.Sprintf("parse input datetime str failed: %v", err),
					}
					return
				}
				baseTime = tmp
			}
		}
	}
	logger.Infof("base time: %s", baseTime.Format(time.RFC3339))

	timeOffsetYear, _ := ipts.GetIntValue(2, 0)
	timeOffsetMonth, _ := ipts.GetIntValue(2, 1)
	timeOffsetDay, _ := ipts.GetIntValue(2, 2)
	timeOffsetHour, _ := ipts.GetIntValue(2, 3)
	timeOffsetMinute, _ := ipts.GetIntValue(2, 4)
	timeOffsetSecond, _ := ipts.GetIntValue(2, 5)
	logger.Infof(
		"get time offset value - year: %d, month: %d, day: %d, hour: %d, minute: %d, second: %d",
		timeOffsetYear, timeOffsetMonth, timeOffsetDay, timeOffsetHour, timeOffsetMinute, timeOffsetSecond,
	)

	optDatetime := time.Date(
		baseTime.Year()+timeOffsetYear,
		time.Month(int(baseTime.Month())+timeOffsetMonth),
		baseTime.Day()+timeOffsetDay,
		baseTime.Hour()+timeOffsetHour,
		baseTime.Minute()+timeOffsetMinute,
		baseTime.Second()+timeOffsetSecond,
		baseTime.Nanosecond(),
		baseTime.Location())
	logger.Infof("output time: %s", optDatetime.Format(time.RFC3339))

	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc: true,
		Description: fmt.Sprintf(
			"input time: %s, opt time: %s",
			baseTime.Format(time.RFC3339), optDatetime.Format(time.RFC3339),
		),
		Detail: map[string]interface{}{
			"year":                optDatetime.Year(),
			"month":               optDatetime.Month(),
			"day":                 optDatetime.Day(),
			"hour":                optDatetime.Hour(),
			"minute":              optDatetime.Minute(),
			"second":              optDatetime.Second(),
			"weekday":             optDatetime.Weekday().String(),
			"timestamp_in_second": optDatetime.Unix(),
			"datetime_str":        optDatetime.Format(time.RFC3339),
		},
	}
}
