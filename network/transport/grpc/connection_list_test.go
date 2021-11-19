/*
 * Copyright (C) 2021 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package grpc

import (
	"github.com/nuts-foundation/nuts-node/network/transport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_connectionList_closeAll(t *testing.T) {
	cn := connectionList{}
	connA, _ := cn.getOrRegister(transport.Peer{ID: "a"}, nil)
	closerA := connA.closer()
	connB, _ := cn.getOrRegister(transport.Peer{ID: "b"}, nil)
	closerB := connB.closer()
	cn.closeAll()

	assert.Len(t, closerA, 1)
	assert.Len(t, closerB, 1)
}

func Test_connectionList_getOrRegister(t *testing.T) {
	t.Run("second call with same peer ID should return same connection", func(t *testing.T) {
		cn := connectionList{}
		connA, created1 := cn.getOrRegister(transport.Peer{ID: "a"}, nil)
		assert.True(t, created1)
		connASecondCall, created2 := cn.getOrRegister(transport.Peer{ID: "a"}, nil)
		assert.False(t, created2)
		assert.Equal(t, connA, connASecondCall)
	})
	t.Run("call with other peer ID should return same connection", func(t *testing.T) {
		cn := connectionList{}
		connA, created1 := cn.getOrRegister(transport.Peer{ID: "a"}, nil)
		assert.True(t, created1)
		connB, created2 := cn.getOrRegister(transport.Peer{ID: "b"}, nil)
		assert.True(t, created2)
		assert.NotEqual(t, connA, connB)
	})
}

func Test_connectionList_remove(t *testing.T) {
	cn := connectionList{}
	connA, _ := cn.getOrRegister(transport.Peer{ID: "a"}, nil)
	connB, _ := cn.getOrRegister(transport.Peer{ID: "b"}, nil)
	connC, _ := cn.getOrRegister(transport.Peer{ID: "c"}, nil)

	assert.Len(t, cn.list, 3)
	cn.remove(connB)
	assert.Len(t, cn.list, 2)
	assert.Contains(t, cn.list, connA)
	assert.Contains(t, cn.list, connC)
}