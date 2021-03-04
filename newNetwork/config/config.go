// Package config contains structs that are useful for initializing hotstuff.
//
// These structs do not implement the hotstuff.Replica or hotstuff.Config interfaces,
// but do contain more or less the same information.
package config

import (
	"crypto/ecdsa"

	hotstuff "github.com/HotstuffWASM/newNetwork"
)

// ReplicaInfo holds information about a replica.
type ReplicaInfo struct {
	ID      hotstuff.ID
	Address string
	PubKey  *ecdsa.PublicKey
}

// ReplicaConfig holds information needed by a replica.
type ReplicaConfig struct {
	ID         hotstuff.ID
	PrivateKey *ecdsa.PrivateKey
	Replicas   map[hotstuff.ID]*ReplicaInfo
}

// NewConfig returns a new ReplicaConfig instance.
func NewConfig(id hotstuff.ID, privateKey *ecdsa.PrivateKey) *ReplicaConfig {
	return &ReplicaConfig{
		ID:         id,
		PrivateKey: privateKey,
		Replicas:   make(map[hotstuff.ID]*ReplicaInfo),
	}
}
