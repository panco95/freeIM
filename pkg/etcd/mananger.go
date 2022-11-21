package etcd

import (
	"context"
	"im/pkg/utils"
	"log"
	"strings"
	"time"

	clientV3 "go.etcd.io/etcd/client/v3"
)

type Manager struct {
	etcdPrefix string
	etcd       *clientV3.Client
	nodeChan   chan nodeOperate
	nodes      []*node
	local      *node
}

func NewManager(addr string, etcd *clientV3.Client) (*Manager, error) {
	var err error

	if addr == "" {
		addr, err = utils.GetOutboundIP()
		if err != nil {
			return nil, err
		}
	}

	m := &Manager{
		etcd:  etcd,
		nodes: make([]*node, 0),
		local: &node{addr: addr},
	}

	m.nodeChan = make(chan nodeOperate)
	go m.nodesWatch()

	if err = m.nodeRegister(true); err != nil {
		return nil, err
	}

	return m, nil
}

type node struct {
	addr string
}

type nodeOperate struct {
	operate string
	addr    string
}

func (m *Manager) GetNodes() []*node {
	return m.nodes
}

func (m *Manager) GetLocalIp() string {
	return m.local.addr
}

func (m *Manager) GetLocalId() string {
	return m.etcdPrefix + "_" + m.local.addr
}

func (m *Manager) nodeRegister(isReconnect bool) error {
	client := m.etcd

	// New lease
	resp, err := client.Grant(context.TODO(), 2)
	if err != nil {
		return err
	}
	// The lease was granted
	if err != nil {
		return err
	}
	_, err = client.Put(context.TODO(), m.GetLocalId(), "0", clientV3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// keep alive
	ch, err := client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return err
	}
	// monitor etcd connection
	go func() {
		for {
			select {
			case resp := <-ch:
				if resp == nil {
					go m.nodeRegister(false)
					return
				}
			}
		}
	}()

	if isReconnect {
		go m.serviceWatcher()
		go func() {
			for {
				m.getAllServices()
				time.Sleep(time.Second * 5)
			}
		}()
	}
	return nil
}

func (m *Manager) serviceWatcher() {
	client := m.etcd

	rch := client.Watch(context.Background(), m.etcdPrefix+"_", clientV3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			arr := strings.Split(string(ev.Kv.Key), "_")
			addr := arr[1]
			switch ev.Type {
			case 0:
				m.addNode(addr)
				log.Printf("add %s", addr)
			case 1:
				m.delNode(addr)
				log.Printf("del %s", addr)
			}
		}
	}
}

func (m *Manager) getAllServices() ([]string, error) {
	client := m.etcd

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := client.Get(ctx, m.etcdPrefix+"_", clientV3.WithPrefix())
	cancel()
	if err != nil {
		return []string{}, nil
	}

	nodes := make([]string, 0)
	for _, ev := range resp.Kvs {
		arr := strings.Split(string(ev.Key), m.etcdPrefix+"_")
		addr := arr[1]
		m.addNode(addr)
	}

	return nodes, nil
}

func (m *Manager) addNode(addr string) {
	c := nodeOperate{
		operate: "addNode",
		addr:    addr,
	}
	m.nodeChan <- c
}

func (m *Manager) delNode(addr string) {
	c := nodeOperate{
		operate: "delNode",
		addr:    addr,
	}
	m.nodeChan <- c
}

func (m *Manager) nodesWatch() {
	for {
		select {
		case c := <-m.nodeChan:
			switch c.operate {
			case "addNode":
				m.nodes = append(m.nodes, &node{addr: c.addr})
			case "delNode":
				for i := 0; i < len(m.nodes); i++ {
					if m.nodes[i].addr == c.addr {
						m.nodes = append(m.nodes[:i], m.nodes[i+1:]...)
						i--
					}
				}
			}
		}
	}
}
