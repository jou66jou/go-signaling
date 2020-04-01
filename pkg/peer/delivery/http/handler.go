package http

import (
	"encoding/json"
	"fmt"
	peer "p2p/pkg/peer"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

// handler http server 控制物件
type handler struct {
	peerMgr *peer.Manager
	peerSvc peer.Servicer
}

// NewHandler http server 控制物件
func NewHandler(peerSvc peer.Servicer, peerMgr *peer.Manager) *handler {
	return &handler{
		peerSvc: peerSvc,
		peerMgr: peerMgr,
	}
}

func (h *handler) webSocketEndpoint(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	id, _ := uuid.NewV4()
	p := &peer.Peer{ID: id.String(), Socket: conn}
	h.peerMgr.Register <- p
	defer func() {
		h.peerMgr.UnRegister <- p
		conn.Close()
	}()

	for {
		// // Write
		// err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		// if err != nil {
		// 	c.Logger().Error(err)
		// }

		// Read
		_, msg, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}
		fmt.Printf("%+v\n", string(msg))
		var wsmsg = peer.SocketMessage{}
		err = json.Unmarshal(msg, &wsmsg)
		if err != nil {
			fmt.Printf("can't difine ws msg : %s \n", msg)
			continue
		}
		switch wsmsg.Type {
		case "offer":
			wsmsg.OfferID = p.ID
			if h.peerMgr.PubPeer != nil {
				err = h.peerMgr.PubPeer.Socket.WriteJSON(&wsmsg)
				if err != nil {
					fmt.Printf("json marshal get err : %+v\n", err)
					continue
				}
			}
		case "answer":
			if _, ok := h.peerMgr.Peers[wsmsg.OfferID]; ok {
				err = h.peerMgr.Peers[wsmsg.OfferID].Socket.WriteJSON(&wsmsg)
				if err != nil {
					fmt.Printf("json marshal get err : %+v\n", err)
					continue
				}
			}
		case "pub":
			h.peerMgr.PubPeer = p
		default:
			fmt.Printf("unknow type : %+v\n", wsmsg)
		}
	}
	return nil
}
