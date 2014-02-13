package main

type NodePool struct {
	dht   *DHT
	cache map[string]*Node
}

func NewNodePool(d *DHT) *NodePool {
	return &NodePool{
		dht:   d,
		cache: make(map[string]*Node),
	}
}

func (np *NodePool) GetNode(id *KeyID, globalAddress string) *Node {
	nodeHash := id.String()

	result, ok := np.cache[nodeHash]

	if !ok {
		result = NewNode(id, globalAddress)
		result.Connect(np.dht) // use dht perhaps here?

		np.cache[nodeHash] = result

    println("nodes in the pool:", len(np.cache))
	}

	return result
}

func (np *NodePool) PutNode(node *Node) *Node {
	nodeHash := node.Id().String()

	if _, ok := np.cache[nodeHash]; !ok {
		np.cache[nodeHash] = node
    println("nodes in the pool:", len(np.cache))
	}

	return node
}
