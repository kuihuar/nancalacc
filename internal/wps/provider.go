package wps

import (
	"github.com/google/wire"
)

var WpsProviderSet = wire.NewSet(NewWps)
