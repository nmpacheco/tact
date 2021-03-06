package tact

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/brunotm/kvs"
	"github.com/brunotm/kvs/goleveldb"
)

var (
	// workingPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	// Registry default collector registry
	Registry *registry
	// Store default persistence store
	Store kvs.Store
)

// init the core
func init() {
	var err error
	Registry = &registry{
		mtx:        &sync.RWMutex{},
		collectors: map[string]*Collector{},
		groups:     map[string][]*Collector{},
	}
	Store, err = goleveldb.New("./statedb") // TODO: make path configurable
	if err != nil {
		panic(err)
	}
}

// Close shutdown and stops the core
func Close() {
	if err := Store.Close(); err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			},
		).Info("Error closing KVStore")
	}
}
