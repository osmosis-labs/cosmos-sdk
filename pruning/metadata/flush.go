package metadata

import (
	"encoding/binary"
	"errors"
	"fmt"

	dbm "github.com/tendermint/tm-db"
)

const(
	uint64Size = 8

	pruneHeightsKey         = "s/pruneheights"
)

type Flusher interface {
	Flush()

	FlushBatch(data interface{}, batch dbm.Batch)

	Load()
}

type PruneHeightFlusher struct {

}

func (*PruneHeightFlusher) FlushBatch(data interface{}, batch dbm.Batch) {
	pruneHeights, ok := data.([]int64)

	if !ok {
		panic(errors.New(""))
	}

	switch t := data.(type) { 
    default:
        panic(fmt.Sprintf("unexpected type %T", t))
    case []int64:
        
    }

	bz := make([]byte, 0, uint64Size * len(pruneHeights))
		for _, ph := range pruneHeights {
			buf := make([]byte, uint64Size)
			binary.BigEndian.PutUint64(buf, uint64(ph))
			bz = append(bz, buf...)
		}
	
		if err := batch.Set([]byte(pruneHeightsKey), bz); err != nil {
			panic(err)
		}
}