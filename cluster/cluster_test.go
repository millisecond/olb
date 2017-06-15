package cluster

import (
	"testing"
	"net"
	"strconv"
	"github.com/hashicorp/memberlist"
	"github.com/facebookgo/ensure"
	"time"
)

func TestBasicCluster(t *testing.T) {
	// Start three nodes
	nodes := []string{"127.0.0.1:7900", "127.0.0.1:7901", "127.0.0.1:7902"}

	servers := []*memberlist.Memberlist{}
	for _, node := range nodes {
		host, port, err := net.SplitHostPort(node)
		ensure.Nil(t, err)
		portI, _ :=  strconv.Atoi(port)
		server := Start(node, host, portI, nodes)
		servers = append(servers, server)
	}
	for _, server := range servers {
		ensure.DeepEqual(t, len(server.Members()), len(nodes))
	}
	if !testing.Short() {
		err := servers[0].Shutdown()
		time.Sleep(time.Millisecond * 10000)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, len(servers[1].Members()), len(nodes)-1)
	}
}
