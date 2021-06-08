// Package ecdsa provides a crypto implementation for HotStuff using Go's 'crypto/ecdsa' package.
package ecdsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	hotstuff "github.com/HotstuffWASM/newNetwork"
)

// ErrHashMismatch is the error used when a partial certificate hash does not match the hash of a block.
var ErrHashMismatch = fmt.Errorf("certificate hash does not match block hash")

// ErrPartialDuplicate is the error used when two or more signatures were created by the same replica.
var ErrPartialDuplicate = fmt.Errorf("cannot add more than one signature per replica")

// PrivateKey is an ECDSA private key.
//
// This struct wraps the regular ecdsa.PrivateKey in order to implement the hotstuff.PrivateKey interface.
type PrivateKey struct {
	*ecdsa.PrivateKey
}

// PublicKey returns the public key associated with the private key
func (pk PrivateKey) PublicKey() hotstuff.PublicKey {
	return pk.Public()
}

var _ hotstuff.PrivateKey = (*PrivateKey)(nil)

// Signature is an ECDSA signature
type Signature struct {
	r, s   *big.Int
	signer hotstuff.ID
}

// NewSignature creates a new Signature struct from the given values.
func NewSignature(r, s *big.Int, signer hotstuff.ID) *Signature {
	return &Signature{r, s, signer}
}

// Signer returns the ID of the replica that generated the signature.
func (sig Signature) Signer() hotstuff.ID {
	return sig.signer
}

// R returns the r value of the signature
func (sig Signature) R() *big.Int {
	return sig.r
}

// S returns the s value of the signature
func (sig Signature) S() *big.Int {
	return sig.s
}

// ToString returns the signature as a string
func (sig Signature) ToString() string {
	sign := sig.r.String() + "-" + sig.s.String() + "-" + strconv.FormatUint(uint64(sig.signer), 10)
	return sign
}

// ToBytes returns a raw byte string representation of the signature
func (sig Signature) ToBytes() []byte {
	var b []byte
	b = append(b, sig.r.Bytes()...)
	b = append(b, sig.s.Bytes()...)
	return b
}

var _ hotstuff.Signature = (*Signature)(nil)

// PartialCert is an ECDSA signature and the hash that was signed.
type PartialCert struct {
	Signature *Signature
	hash      hotstuff.Hash
}

// NewPartialCert initializes a PartialCert struct from the given values.
func NewPartialCert(signature *Signature, hash hotstuff.Hash) *PartialCert {
	return &PartialCert{signature, hash}
}

// GetSignature returns the signature.
func (cert PartialCert) GetSignature() hotstuff.Signature {
	return cert.Signature
}

// GetStringSignature returns the string representation of the signature
func (cert PartialCert) GetStringSignature() string {
	return cert.Signature.ToString()
}

// BlockHash returns the hash of the block that was signed.
func (cert PartialCert) BlockHash() hotstuff.Hash {
	return cert.hash
}

// ToBytes returns a byte representation of the partial certificate.
func (cert PartialCert) ToBytes() []byte {
	return append(cert.hash[:], cert.Signature.ToBytes()...)
}

func (cert PartialCert) String() string {
	return fmt.Sprintf("PartialCert{ Block: %.6s, Signer: %d }", cert.hash.String(), cert.Signature.signer)
}

var _ hotstuff.PartialCert = (*PartialCert)(nil)

// QuorumCert is a set of signature that form a quorum certificate for a block.
type QuorumCert struct {
	signatures map[hotstuff.ID]*Signature
	hash       hotstuff.Hash
}

// NewQuorumCert initializes a new QuorumCert struct from the given values.
func NewQuorumCert(signatures map[hotstuff.ID]*Signature, hash hotstuff.Hash) *QuorumCert {
	return &QuorumCert{signatures, hash}
}

// GetSignatures returns the signatures within the quorum certificate.
func (qc QuorumCert) GetSignatures() map[hotstuff.ID]*Signature {
	return qc.signatures
}

// GetStringSignatures returns the map of signatures as a string
func (qc QuorumCert) GetStringSignatures() string {
	b := ""
	m := qc.GetSignatures()
	for key := range m {
		// fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
		b += strconv.FormatUint(uint64(key), 10) + "=" + m[key].ToString() + "\n"
	}
	return b
}

// BlockHash returns the hash of the block for which the certificate was created.
func (qc QuorumCert) BlockHash() hotstuff.Hash {
	return qc.hash
}

// ToBytes returns a byte representation of the quorum certificate.
func (qc QuorumCert) ToBytes() []byte {
	b := qc.hash[:]
	// sort signatures by id to ensure determinism
	sigs := make([]*Signature, 0, len(qc.signatures))
	for _, sig := range qc.signatures {
		i := sort.Search(len(sigs), func(i int) bool {
			return sig.signer < sigs[i].signer
		})
		sigs = append(sigs, nil)
		copy(sigs[i+1:], sigs[i:])
		sigs[i] = sig
	}
	for _, sig := range sigs {
		b = append(b, sig.ToBytes()...)
	}
	return b
}

