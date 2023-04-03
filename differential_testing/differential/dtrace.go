package differential

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type DTraceResult struct {
	TxBlockNum      uint64
	TxPositionNum   int
	Result          bool
	PreImplAddress  common.Address
	PostImplAdrress common.Address
	Storage         map[common.Hash]common.Hash
	ProxyStorage    map[common.Hash]common.Hash
}

func NewDTraceResult(block uint64, tx int, result bool) *DTraceResult {
	return &DTraceResult{
		TxBlockNum:      block,
		TxPositionNum:   tx,
		Result:          result,
		PreImplAddress:  OriVersion,
		PostImplAdrress: ModiVersion,
		Storage:         make(map[common.Hash]common.Hash),
		ProxyStorage:    make(map[common.Hash]common.Hash),
		// Storage: nil,
	}
}

type DTraceResultJSON struct {
	TxBlockNum      math.HexOrDecimal64         `json:"txblocknum,omitempty"`
	TxPositionNum   int                         `json:"txpositionnum,omitempty"`
	Result          bool                        `json:"result,omitempty"`
	PreImplAddress  common.Address              `json:"preimpladdress,omitempty"`
	PostImplAdrress common.Address              `json:"postimpladdress,omitempty"`
	Storage         map[common.Hash]common.Hash `json:"storage"`
	ProxyStorage    map[common.Hash]common.Hash `json:"proxystorage"`
}

func NewDTraceResultJSON(dt *DTraceResult) *DTraceResultJSON {
	return &DTraceResultJSON{
		TxBlockNum:      math.HexOrDecimal64(dt.TxBlockNum),
		TxPositionNum:   dt.TxPositionNum,
		Result:          dt.Result,
		PreImplAddress:  dt.PreImplAddress,
		PostImplAdrress: dt.PostImplAdrress,
		Storage:         dt.Storage,
		ProxyStorage:    dt.ProxyStorage,
	}
}

func (dt DTraceResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewDTraceResultJSON(&dt))
}
