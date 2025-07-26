package wps

import (
	"github.com/google/wire"
)

var WpsProviderSet = wire.NewSet(NewWpsSync, NewWps)
