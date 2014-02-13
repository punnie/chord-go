package main

import (
	"net"
	"time"
)

const (
	BITS           = 256 // sha1
	WORKER_THREADS = 4
)

type DHT struct {
	self          *Node
	finger        []*Node
	predecessor   *Node
	successor     *Node
	globalInbound chan *Envelope
	nodePool      *NodePool
	ticker        *time.Ticker
	debugTicker   *time.Ticker
}

func NewDHT(self *Node) *DHT {
	dht := &DHT{
		self:          self,
		finger:        make([]*Node, BITS),
		predecessor:   nil,
		successor:     self,
		globalInbound: make(chan *Envelope, 100),
		ticker:        time.NewTicker(10 * time.Second),
		debugTicker:   time.NewTicker(5 * time.Second),
	}

	dht.nodePool = NewNodePool(dht)

	// TODO: implement worker number
	// also: make threadsafe, which _isn't_
	for i := 0; i < WORKER_THREADS; i++ {
		go dht.GlobalInboundWorker()
	}

	return dht
}

func (d *DHT) Store(object []byte) error {
	return nil
}

func (d *DHT) Retrieve(id int64) ([]byte, error) {
	return nil, nil
}

func (d *DHT) Join(node *Node) {
	err := node.Connect(d)

	if err != nil {
		panic(err)
	}

	d.successor, err = node.RequestSuccessor(d.self.Id()) // blocks
	d.successor.Connect(d)
	d.stabilize()
}

func (d *DHT) Listen() {
	sock, err := net.Listen("tcp", d.self.GlobalAddress())

	if err != nil {
		panic(err)
	}

	println("Listening on", d.self.GlobalAddress())

	for {
		conn, err := sock.Accept()

		if err != nil {
			println("Error accepting!")
		}

		node := NewNode(nil, "")
		node.Accept(d, conn)
	}
}

func (d *DHT) GlobalInboundWorker() {
	for {
		select {
		case e := <-d.globalInbound:
			m := e.message
      //println("receiving      :", m.String())

			switch m.Intent { // perhaps make functions out of this
			case _M_REQUEST_SUCCESSOR:
				queryKey := NewKeyID().SetHash(m.Parameters[0])
        println("requested successor of:", queryKey.String())
				replyNode, err := d.findSuccessor(queryKey)

				if err != nil {
					panic(err)
				}

				e.sender.ReplySuccessor(replyNode)

			case _M_REQUEST_PREDECESSOR:
				e.sender.ReplyPredecessor(d.predecessor)

			case _M_NOTIFY:
				d.notify(e.sender)

			case _M_REQUEST_PING:
				e.sender.ReplyPing()
			}

		case <-d.ticker.C:
			d.stabilize()
      d.checkPredecessor()

		case <-d.debugTicker.C:
      println("successor:", d.successor.Id().String())

      if d.predecessor != nil {
        println("predecessor:", d.predecessor.Id().String())
      }
		}
	}
}

//
//
//
//
//

func (d *DHT) findSuccessor(id *KeyID) (*Node, error) {
	if id.elementOf(d.self.Id(), d.successor.Id()) { // this interval is (] / correct
		return d.successor, nil
	} else {
		//queryNode := d.closestPrecedingNode(id)
    queryNode := d.successor
		resultNode, err := queryNode.RequestSuccessor(id)

		if err != nil {
			return nil, err
		}

		return resultNode, nil
	}
}

func (d *DHT) closestPrecedingNode(id *KeyID) *Node {
	for i := BITS; i > 0; i-- {
		if d.finger[i].Id().elementOf(d.self.Id(), id) { // this interval is () / false
			return d.finger[i]
		}
	}

	return d.self
}

func (d *DHT) stabilize() {
	if d.successor == d.self { // if this is the node that created the dht
		node := d.predecessor

		if node != nil && node.Id().elementOf(d.self.Id(), d.successor.Id()) {
			d.successor = node
			d.successor.Notify(node)
		}
	} else {
		node, err := d.successor.RequestPredecessor()

		if err != nil {
			panic(err)
		}

		if node != nil && node.Id().elementOf(d.self.Id(), d.successor.Id()) {
			d.successor = node
		}

		d.successor.Notify(node)
	}
}

func (d *DHT) notify(node *Node) {
	if d.predecessor == nil || node.Id().elementOf(d.predecessor.Id(), d.self.Id()) { // this interval is () / false
		d.predecessor = node
	}
}

func (d *DHT) fixFingers() {
}

func (d *DHT) checkPredecessor() {
  if d.predecessor != nil {
    err := d.predecessor.SendPing()

    if err != nil {
      d.predecessor = nil
    }
  }
}
