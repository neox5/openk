package crypto

import "errors"

var (
    ErrNilParameter  = errors.New("parameter cannot be nil")
    ErrInvalidNonce = errors.New("nonce must be 12 bytes")
    ErrInvalidTag   = errors.New("tag must be 16 bytes")
)

const (
    NonceSize = 12 // 96 bits
    TagSize   = 16 // 128 bits
)

type Ciphertext struct {
    Nonce []byte // 96 bits (12 bytes)
    Data  []byte // The encrypted data itself
    Tag   []byte // 128 bits (16 bytes) authentication tag
}

func NewCiphertext(nonce, data, tag []byte) (*Ciphertext, error) {
    if nonce == nil || data == nil || tag == nil {
        return nil, ErrNilParameter
    }
    if len(nonce) != NonceSize {
        return nil, ErrInvalidNonce
    }
    if len(tag) != TagSize {
        return nil, ErrInvalidTag
    }
    return &Ciphertext{
        Nonce: nonce,
        Data:  data,
        Tag:   tag,
    }, nil
}