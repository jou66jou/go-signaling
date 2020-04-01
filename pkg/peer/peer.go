package peer

import "github.com/gorilla/websocket"

// Peer 單一節點
type Peer struct {
	ID     string
	Socket *websocket.Conn
}

// SocketMessage ws 訊息格式
type SocketMessage struct {
	PeerID   string   `json:"peerID"`
	Type     string   `json:"type"`
	SDP      string   `json:"sdp"`
	OfferID  string   `json:"offerID"`
	FileName string   `json:"fileName"`
	FileList []string `json:"fileList"`
}

// Servicer peer 業務邏輯需求
type Servicer interface {
	SetPub()
	GetPub()
	SendAnswer()
	SendOffer()
}
