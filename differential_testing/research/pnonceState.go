package research

import (
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	leveldb_opt "github.com/syndtr/goleveldb/leveldb/opt"
)

// modi-substate: pthread which is used to implement multithreaded replay

var nonceSubstateDB *leveldb.DB

func OpenNonceSubstateDB() {
	fmt.Println("modi-substate: openNonceSubstateDB")

	var nonceDir = filepath.Join("nonce-substate")
	var err error
	var opt leveldb_opt.Options

	opt.BlockCacheCapacity = 1 * leveldb_opt.GiB
	opt.OpenFilesCacheCapacity = 50

	nonceSubstateDB, err = leveldb.OpenFile(nonceDir, &opt)
	if _, corrupted := err.(*leveldb_errors.ErrCorrupted); corrupted {
		nonceSubstateDB, err = leveldb.RecoverFile(nonceDir, &opt)
	}
	if err != nil {
		panic(fmt.Errorf("error opening substate leveldb %s: %v", nonceDir, err))
	}

	fmt.Printf("modi-substate: successfully opened %s leveldb\n", nonceDir)
}

func CloseNonceSubstateDB() {
	defer fmt.Println("modi-substate: CloseNoncesubstateDB")

	nonceSubstateDB.Close()
	fmt.Println("modi-substate: successfully closed nonce leveldb")
}

func GetPnonce(addr common.Address) uint64 {
	var nonce uint64
	value, err := nonceSubstateDB.Get(addr.Bytes(), nil)
	if err != nil {
		// panic(fmt.Errorf("modi-substate: error getting nonce %s: %v", addr.String(), err))
		// If the value dose not exist in the database, it is considerd 0
		fmt.Println("don't get value")
		return 0
	}
	err = rlp.DecodeBytes(value, &nonce)
	if err != nil {
		fmt.Println(fmt.Errorf("error happens in decode"))
	}
	return nonce
}

func PutPnonce(addr common.Address, pnonce uint64) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Errorf("modi-substate: error putting substate %s into substate DB", addr))
		}
	}()
	value, err := rlp.EncodeToBytes(pnonce)
	if err != nil {
		return
	}
	err = nonceSubstateDB.Put(addr.Bytes(), value, nil)
	if err != nil {
		return
	}
}
