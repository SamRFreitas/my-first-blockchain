package src

import (
	"errors"
	"fmt"
	"net/url"
)

type Network struct {
	Nodes map[url.URL]struct{} `json:"nodes"`
}

func NewNetwork() *Network {
	var network Network
	network.Nodes = make(map[url.URL]struct{})
	return &network
}

func (n *Network) Add(element url.URL) {
	n.Nodes[element] = struct{}{}
}

func (n *Network) Delete(element url.URL) error {
	if _, existis := n.Nodes[element]; !existis {
		return errors.New("element not present in set")
	}

	delete(n.Nodes, element)
	return nil
}

func (n *Network) Contains(element url.URL) bool {
	_, existis := n.Nodes[element]
	return existis
}

func (n *Network) List() {
	for k, _ := range n.Nodes {
		fmt.Println(k)
	}
}
