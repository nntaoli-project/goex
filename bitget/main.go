package bitget

import "github.com/nntaoli-project/goex/v2/bitget/spot"

type Bitget struct {
	Spot *spot.Spot
}

func New() *Bitget {
	return &Bitget{
		Spot: spot.New(),
	}
}
