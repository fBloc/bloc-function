package time_about

import (
	bloc_client "github.com/fBloc/bloc-client-go"
)

type choose int

const (
	chooseInputFirstOrNowAsBase choose = iota
	chooseNowAsBase
	chooseInterceptIfNotInput
	max
)

func (c choose) String() string {
	switch c {
	case chooseNowAsBase:
		return "use now as base"
	case chooseInputFirstOrNowAsBase:
		return "use input time as base first、or now as base"
	case chooseInterceptIfNotInput:
		return "use input time as base first、or intercept below run"
	}
	return "unknown"
}

func chooseSelections() []bloc_client.SelectOption {
	ret := make([]bloc_client.SelectOption, 0, max-1)
	for i := 0; i < int(max); i++ {
		ret = append(
			ret,
			bloc_client.SelectOption{
				Label: choose(i).String(),
				Value: i,
			},
		)
	}
	return ret
}
