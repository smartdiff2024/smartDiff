package research

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	BEFORE uint = iota
	AFTER
)

type TxRecord struct {
	Addr                   string
	Blocknumber            uint64
	Txnum                  int
	Gasbefore              uint64
	Gasafter               uint64
	Modified_before_status uint64
	Modified_after_status  uint64
	Before_status          uint64
	After_status           uint64
}

var FileNameMap map[common.Address]*os.File

func InitFileMap() {
	FileNameMap = make(map[common.Address]*os.File)
}

func InitFile(addr common.Address) {
	headers := []string{
		"address", "blocknumber",
		"txnumber", "origin gas", "modified gas",
		"origin_before_status", "origin_after_status",
		"modi_before_status", "modi_after_status",
	}
	filename := fmt.Sprintf("%s%s", addr.String(), ".csv")
	file, openFileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	PanicOnError(openFileErr)
	stat, statErr := file.Stat()
	PanicOnError(statErr)
	if stat.Size() == 0 {
		_, err := file.WriteString(strings.Join(headers, ",") + "\n")
		PanicOnError(err)
	}
	FileNameMap[addr] = file
}

func CloseFile() {
	for _, file := range FileNameMap {
		closefileerr := file.Close()
		PanicOnError(closefileerr)
	}
}

func WriteTx(addr *[]common.Address, block uint64, tx int, usedgas uint64, before uint64, after uint64, records *[]TxRecord, beforeorafter uint) {
	if beforeorafter == AFTER {
		WriteModiTx(addr, block, tx, usedgas, before, after, records)
	}
	if beforeorafter == BEFORE {
		WriteOriginTx(addr, block, tx, usedgas, before, after, records)
	}
}

func WriteModiTx(addrs *[]common.Address, block uint64, tx int, usedgas uint64, before uint64, after uint64, records *[]TxRecord) {
	for i, addr := range *addrs {
		if i >= len(*records) {
			*records = append(*records, TxRecord{})
		}
		(*records)[i].Addr = addr.String()
		(*records)[i].Blocknumber = block
		(*records)[i].Txnum = tx
		(*records)[i].Gasafter = usedgas
		(*records)[i].Modified_before_status = before
		(*records)[i].Modified_after_status = after
	}
}

func WriteOriginTx(addrs *[]common.Address, block uint64, tx int, usedgas uint64, before uint64, after uint64, records *[]TxRecord) {
	// if addr.String() == record.Addr && block == record.Blocknumber && tx == record.Txnum {
	// 	record.Gasbefore = usedgas
	// 	record.Before_status = before
	// 	record.After_status = after
	// } else {
	// 	fmt.Fprintln(os.Stderr, "block_tx", block, tx, "block_tx2: ", record.Blocknumber, record.Txnum, "addr: ", addr, "addr before: ", record.Addr)
	// 	panic("inconsistent")
	// }
	for _, record := range *records {
		for _, addr := range *addrs {
			if addr.String() == record.Addr {
				record.Gasbefore = usedgas
				record.Before_status = before
				record.After_status = after
				WriteToCsv(addr, &record)
			}
		}
	}
}

func WriteToCsv(addr common.Address, record *TxRecord) {
	file := FileNameMap[addr]
	lockByName(file.Name())
	defer unlockByName(file.Name())

	//
	writer := csv.NewWriter(file)
	blocknumstring := strconv.FormatUint(record.Blocknumber, 10)
	txnumstr := strconv.Itoa(record.Txnum)
	gasbeforestr := strconv.FormatUint(record.Gasbefore, 10)
	gasafterstr := strconv.FormatUint(record.Gasafter, 10)
	mobsstr := strconv.FormatUint(record.Modified_before_status, 10)
	moasstr := strconv.FormatUint(record.Modified_after_status, 10)
	bsstr := strconv.FormatUint(record.Before_status, 10)
	asstr := strconv.FormatUint(record.After_status, 10)
	line := []string{record.Addr, blocknumstring, txnumstr, gasbeforestr, gasafterstr, mobsstr, moasstr, bsstr, asstr}
	err := writer.Write(line)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}
