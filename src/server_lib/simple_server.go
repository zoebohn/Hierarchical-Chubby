package main

import(
	"fmt"
	"raft"
	"io"
	"io/ioutil"
	"time"
	"log"
	"os"
    "os/signal"
)

type simpleFSM struct{
}

type simpleSnapshot struct{
}

func (s *simpleFSM) Apply(log *raft.Log) interface{} {
	fmt.Println(log.Data)
	return nil
}

func (s *simpleFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &simpleSnapshot{}, nil
}

func (s *simpleFSM) Restore(io.ReadCloser) error {
	return nil
}

func (s *simpleSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s *simpleSnapshot) Release() {
}

func main() {
    makeCluster(3)
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func makeCluster(n int) *cluster {
    conf := raft.DefaultConfig()
    bootstrap := true

	c := &cluster{
		conf:          conf,
		// Propagation takes a maximum of 2 heartbeat timeouts (time to
		// get a new heartbeat that would cause a commit) plus a bit.
		propagateTimeout: conf.HeartbeatTimeout*2 + conf.CommitTimeout,
		longstopTimeout:  5 * time.Second,
	}
	var configuration raft.Configuration

	// Setup the stores and transports
	for i := 0; i < n; i++ {
		dir, err := ioutil.TempDir("", "raft")
		if err != nil {
			fmt.Println("[ERR] err: %v ", err)
		}

		store := raft.NewInmemStore()
		c.dirs = append(c.dirs, dir)
		c.stores = append(c.stores, store)
		c.fsms = append(c.fsms, &simpleFSM{})


	    snap, err := raft.NewFileSnapshotStore(dir, 3, nil)
		c.snaps = append(c.snaps, snap)

        trans, err := raft.NewTCPTransport("127.0.0.1:0", nil, 2, time.Second, nil)
        if err != nil {
            fmt.Println("err: %v", err)
        }
        addr := trans.LocalAddr()
        c.trans = append(c.trans, trans)
		localID := raft.ServerID(fmt.Sprintf("server-%s", addr))
		configuration.Servers = append(configuration.Servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       localID,
			Address:  addr,
		})
	}

	// Create all the rafts
	c.startTime = time.Now()
	for i := 0; i < n; i++ {
		logs := c.stores[i]
		store := c.stores[i]
		snap := c.snaps[i]
		trans := c.trans[i]

		peerConf := conf
		peerConf.LocalID = configuration.Servers[i].ID
        peerConf.Logger = log.New(os.Stdout, string(peerConf.LocalID) + " : ", log.Lmicroseconds)

		if bootstrap {
			err := raft.BootstrapCluster(peerConf, logs, store, snap, trans, configuration)
			if err != nil {
				fmt.Println("[ERR] BootstrapCluster failed: %v", err)
			}
		}

		raft, err := raft.NewRaft(peerConf, c.fsms[i], logs, store, snap, trans)
		if err != nil {
		    fmt.Println("[ERR] NewRaft failed: %v", err)
		}

		raft.AddVoter(peerConf.LocalID, trans.LocalAddr(), 0, 0)
		c.rafts = append(c.rafts, raft)
	}

	return c
}

type cluster struct {
	dirs             []string
	stores           []*raft.InmemStore
	fsms             []*simpleFSM
	snaps            []*raft.FileSnapshotStore
	trans            []raft.Transport
	rafts            []*raft.Raft
	conf             *raft.Config
	propagateTimeout time.Duration
	longstopTimeout  time.Duration
	startTime        time.Time
}
