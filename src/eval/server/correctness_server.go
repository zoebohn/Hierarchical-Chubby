package main

import(
    "locks"
	"raft"
	"os"
    "os/signal"
)

func main() {
    addrs := []raft.ServerAddress{"127.0.0.1:40000", "127.0.0.1:40001", "127.0.0.1:40002"}
   var recruitAddrs [][]raft.ServerAddress = [][]raft.ServerAddress{{"127.0.0.1:50000", "127.0.0.1:50001", "127.0.0.1:50002"}, {"127.0.0.1:50003", "127.0.0.1:50004", "127.0.0.1:50005"}, {"127.0.0.1:50006", "127.0.0.1:50007", "127.0.0.1:50008"}, {"127.0.0.1:50009", "127.0.0.1:50010", "127.0.0.1:50011"}, {"127.0.0.1:50012", "127.0.0.1:50013", "127.0.0.1:50014"}, {"127.0.0.1:50015", "127.0.0.1:50016", "127.0.0.1:50017"}, {"127.0.0.1:50018", "127.0.0.1:50019", "127.0.0.1:50020"}} 
    locks.MakeCluster(3, locks.CreateMasters(len(addrs), addrs, recruitAddrs, 10, true), addrs)
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
