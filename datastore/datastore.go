package datastore

import (
	"os"

	log "github.com/socketplane/socketplane/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/socketplane/socketplane/Godeps/_workspace/src/github.com/socketplane/ecc"
)

const dataDir = "/tmp/socketplane"

func Init(bindInterface string, bootstrap bool, listener ecc.Listener) error {
	err := ecc.Start(bootstrap, bootstrap, bindInterface, dataDir)
	if err == nil {
		go ecc.RegisterForNodeUpdates(listener)
	}
	return err
}

func Join(address string) error {
	return ecc.Join(address)
}

func Leave() error {
	if err := ecc.Leave(); err != nil {
		log.Error(err)
		return err
	}
	if err := os.RemoveAll(dataDir); err != nil {
		log.Errorf("Error deleting data directory %s", err)
		return err
	}
	return nil
}
