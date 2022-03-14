package storage

import (
	"testing"
)

func NewTestWarehouse(t *testing.T, testDirectory string) Warehouse {
	w, err :=  NewWarehouse(testDirectory)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		_ = w.Close()
	})
	return w
}
