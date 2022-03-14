package storage

import (
	"errors"
	"fmt"
	"github.com/nuts-foundation/go-leia/v2"
	"github.com/nuts-foundation/nuts-node/core"
	"github.com/nuts-foundation/nuts-node/storage/log"
	"go.etcd.io/bbolt"
	"os"
	"path"
	"strings"
)

func NewWarehouse(datadir string) (Warehouse, error) {
	err := os.MkdirAll(datadir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &warehouse{
		datadir:        datadir,
		kvStores:       map[storeKey]*bbolt.DB{},
		documentStores: map[storeKey]leia.Store{},
	}, nil
}

type warehouse struct {
	datadir        string
	kvStores       map[storeKey]*bbolt.DB
	documentStores map[storeKey]leia.Store
}

func NoSyncOption(noSync bool) Option {
	return func(config interface{}) {
		switch c := config.(type) {
		case *leiaStoreConfig:
			c.noSync = noSync
		case *bbolt.Options:
			c.NoSync = noSync
		default:
			log.Logger().Errorf("Unsupported option NoSync for %T", config)
		}
	}
}

func applyOptions(config interface{}, opts []Option) {
	for _, opt := range opts {
		opt(config)
	}
}

func (s *warehouse) GetKVStore(engine core.Named, name string, opts ...Option) (*bbolt.DB, error) {
	key := getStorageKey(engine, name)

	// See if store is already there
	store := s.kvStores[key]
	if store != nil {
		return store, nil
	}

	// Derive config from options
	cfg := *bbolt.DefaultOptions // make sure to copy the defaults before making alterations
	applyOptions(&cfg, opts)

	filePath := path.Join(s.datadir, key.engine, fmt.Sprintf("%s.db", key.name))
	db, err := bbolt.Open(filePath, os.ModePerm, &cfg)
	if err == nil {
		s.kvStores[key] = db
	}
	return db, err
}

type leiaStoreConfig struct {
	noSync bool
}

func (s *warehouse) GetDocumentStore(engine core.Named, name string, opts ...Option) (leia.Store, error) {
	key := getStorageKey(engine, name)

	// See if store is already there
	store := s.documentStores[key]
	if store != nil {
		return store, nil
	}

	// Derive config from options
	cfg := &leiaStoreConfig{}
	applyOptions(cfg, opts)

	filePath := path.Join(s.datadir, key.engine, fmt.Sprintf("%s.db", key.name))
	db, err := leia.NewStore(filePath, cfg.noSync)
	if err == nil {
		s.documentStores[key] = db
	}
	return db, err
}

func getStorageKey(engine core.Named, name string) storeKey {
	return storeKey{engine: strings.ToLower(engine.Name()), name: strings.ToLower(name)}
}

func (s *warehouse) Close() error {
	success := true
	for key, db := range s.kvStores {
		log.Logger().Debugf("Closing KV store: %s", key)
		err := db.Close()
		if err != nil {
			log.Logger().Errorf("Error closing KV store %s: %v", key, err)
			success = false
		}
	}
	for key, db := range s.documentStores {
		log.Logger().Debugf("Closing document store: %s", key)
		err := db.Close()
		if err != nil {
			log.Logger().Errorf("Error closing document store %s: %v", key, err)
			success = false
		}
	}
	if !success {
		return errors.New("not all stores could be properly closed")
	}
	return nil
}
