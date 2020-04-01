package peer

import "fmt"

// Manager 節點管理器
type Manager struct {
	Peers      map[string]*Peer
	FileMap    map[string][]*Peer
	PubPeer    *Peer
	Register   chan *Peer
	UnRegister chan *Peer
}

// InitManager 初始化節點管理器
func InitManager() *Manager {
	var mgr = &Manager{
		Register:   make(chan *Peer),
		UnRegister: make(chan *Peer),
		Peers:      make(map[string]*Peer),
		FileMap:    make(map[string][]*Peer),
	}
	go mgr.start()
	return mgr
}

// start 監聽節點事件
func (Manager *Manager) start() {
	for {
		select {
		case peer := <-Manager.Register:
			Manager.Peers[peer.ID] = peer

		case peer := <-Manager.UnRegister:
			if _, ok := Manager.Peers[peer.ID]; ok {
				fmt.Printf("%+v unRegister peer\n", peer.ID)
				delete(Manager.Peers, peer.ID)
				// TODO del fileMap
			}
			// case message := <-Manager.broadcast:
			// 	for peer := range Manager.Peers {
			// 		select {
			// 		case peer.send <- message:
			// 		}
			// 	}
		}
	}
}
