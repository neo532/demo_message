// generate by wireGenerate.sh with '^func New' in on package
package api

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCampaignApi,
	NewMessageApi,
	NewSystemApi,
)
