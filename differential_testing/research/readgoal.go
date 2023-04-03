package research

import "github.com/ethereum/go-ethereum/differential"

func NeedImpl(alloc SubstateAlloc) bool {
	for addr := range alloc {
		_, ok := differential.ProxyInits[addr]
		if ok {
			return true
		}
	}
	return false
}

func NeedTx(alloc SubstateAlloc) bool {
	for addr := range alloc {
		if addr == differential.OriVersion || addr == differential.ProxyAddress {
			return true
		}
	}
	return false
}

func ProxyIncludeCheck(alloc SubstateAlloc) {
	for addr := range alloc {
		if addr == differential.ProxyAddress {
			differential.ProxyInclude = true
			return
		}
	}
	differential.ProxyInclude = false
}

func ImplIncludeCheck(alloc SubstateAlloc) {
	for addr := range alloc {
		if addr == differential.OriVersion {
			differential.ImplInclude = true
			return
		}
	}
	differential.ImplInclude = false
}
