package metricplugin

import (
	"time"

	bsmsg "github.com/ipfs/go-bitswap/message"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

// A BitswapMessage is the type pushed to remote clients for recorded incoming
// Bitswap messages.
type BitswapMessage struct {
	// Wantlist entries sent with this message.
	WantlistEntries []bsmsg.Entry `json:"wantlist_entries"`

	// Whether the wantlist entries are a full new wantlist.
	FullWantList bool `json:"full_wantlist"`

	// Blocks sent with this message.
	Blocks []cid.Cid `json:"blocks"`

	// Block presence indicators sent with this message.
	BlockPresences []BlockPresence `json:"block_presences"`
}

// A BlockPresence indicates the presence or absence of a block.
type BlockPresence struct {
	Cid  cid.Cid           `json:"cid"`
	Type BlockPresenceType `json:"block_presence_type"`
}

// BlockPresenceType is an enum for presence or absence notifications.
type BlockPresenceType int

const (
	// Have indicates that the peer has the block.
	Have BlockPresenceType = 0
	// DontHave indicates that the peer does not have the block.
	DontHave BlockPresenceType = 1
)

// ConnectionEventType specifies the type of connection event.
type ConnectionEventType int

const (
	// Connected specifies that a connection was opened.
	Connected ConnectionEventType = 0
	// Disconnected specifies that a connection was closed.
	Disconnected ConnectionEventType = 1
)

// A ConnectionEvent is the type pushed to remote clients for recorded
// connection events.
type ConnectionEvent struct {
	// The multiaddress of the remote peer.
	Remote ma.Multiaddr `json:"remote"`

	// The type of this event.
	ConnectionEventType ConnectionEventType `json:"connection_event_type"`
}

// An EventSubscriber can handle events generated by the plugin.
type EventSubscriber interface {
	// ID should return a unique identifier.
	// The identifier is used to keep track of subscribers.
	// The identifier may be reused, but only after the old use has been
	// unsubscribed.
	ID() string

	// BitswapMessageReceived handles a Bitswap message that was recorded by the
	// plugin.
	// TODO for now this blocks, maybe we have to figure something out.
	BitswapMessageReceived(timestamp time.Time, peer peer.ID, msg BitswapMessage)

	// ConnectionEventRecorded handles a connection event that was recorded by
	// the plugin.
	// TODO for now this blocks, maybe we have to figure something out.
	ConnectionEventRecorded(timestamp time.Time, peer peer.ID, connEvent ConnectionEvent)
}

// ErrAlreadySubscribed is returned by Subscribe if the given EventSubscriber is
// already subscribed.
var ErrAlreadySubscribed = errors.New("already subscribed")

// PluginAPI describes the functionality provided by this plugin to remote
// clients.
type PluginAPI interface {
	// Subscribe adds a subscriber to the event subscription service.
	// Returns ErrAlreadySubscribed if the given subscriber is already subscribed.
	Subscribe(subscriber EventSubscriber) error

	// Unsubscribe removes a subscriber from the event subscription service.
	// It is safe to call this multiple times with the same subscriber.
	Unsubscribe(subscriber EventSubscriber)

	// Ping is a no-op.
	Ping()

	// TODO additional methods
}