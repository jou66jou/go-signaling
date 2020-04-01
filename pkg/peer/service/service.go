package service

import "p2p/pkg/peer"

type service struct {
	peerMgr *peer.Manager
}

// NewService returns a user.Servicer implementation.
func NewService(peerMgr *peer.Manager) *service {
	return &service{
		peerMgr: peerMgr,
	}
}

func (s *service) SetPub() {

}

func (s *service) GetPub() {

}

func (s *service) SendAnswer() {

}

func (s *service) SendOffer() {

}
