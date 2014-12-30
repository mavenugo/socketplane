package datastore

import (
	"errors"
	"net"

	log "github.com/socketplane/socketplane/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/socketplane/socketplane/Godeps/_workspace/src/github.com/socketplane/ecc"
	"github.com/socketplane/socketplane/ovs"
)

const (
	dataDir = "/tmp/socketplane"
)

type eccListener struct{}

func Init(bindInterface string, bootstrap bool) error {
	err := ecc.Start(bootstrap, bootstrap, bindInterface, dataDir)
	if err == nil {
		go ecc.RegisterForNodeUpdates(eccListener{})
	}
	return err
}

func Join(address string) error {
	return ecc.Join(address)
}

func Leave() error {
	return ecc.Leave()
}

func (e eccListener) NotifyNodeUpdate(nType ecc.NotifyUpdateType, nodeAddr string) {
	log.Debug("CLIENT UPDATE :", nType, nodeAddr)
	ip := net.ParseIP(nodeAddr).To4()
	if nType == ecc.NOTIFY_UPDATE_ADD {
		NewClusterNode(ip)
	} else {
		RemoveClusterNode(ip)
	}
}

func (e eccListener) NotifyKeyUpdate(nType ecc.NotifyUpdateType, key string, data []byte) {
}

func (e eccListener) NotifyStoreUpdate(nType ecc.NotifyUpdateType, store string, data map[string][]byte) {
}

func NewClusterNode(addr net.IP) error {
	log.Info("New Member Added : ", addr)
	Join(addr.String())
	// add the local node tunnels
	err := ovs.AddPeer(addr.String())
	if err != nil {
		return errors.New("Failed to adding new node")
	}
	log.Info("Added cluster member : ", addr)
	return nil
}

// Added to remove orphaned hosts
func RemoveClusterNode(addr net.IP) error {
	log.Info("Member Left : ", addr)
	// remove the local node tunnels
	err := ovs.DeletePeer(addr.String())
	if err != nil {
		return errors.New("Failed to adding new node")
	}
	log.Info("Deleted cluster member : ", addr)
	return nil
}

func parseAddrStr(ipStr net.IP) string {
	ipAddr, _ := net.ResolveIPAddr("ip", ipStr.String())
	return ipAddr.String()
}
