package dingtalk

import (
	"github.com/google/wire"
)

var DingtalkProviderSet = wire.NewSet(NewDingTalkRepo)
