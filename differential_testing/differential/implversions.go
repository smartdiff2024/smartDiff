package differential

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var ModiVersions []common.Address
var OriVersion common.Address
var ProxyAddress common.Address
var ProxyInclude bool
var ImplInclude bool
var ModiVersion common.Address

func UpdateVersions(replaces string) {
	versions := strings.Split(replaces, ",")
	for _, version := range versions {
		ModiVersions = append(ModiVersions, common.HexToAddress(version))
	}
}

func HasImplVersions(addr common.Address) bool {
	for _, version := range ModiVersions {
		if version == addr {
			return true
		}
	}
	return false
}
