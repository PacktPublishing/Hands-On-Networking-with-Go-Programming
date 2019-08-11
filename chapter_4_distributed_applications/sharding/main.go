package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

func main() {
	s, _ := NewSharder(32)
	fmt.Println(s.Shard("id1"))
	fmt.Println(s.Shard("id2"))
	fmt.Println(s.Shard("id3"))
	fmt.Println(s.Shard("id4"))
	fmt.Println(s.Shard("id5"))
	fmt.Println(s.Shard("id6"))
	fmt.Println(s.Shard("id7"))
	fmt.Println(s.Shard("id8"))
	fmt.Println(s.Shard("id9"))
}

// Sharder creates shards.
type Sharder struct {
	divider byte
}

// ErrMustShardToMoreThanOne is returned when not enough shards are being made.
var ErrMustShardToMoreThanOne = errors.New("sharder: must shard to more than one shard")

// ErrMustShardToLessThan256 is returned when too many shards are being made.
var ErrMustShardToLessThan256 = errors.New("sharder: must shard to less than 256 shards")

// NewSharder allows sharding input data based on its SHA256 hash.
func NewSharder(shards int) (s Sharder, err error) {
	if shards <= 1 {
		err = ErrMustShardToMoreThanOne
		return
	}
	if shards > 255 {
		err = ErrMustShardToLessThan256
		return
	}
	s = Sharder{
		divider: byte(255) / byte(shards-1),
	}
	return
}

// Shard selects the shard an ID should be within.
func (s Sharder) Shard(of string) (index int) {
	hash := sha256.Sum256([]byte(of))
	firstByte := hash[0]
	return s.shardIndex(firstByte)
}

func (s Sharder) shardIndex(b byte) int {
	return int(b / s.divider)
}
