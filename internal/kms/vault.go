package kms

import (
	"crypto/subtle"
	"errors"
	"runtime"
)

var (
	ErrNoKey      = errors.New("no key available")
	ErrKeyPresent = errors.New("vault already contains a key")
	ErrInvalidKey = errors.New("invalid key material")
)

// Vault provides secure key material storage
type Vault struct {
	key []byte
}

func NewVault() *Vault {
	return &Vault{}
}

// Store securely stores key material
func (v *Vault) Store(key []byte) error {
	if v.HasKey() {
		return ErrKeyPresent
	}

	if len(key) == 0 {
		return ErrInvalidKey
	}

	// Allocate and copy key material
	v.key = make([]byte, len(key))
	subtle.ConstantTimeCopy(1, v.key, key)

	return nil
}

// HasKey returns wether the vault contains a HasKey
func (v *Vault) HasKey() bool {
	return len(v.key) > 0
}

// UseKey performs an operation with stored key material
func (v *Vault) UseKey(op func(key []byte) error) error {
	if !v.HasKey() {
		return ErrNoKey
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	return op(v.key)
}

// Cleanup securely wipes stored key material
func (v *Vault) Cleanup() error {
	if v.HasKey() {
		subtle.ConstantTimeCopy(1, v.key, make([]byte, len(v.key)))
		v.key = nil
	}

	runtime.GC()
	return nil
}
