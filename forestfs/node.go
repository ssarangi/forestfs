package forestfs

import (
	"log"
	"sync"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	"github.com/ssarangi/forestfs/forest/config"
)

type Node struct {
	sync.RWMutex
	logger log.Logger

	// Raft members
	raft          *raft.Raft
	raftStore     *raftboltdb.BoltStore
	raftTransport *raft.NetworkTransport
	raftInmem     *raft.InmemStore
	// raftNotifyCh ensures we get reliable leader transition notifications from the raft layer
	raftNotifyCh <-chan bool
	// reconcileCh is used to pass events from the serf handler to the raft leader to update its state
	reconcileCh <-chan serf.Member
	serf        *serf.Serf
	fsm         *fsm.FSM
	eventChLan  chan serf.Event

	shutdownCh   chan struct{}
	shutdown     bool
	showdownLock sync.Mutex
}

func NewNode(config *config.NodeConfig, logger log.Logger) (*Node, error) {

}
