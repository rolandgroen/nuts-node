package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/nuts-foundation/nuts-node/crypto"
	v1 "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-node/vdr/doc"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestGenerateLoadsOfDIDs_MultipleNodes(t *testing.T) {
	const numberOfDIDs = 100
	const numberOfUpdatesPerDID = 10
	var nodeAddresses = []string{
		"http://localhost:1001",
		"http://localhost:1002",
		"http://localhost:1003",
		"http://localhost:1004",
		"http://localhost:1005",
		//"http://localhost:1006",
		//"http://localhost:1007",
		//"http://localhost:1008",
		//"http://localhost:1009",
	}

	// Create root TX
	v1.HTTPClient{
		ServerAddress: nodeAddresses[0],
		Timeout:       60 * time.Second,
	}.Create(v1.DIDCreateRequest{})
	const sleep = 5 * time.Second
	t.Logf("Waiting %s seconds for root TX to be distributed...", sleep)
	time.Sleep(sleep)

	wg := sync.WaitGroup{}
	wg.Add(len(nodeAddresses))

	for _, nodeAddress := range nodeAddresses {
		go func(nodeAddress string) {
			t.Logf("Creating data on %s", nodeAddress)
			client := v1.HTTPClient{
				ServerAddress: nodeAddress,
				Timeout:       60 * time.Second,
			}
		outer:
			for j := 0; j < numberOfDIDs; j++ {
				t.Logf("[ %s ] Creating DID Document...", nodeAddress)
				document, err := client.Create(v1.DIDCreateRequest{})
				if !assert.NoError(t, err) {
					break
				}
				t.Logf("[ %s ] DID Document created: %s", nodeAddress, document.ID)
				for x := 0; x < numberOfUpdatesPerDID; x++ {
					vm, err := doc.CreateNewVerificationMethodForDID(document.ID, testKeyCreator{})
					if !assert.NoError(t, err) {
						break outer
					}
					document, metadata, err := client.Get(document.ID.String())
					if !assert.NoError(t, err) {
						break outer
					}
					document.AddAuthenticationMethod(vm)
					_, err = client.Update(document.ID.String(), metadata.Hash.String(), *document)
					if !assert.NoError(t, err) {
						break outer
					}
					t.Logf("[ %s ] DID Document updated (%d/%d): %s", nodeAddress, x+1, numberOfUpdatesPerDID, document.ID)
				}
			}
			wg.Done()
		}(nodeAddress)
	}

	t.Logf("Waiting for %d routines to finish...", len(nodeAddresses))
	wg.Wait()
	t.Log("Done.")
}

func TestGenerateLoadsOfDIDs_SingleNode(t *testing.T) {
	const numberOfDIDs = 100
	const numberOfUpdatesPerDID = 10
	const numberOfThreads = 5
	const nodeAddress = "http://localhost:1323"

	client := v1.HTTPClient{
		ServerAddress: nodeAddress,
		Timeout:       60 * time.Second,
	}

	wg := sync.WaitGroup{}
	wg.Add(numberOfThreads)

	for i := 0; i < numberOfThreads; i++ {
		go func(goroutineNum int) {
		outer:
			for j := 0; j < numberOfDIDs/numberOfThreads; j++ {
				t.Logf("[ %d ] Creating DID Document...", goroutineNum)
				document, err := client.Create(v1.DIDCreateRequest{})
				if !assert.NoError(t, err) {
					break
				}
				t.Logf("[ %d ] DID Document created: %s", goroutineNum, document.ID)
				for x := 0; x < numberOfUpdatesPerDID; x++ {
					vm, err := doc.CreateNewVerificationMethodForDID(document.ID, testKeyCreator{})
					if !assert.NoError(t, err) {
						break outer
					}
					document, metadata, err := client.Get(document.ID.String())
					if !assert.NoError(t, err) {
						break outer
					}
					document.AddAuthenticationMethod(vm)
					_, err = client.Update(document.ID.String(), metadata.Hash.String(), *document)
					if !assert.NoError(t, err) {
						break outer
					}
					t.Logf("[ %d ] DID Document updated (%d/%d): %s", goroutineNum, x+1, numberOfUpdatesPerDID, document.ID)
				}
			}
			wg.Done()
		}(i)
	}

	t.Logf("Waiting for %d routines to finish...", numberOfThreads)
	wg.Wait()
	t.Log("Done.")
}

type testKeyCreator struct {
}

func (t testKeyCreator) New(namingFunc crypto.KIDNamingFunc) (crypto.Key, error) {
	pair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	kid, err := namingFunc(pair.Public())
	if err != nil {
		return nil, err
	}
	return crypto.TestKey{
		PrivateKey: pair,
		Kid:        kid,
	}, nil
}
