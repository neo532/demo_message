// generate by wireGenerate.sh with '^func New' in on package
package script

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMessageScript,
)
