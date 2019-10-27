// Copyright 2019 zigma authors
// This file is part of the zigma library.
//
// The zigma library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The zigma library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the zigma library. If not, see <http://www.gnu.org/licenses/>.

package dht

import "time"

const (
	// the time after which a key/value pair expires; this is a time-to-live (TTL)
	// from the original publication date; this is normally 86400s
	tExpire = time.Hour * 24
	// the time after which an otherwise unaccessed bucket must be refreshed
	tRefresh = time.Hour
	// the interval between kademlia replication events, when a node is required to
	// publish its entire database
	tReplicate = time.Hour
	// the time after which the original publisher must republish a key/value pair
	tRepublish = time.Hour * 24
)

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
}
