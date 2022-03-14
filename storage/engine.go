package storage

import (
	"fmt"
	"github.com/nuts-foundation/go-leia/v2"
	"github.com/nuts-foundation/nuts-node/core"
	"go.etcd.io/bbolt"
)

func New() StorageEngine {
	return &engine{}
}

type engine struct {
	warehouse Warehouse
}

func (s *engine) GetKVStore(engine core.Named, name string, opts ...Option) (*bbolt.DB, error) {
	return s.GetKVStore(engine, name, opts...)
}

func (s *engine) GetDocumentStore(engine core.Named, name string, opts ...Option) (leia.Store, error) {
	return s.GetDocumentStore(engine, name, opts...)
}

func (s *engine) Close() error {
	return s.Shutdown()
}

func (s *engine) Name() string {
	return "Storage"
}

func (s *engine) Routes(router core.EchoRouter) {
	// TODO
}

func (s *engine) Configure(config core.ServerConfig) error {
	warehouse, err := NewWarehouse(config.Datadir)
	if err != nil {
		return fmt.Errorf("unable to configure storage: %w", err)
	}
	s.warehouse = warehouse
	return nil
}

func (s *engine) Start() error {
	return nil
}

func (s *engine) Shutdown() error {
	return s.warehouse.Close()
}

func (s warehouse) Name() string {
	return "storage"
}
