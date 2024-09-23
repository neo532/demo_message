// generate by wireGenerate.sh with '^func New' in on package
package biz

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCampaignUsecase,
	NewMessageUsecase,
)
