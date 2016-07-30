package memcacheha

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

var MEMCACHEHA_HEADER []byte = []byte{0xfd, 0x37, 0xd3, 0x1b}
var ErrNotMemcacheHAKey = errors.New("not a memcacheha key")

type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string

	// Value is the Item's value.
	Value []byte

	// Flags are server-opaque flags whose semantics are entirely
	// up to the app.
	Flags uint32

	// Expiration is either nil (no expiry) or an absolute expiry time
	Expiration *time.Time
}

func NewItemFromMemcacheItem(item *memcache.Item) (*Item, error) {

	// Check basic header length
	if len(item.Value) < 8 {
		return nil, ErrNotMemcacheHAKey
	}

	// Check header
	for i, x := range MEMCACHEHA_HEADER {
		if item.Value[i] != x {
			return nil, ErrNotMemcacheHAKey
		}
	}

	// Read Expiration
	var mcExpiry uint32
	mcExpiry = mcExpiry | uint32(item.Value[4])<<24
	mcExpiry = mcExpiry | uint32(item.Value[5])<<16
	mcExpiry = mcExpiry | uint32(item.Value[6])<<8
	mcExpiry = mcExpiry | uint32(item.Value[7])

	var haExpiry *time.Time
	if mcExpiry != 0 {
		x := time.Unix(int64(mcExpiry), 0)
		haExpiry = &x
	}

	return &Item{
		Key:        item.Key,
		Value:      item.Value[8:],
		Flags:      item.Flags,
		Expiration: haExpiry,
	}, nil
}

func (item *Item) AsMemcacheItem() *memcache.Item {
	var mcExpiry int32
	var binTime []byte = make([]byte, 4)

	if item.Expiration != nil {
		// Write Expiration Int32
		mcExpiry = int32(item.Expiration.Unix())
		binTime[0] = byte((mcExpiry >> 24) & 0xFF)
		binTime[1] = byte((mcExpiry >> 16) & 0xFF)
		binTime[2] = byte((mcExpiry >> 8) & 0xFF)
		binTime[3] = byte(mcExpiry & 0xFF)

		// Change to relative for memcached
		mcExpiry = mcExpiry - int32(time.Now().Unix())

		// Catch negative offset (expire now)
		if mcExpiry < 1 {
			mcExpiry = 1
		}
	}

	var value []byte

	// Write Header
	value = append(value, MEMCACHEHA_HEADER...)

	// Write expiry time
	value = append(value, binTime...)

	// Write Data
	value = append(value, item.Value...)

	return &memcache.Item{
		Key:        item.Key,
		Value:      value,
		Flags:      item.Flags,
		Expiration: mcExpiry,
	}
}
