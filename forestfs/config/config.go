package config

import (
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	DefaultLANSerfPort = 8301
)

type NodeConfig struct {
	ID                int32
	NodeName          string
	DataDir           string
	DevMode           bool
	Addr              string
	SerfLANConfig     *serf.Config
	RaftConfig        *raft.Config
	Bootstrap         bool
	BootstrapExpect   int
	StartAsLeader     bool
	StartJoinAddrsLAN []string
	StartJoinAddrsWAN []string
	NonVoter          bool
	RaftAddr          string
	LeaveDrainTime    time.Duration
	ReconcileInterval time.Duration
}

func DefaultNodeConfig() *NodeConfig {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	conf := &NodeConfig{
		DevMode:           false,
		NodeName:          hostname,
		SerfLANConfig:     serfDefaultConfig(),
		RaftConfig:        raft.DefaultConfig(),
		LeaveDrainTime:    5 * time.Second,
		ReconcileInterval: 60 * time.Second,
	}

	conf.SerfLANConfig.ReconnectTimeout = 3 * 24 * time.Hour
	conf.SerfLANConfig.MemberlistConfig.BindPort = DefaultLANSerfPort

	return conf
}

func serfDefaultConfig() *serf.Config {
	base := serf.DefaultConfig()
	base.QueueDepthWarning = 1000000
	return base
}

type ServerConfig struct {
	NodeAddr string
}