func (qc QuorumCert) String() string {
	var sb strings.Builder
	for _, sig := range qc.signatures {
		sb.WriteString(" " + strconv.Itoa(int(sig.signer)) + " ")
	}
	return fmt.Sprintf("QC{ Block: %.6s, Sigs: [%s] }", qc.hash.String(), sb.String())
}

var _ hotstuff.QuorumCert = (*QuorumCert)(nil)

// TODO: consider adding caching back

type ecdsaCrypto struct {
	cfg hotstuff.Config
}

// New returns a new Signer and a new Verifier.
func New(cfg hotstuff.Config) (hotstuff.Signer, hotstuff.Verifier) {
	ec := &ecdsaCrypto{cfg}
	return ec, ec
}

func (ec *ecdsaCrypto) getPrivateKey() *PrivateKey {
	pk := ec.cfg.PrivateKey()
	return pk.(*PrivateKey)
}

// Sign signs a single block and returns a partial certificate.
func (ec *ecdsaCrypto) Sign(block *hotstuff.Block) (cert hotstuff.PartialCert, err error) {
	hash := block.Hash()
	r, s, err := ecdsa.Sign(rand.Reader, ec.getPrivateKey().PrivateKey, hash[:])
	if err != nil {
		fmt.Println("Error 21")
		return nil, err
	}
	return &PartialCert{
		&Signature{r, s, ec.cfg.ID()},
		hash,
	}, nil
}

// CreateQuorumCert creates a quorum certificate from a block and a set of signatures.
func (ec *ecdsaCrypto) CreateQuorumCert(block *hotstuff.Block, signatures []hotstuff.PartialCert) (cert hotstuff.QuorumCert, err error) {
	hash := block.Hash()
	qc := &QuorumCert{
		signatures: make(map[hotstuff.ID]*Signature),
		hash:       hash,
	}
	for _, s := range signatures {
		blockHash := s.BlockHash()
		if !bytes.Equal(hash[:], blockHash[:]) {
			fmt.Println("Error 22")
			return nil, ErrHashMismatch
		}
		if _, ok := qc.signatures[s.GetSignature().Signer()]; ok {
			fmt.Println("Error 23")
			return nil, ErrPartialDuplicate
		}
		qc.signatures[s.GetSignature().Signer()] = s.(*PartialCert).Signature
	}
	return qc, nil
}

// VerifyPartialCert verifies a single partial certificate.
func (ec *ecdsaCrypto) VerifyPartialCert(cert hotstuff.PartialCert) bool {
	// TODO: decide how to handle incompatible types. For now we'll simply panic
	sig := cert.GetSignature().(*Signature)
	replica, ok := ec.cfg.Replica(sig.Signer())
	if !ok {
		// logger.Info("ecdsaCrypto: got signature from replica whose ID (%d) was not in the config.")
		fmt.Println("Error 24")
		return false
	}
	pk := replica.PublicKey().(*ecdsa.PublicKey)
	hash := cert.BlockHash()
	return ecdsa.Verify(pk, hash[:], sig.R(), sig.S())
}

// VerifyQuorumCert verifies a quorum certificate.
func (ec *ecdsaCrypto) VerifyQuorumCert(cert hotstuff.QuorumCert) bool {
	// If QC was created for genesis, then skip verification.
	if cert.BlockHash() == hotstuff.GetGenesis().Hash() {
		return true
	}

	qc := cert.(*QuorumCert)
	if len(qc.GetSignatures()) < ec.cfg.QuorumSize() {
		fmt.Println("Error 25")
		return false
	}
	hash := qc.BlockHash()
	var wg sync.WaitGroup
	var numVerified uint64 = 0
	for id, pSig := range qc.GetSignatures() {
		info, ok := ec.cfg.Replica(id)
		if !ok {
			// logger.Error("VerifyQuorumSig: got signature from replica whose ID (%d) was not in config.", id)
		}
		pubKey := info.PublicKey().(*ecdsa.PublicKey)
		wg.Add(1)
		go func(pSig *Signature) {
			if ecdsa.Verify(pubKey, hash[:], pSig.R(), pSig.S()) {
				atomic.AddUint64(&numVerified, 1)
			}
			wg.Done()
		}(pSig)
	}
	wg.Wait()
	return numVerified >= uint64(ec.cfg.QuorumSize())
}

var _ hotstuff.Signer = (*ecdsaCrypto)(nil)
var _ hotstuff.Verifier = (*ecdsaCrypto)(nil)
