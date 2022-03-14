package storage

import (
	"github.com/nuts-foundation/go-leia/v2"
	"github.com/nuts-foundation/nuts-node/core"
	"go.etcd.io/bbolt"
	"io"
)

type Option func(config interface{})

type StorageEngine interface {
	core.Runnable
	core.Configurable
	core.Routable
	core.Named
}

type Warehouse interface {
	GetKVStore(engine core.Named, name string, opts ...Option) (*bbolt.DB, error)
	GetDocumentStore(engine core.Named, name string, opts ...Option) (leia.Store, error)
	io.Closer
}
