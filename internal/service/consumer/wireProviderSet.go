// generate by wireGenerate.sh with '^func New' in on package
package consumer

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMessageConsumer,
)
