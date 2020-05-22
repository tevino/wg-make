package rendering

import (
	"time"

	"github.com/tevino/wg-make/config"
)

// PeerConfigTplContext contains context for peer configuration file rendering.
type PeerConfigTplContext struct {
	Network     *config.Network
	GeneratedAt time.Time
	Interface   *config.Peer
	Peers       []config.Peer
}
