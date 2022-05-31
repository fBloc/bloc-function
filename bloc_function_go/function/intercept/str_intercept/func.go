package str_intercept

import (
	"context"
	"fmt"
	"strings"

	bloc_client "github.com/fBloc/bloc-client-go"
)

func init() {
	var _ bloc_client.BlocFunctionNodeInterface = &StrIntercept{}
}

type StrIntercept struct {
}

func (sI *StrIntercept) AllProgressMilestones() []string {
	return []string{}
}

func (sI *StrIntercept) IptConfig() bloc_client.Ipts {
	return bloc_client.Ipts{
		{
			Key:     "source_string",
			Display: "source string",
			Must:    true,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
				},
			},
		},
		{
			Key:     "blank_check",
			Display: "blank string check",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "condition",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   blankSelections(),
				},
				{
					Hint:            "intercept_or_pass",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   interceptSelections(),
				},
			},
		},
		{
			Key:     "intercept_condition",
			Display: "intercept condition",
			Must:    false,
			Components: []*bloc_client.IptComponent{
				{
					Hint:            "compare_type",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   strCompareSelections(),
				},
				{
					Hint:            "value",
					ValueType:       bloc_client.StringValueType,
					FormControlType: bloc_client.InputFormControl,
				},
				{
					Hint:            "intercept_or_pass",
					ValueType:       bloc_client.IntValueType,
					FormControlType: bloc_client.SelectFormControl,
					SelectOptions:   interceptSelections(),
				},
			},
		},
	}
}

func (sI *StrIntercept) OptConfig() bloc_client.Opts {
	return bloc_client.Opts{}
}

func (sI *StrIntercept) Run(
	ctx context.Context,
	ipts bloc_client.Ipts,
	progressReportChan chan bloc_client.HighReadableFunctionRunProgress,
	blocOptChan chan *bloc_client.FunctionRunOpt,
	logger *bloc_client.Logger,
) {
	iptStr, err := ipts.GetStringValue(0, 0)
	if err != nil {
		errorMsg := fmt.Sprintf("get ipt string failed: %v", err)
		logger.Warningf(errorMsg)
		blocOptChan <- &bloc_client.FunctionRunOpt{
			InterceptBelowFunctionRun: true,
			ErrorMsg:                  errorMsg,
		}
		return
	}
	iptStrLower := strings.ToLower(iptStr)
	logger.Infof("get input string data: %s", iptStr)

	blankCheckIntVal, _ := ipts.GetIntValue(1, 0)
	blankCheckVal := blank(blankCheckIntVal)
	blankCheckInterceptOrPassIntVal, _ := ipts.GetIntValue(1, 1)
	blankCheckInterceptOrPassVal := interceptOrPass(blankCheckInterceptOrPassIntVal)
	if blankCheckVal.IsValid() && blankCheckInterceptOrPassVal.IsValid() {
		if blankCheckVal == isBlank {
			if iptStr == "" && blankCheckInterceptOrPassVal == intercept { // 拦截空
				blocOptChan <- &bloc_client.FunctionRunOpt{
					Suc:                       true,
					InterceptBelowFunctionRun: true,
					Description:               "intercept blank string",
				}
				return
			}
			if iptStr != "" && blankCheckInterceptOrPassVal == pass { // 放过空
				blocOptChan <- &bloc_client.FunctionRunOpt{
					Suc:                       true,
					InterceptBelowFunctionRun: true,
					Description:               "only pass blank string",
				}
				return
			}
		} else {
			if iptStr != "" && blankCheckInterceptOrPassVal == intercept { // 拦截非空
				blocOptChan <- &bloc_client.FunctionRunOpt{
					Suc:                       true,
					InterceptBelowFunctionRun: true,
					Description:               "intercept not blank string",
				}
				return
			}
			if iptStr == "" && blankCheckInterceptOrPassVal == pass { // 放过非空（拦截空）
				blocOptChan <- &bloc_client.FunctionRunOpt{
					Suc:                       true,
					InterceptBelowFunctionRun: true,
					Description:               "intercept blank string",
				}
				return
			}
		}
	}

	compareTypeIntVal, _ := ipts.GetIntValue(2, 0)
	compareTypeVal := strCompare(compareTypeIntVal)
	compareTargetStrVal, _ := ipts.GetStringValue(2, 1)
	compareTargetStrValLower := strings.ToLower(compareTargetStrVal)
	compareInterceptOrPassIntVal, _ := ipts.GetIntValue(2, 2)
	compareInterceptOrPassVal := interceptOrPass(compareInterceptOrPassIntVal)

	intercepted := false

	if iptStr != "" && compareTypeVal.IsValid() && compareInterceptOrPassVal.IsValid() {
		switch compareTypeVal {
		case strCompareEq:
			if iptStr == compareTargetStrVal && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if iptStr != compareTargetStrVal && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotEq:
			if iptStr != compareTargetStrVal && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if iptStr == compareTargetStrVal && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareContains:
			if strings.Contains(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.Contains(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotContains:
			if !strings.Contains(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.Contains(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareStartsWith:
			if strings.HasPrefix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.HasPrefix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotStartsWith:
			if !strings.HasPrefix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.HasPrefix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareEndsWith:
			if strings.HasSuffix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.HasSuffix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotEndsWith:
			if !strings.HasSuffix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.HasSuffix(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareEqIgnoreCase:
			if strings.EqualFold(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.EqualFold(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotEqIgnoreCase:
			if !strings.EqualFold(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.EqualFold(iptStr, compareTargetStrVal) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareContainsIgnoreCase:
			if strings.Contains(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.Contains(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotContainsIgnoreCase:
			if !strings.Contains(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.Contains(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareStartsWithIgnoreCase:
			if strings.HasPrefix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.HasPrefix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotStartsWithIgnoreCase:
			if !strings.HasPrefix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.HasPrefix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareEndsWithIgnoreCase:
			if strings.HasSuffix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if !strings.HasSuffix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		case strCompareNotEndsWithIgnoreCase:
			if !strings.HasSuffix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == intercept {
				intercepted = true
			}
			if strings.HasSuffix(iptStrLower, compareTargetStrValLower) && compareInterceptOrPassVal == pass {
				intercepted = true
			}
		}
	}

	if intercepted {
		var msg string
		if compareInterceptOrPassVal == intercept {
			msg = fmt.Sprintf(
				"intercept because not allowed %s %s %s",
				iptStr, compareTypeVal.String(), compareTargetStrVal)
		} else {
			msg = fmt.Sprintf(
				"intercept because need %s %s %s",
				iptStr, compareTypeVal.String(), compareTargetStrVal)
		}
		blocOptChan <- &bloc_client.FunctionRunOpt{
			Suc:                       true,
			InterceptBelowFunctionRun: true,
			Description:               msg,
		}
		return
	}
	blocOptChan <- &bloc_client.FunctionRunOpt{
		Suc: true,
	}
}
