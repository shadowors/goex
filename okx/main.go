package okx

import (
	"github.com/shadowors/goex/v2/okx/common"
	"github.com/shadowors/goex/v2/okx/futures"
	"github.com/shadowors/goex/v2/okx/spot"
)

type OKx struct {
	Spot    *spot.Spot
	Futures *futures.Futures
	Swap    *futures.Swap
	Asset   *common.OKxV5
}

func New() *OKx {
	okxV5 := spot.New().OKxV5
	return &OKx{
		Spot:    spot.New(),
		Futures: futures.New(),
		Swap:    futures.NewSwap(),
		Asset:   okxV5,
	}
}
