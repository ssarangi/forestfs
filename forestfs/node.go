package forestfs

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	"github.com/ssarangi/forestfs/forestfs/config"
)

var (
	ErrInvalidArgument = errors.New("No Logger set")
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
	node := &Node{
		config:      config,
		logger:      logger.With(log.Int32("id", config.ID), log.String("raft addr", config.RaftAddr)),
		shutdownCh:  make(chan struct{}),
		eventChLan:  make(chan serf.Event, 256),
		reconcileCh: make(chan serf.Member, 32),
	}

	if node.logger == nil {
		return nil, ErrInvalidArgument
	}

	node.logger.Info("Test Log....")
	if err := node.setupRaft(); err != nil {
		node.Shutdown()
		return nil, fmt.Errorf("Failed to start Raft: %v", err)
	}

	var err Error
	node.serf, err := node.setupSerf(config.SerfLANConfig, node.eventChLan, serfLANSnapshot)
	if err != nil {
		return nil, err
	}

	go node.lanEventHandler()
	go node.monitorLeadership()
	return node, nil
}

// Runs a loop to handle requests send back responses
func (node *Node) Run(ctx context.Context, requestc <-chan Request, responsec chan<- Response) {

}
