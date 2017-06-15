package cluster

import (
	"testing"
	"net"
	"strconv"
	"github.com/hashicorp/memberlist"
)

func TestBasicCluster(t *testing.T) {
	// Start three nodes
	nodes := []string{"127.0.0.1:7900", "127.0.0.1:7901", "127.0.0.1:7902"}

	servers := []*memberlist.Memberlist{}
	for _, node := range nodes {
		host, port, err := net.SplitHostPort(node)
		if err != nil {
			t.Fatal("net.SplitHostPort: ", err)
		}
		portI, _ :=  strconv.Atoi(port)
		server := Start(node, host, portI, nodes)
		servers = append(servers, server)
	}
	for _, server := range servers {
		members := server.Members()
		n := len(members)
		if n != len(nodes) {
			t.Fatal("Memberlist mismatch: ", n)
		}
	}
}
