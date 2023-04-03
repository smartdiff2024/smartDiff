package research

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/differential"
	"github.com/sirupsen/logrus"
)

func CreateOrCallCheck(substate *Substate, unsubstate *Substate) (*Substate, bool) {
	outalloc := substate.OutputAlloc
	var modiAddress common.Address
	needChange := 0
	for addr := range outalloc {
		if addr == differential.ProxyAddress {
			needChange++
			modiAddress = addr
		} else if addr == differential.OriVersion {
			needChange++
		} else {
			continue
		}
	}
	inalloc := substate.InputAlloc
	for addr := range inalloc {
		if addr == differential.OriVersion {
			call_replace(substate, unsubstate, addr)
			needChange--
		} else if addr == differential.ProxyAddress {
			proxy_replace(substate, unsubstate)
			needChange--
		} else {
			continue
		}
	}
	if needChange == 0 {
		return substate, true
	} else if needChange == 1 {
		if modiAddress == differential.ProxyAddress {
			return substate, true
		} else {
			return create_replace(substate)
		}
	} else {
		panic("error replace times")
	}
}

// update the storage of proxy address
func proxy_replace(substate *Substate, unsubstate *Substate) {
	logrus.Debug("invoke proxy change")
	goalcontent := differential.ProxyByte[differential.ProxyAddress]
	goalex := (*goalcontent)[differential.ModiVersion]
	inpalloc := (*substate).InputAlloc[differential.ProxyAddress]
	inpalloc.Storage = make(map[common.Hash]common.Hash)
	for key, value := range goalex.ProxyStorage {
		inpalloc.Storage[key] = value
	}
}

// update the createbin
func create_replace(substate *Substate) (*Substate, bool) {
	goalcontent := differential.ProxyByte[differential.ProxyAddress]
	goalex := (*goalcontent)[differential.ModiVersion]
	(*substate).Message.Data = goalex.CreateBin
	return substate, true
}

// update the runtimebin and storage slots
func call_replace(substate *Substate, unsubstate *Substate, addr common.Address) {
	goalcontent := differential.ProxyByte[differential.ProxyAddress]
	goalex := (*goalcontent)[differential.ModiVersion]
	inpalloc := (*substate).InputAlloc[addr]

	inpalloc.Code = goalex.RuntimeBin
	inpalloc.Storage = make(map[common.Hash]common.Hash)
	for key, value := range goalex.Storage {
		inpalloc.Storage[key] = value
	}
}

// update storage after the execution of the transaction
func UpdateStorage(evmAlloc *SubstateAlloc) {
	mapInit := false
	ModiProxyByte := (*differential.ProxyByte[differential.ProxyAddress])[differential.ModiVersion]
	for address, value := range *evmAlloc {
		if address == differential.OriVersion {
			for storagekey, storagevalue := range value.Storage {
				if !mapInit {
					if ModiProxyByte.Storage == nil {
						ModiProxyByte.Storage = make(map[common.Hash]common.Hash)
					}
					mapInit = true
				}
				ModiProxyByte.Storage[storagekey] = storagevalue
			}
		}
		mapInit = false
		if address == differential.ProxyAddress {
			for storagekey, storagevalue := range value.Storage {
				if !mapInit {
					if ModiProxyByte.ProxyStorage == nil {
						ModiProxyByte.ProxyStorage = make(map[common.Hash]common.Hash)
					}
					mapInit = true
				}
				ModiProxyByte.ProxyStorage[storagekey] = storagevalue
			}
		}
	}
}
