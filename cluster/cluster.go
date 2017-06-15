package cluster

import (
	"github.com/hashicorp/memberlist"
	"log"
)

func Start(name string, addr string, port int, others []string) *memberlist.Memberlist {
	list, err := memberlist.Create(config(name, addr, port))
	if err != nil {
		panic("Failed to init memberlist: " + err.Error())
	}
	n, err := list.Join(others)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
	log.Println("Joined MemberList cluster with nodecount: ", n)
	return list
}

func config(name string, addr string, port int) *memberlist.Config {
	conf := memberlist.DefaultLANConfig()
	conf.Name = name
	conf.BindAddr = addr
	conf.BindPort = port
	return conf
}
