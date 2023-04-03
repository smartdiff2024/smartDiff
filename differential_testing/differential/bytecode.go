package differential

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
)

type BytecodeAccount struct {
	BlockNum     uint64
	PositionNum  int
	CreateBin    []byte
	RuntimeBin   []byte
	Storage      map[common.Hash]common.Hash
	ProxyStorage map[common.Hash]common.Hash
}

type BytecodeAccounts map[common.Address]*BytecodeAccount

type ProxyAccount map[common.Address]*BytecodeAccounts

type BytecodeAccountJSON struct {
	BlockNum     math.HexOrDecimal64         `json:"blocknum,omitempty"`
	PositionNum  int                         `json:"positionnum,omitempty"`
	CreateBin    hexutil.Bytes               `json:"createbin,omitempty"`
	RuntimeBin   hexutil.Bytes               `json:"runtimebin,omitempty"`
	Storage      map[common.Hash]common.Hash `json:"storage"`
	ProxyStorage map[common.Hash]common.Hash `json:"proxyjson"`
}

func NewBytecodeAccountJSON(bc *BytecodeAccount) *BytecodeAccountJSON {
	return &BytecodeAccountJSON{
		BlockNum:     math.HexOrDecimal64(bc.BlockNum),
		PositionNum:  bc.PositionNum,
		CreateBin:    bc.CreateBin,
		RuntimeBin:   bc.RuntimeBin,
		Storage:      bc.Storage,
		ProxyStorage: bc.ProxyStorage,
	}
}

func (bc *BytecodeAccount) SetJson(bcJSON *BytecodeAccountJSON) {
	bc.BlockNum = uint64(bcJSON.BlockNum)
	bc.PositionNum = bcJSON.PositionNum
	bc.RuntimeBin = bcJSON.RuntimeBin
	bc.CreateBin = bcJSON.CreateBin
	bc.Storage = make(map[common.Hash]common.Hash)
	if bcJSON.Storage != nil {
		bc.Storage = bcJSON.Storage
	}
	bc.ProxyStorage = make(map[common.Hash]common.Hash)
	if bcJSON.ProxyStorage != nil {
		bc.ProxyStorage = bcJSON.Storage
	}
}

func (bc BytecodeAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewBytecodeAccountJSON(&bc))
}

func (bc *BytecodeAccount) UnmarshalJSON(b []byte) error {
	var err error
	var bcJSON BytecodeAccountJSON

	err = json.Unmarshal(b, &bcJSON)
	if err != nil {
		return err
	}

	bc.SetJson(&bcJSON)

	return nil
}

type BytecodeAccountsJSON map[common.Address]*BytecodeAccountJSON

func NewBytecodeAccountsJSON(ba BytecodeAccounts) BytecodeAccountsJSON {
	baJSON := make(BytecodeAccountsJSON)
	for addr, account := range ba {
		baJSON[addr] = NewBytecodeAccountJSON(account)
	}
	return baJSON
}

func (bc *BytecodeAccounts) SetJSON(bcJSON BytecodeAccountsJSON) {
	*bc = make(BytecodeAccounts)
	for addr, bJSON := range bcJSON {
		var b BytecodeAccount

		b.BlockNum = uint64(bJSON.BlockNum)
		b.PositionNum = bJSON.PositionNum
		b.CreateBin = bJSON.CreateBin
		b.RuntimeBin = bJSON.RuntimeBin

		b.Storage = make(map[common.Hash]common.Hash)
		if bJSON.Storage != nil {
			b.Storage = bJSON.Storage
		}

		b.ProxyStorage = make(map[common.Hash]common.Hash)
		if bJSON.ProxyStorage != nil {
			b.ProxyStorage = bJSON.ProxyStorage
		}

		(*bc)[addr] = &b

	}
}

func (bc BytecodeAccounts) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewBytecodeAccountsJSON(bc))
}

func (bc *BytecodeAccounts) UnmarshalJSON(b []byte) error {
	var err error
	var bcJSON BytecodeAccountsJSON

	err = json.Unmarshal(b, &bcJSON)
	if err != nil {
		return err
	}

	bc.SetJSON(bcJSON)

	return nil
}

type ProxyAccountJSON map[common.Address]BytecodeAccountsJSON

func NewProxyAccountsJSON(proxy ProxyAccount) ProxyAccountJSON {
	proxyJSON := make(ProxyAccountJSON)
	for addr, account := range proxy {
		proxyJSON[addr] = NewBytecodeAccountsJSON(*account)
	}
	return proxyJSON
}

func (pa *ProxyAccount) SetJSON(paJSON ProxyAccountJSON) {
	*pa = make(ProxyAccount)
	for addr, pJSON := range paJSON {
		bas := make(BytecodeAccounts)
		for implAddr, bJSON := range pJSON {
			var ba BytecodeAccount
			ba.BlockNum = uint64(bJSON.BlockNum)
			ba.PositionNum = bJSON.PositionNum
			ba.CreateBin = bJSON.CreateBin
			ba.RuntimeBin = bJSON.RuntimeBin
			bas[implAddr] = &ba
		}
		(*pa)[addr] = &bas
	}
}

func (pa ProxyAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewProxyAccountsJSON(pa))
}

func (pa *ProxyAccount) UnmarshalJSON(b []byte) error {
	var err error
	var paJSON ProxyAccountJSON

	err = json.Unmarshal(b, &paJSON)
	if err != nil {
		return err
	}

	pa.SetJSON(paJSON)

	return nil
}

var (
	ProxyByte ProxyAccount
)

type ImplPosition struct {
	ProxyAddress []common.Address
	ImplAddress  common.Address
}
type BlockProxy map[int]*ImplPosition

var BlockImpl map[uint64]BlockProxy

func init() {
	ProxyByte = make(ProxyAccount)
}

func ReadBytecodeCSV() {
	BlockImpl = make(map[uint64]BlockProxy)
	if CSVfileToStoreResult == "" {
		panic("Empty CSV FILE!")
	}
	csvFile, _ := os.Open(CSVfileToStoreResult)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	reader.FieldsPerRecord = -1

	for {
		implBytes := make(BytecodeAccounts)
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		lineLength := len(line)
		proxyAddress := line[5]

		for linenum := 8; linenum < lineLength; linenum += 5 {
			blockNumTemp, _ := strconv.Atoi(line[linenum+3])
			positionTemp, _ := strconv.Atoi(line[linenum+4])
			implBytes[common.HexToAddress(line[linenum])] = &BytecodeAccount{BlockNum: uint64(blockNumTemp), PositionNum: positionTemp}

			var value BlockProxy
			value, knownblock := BlockImpl[uint64(blockNumTemp)]
			if knownblock {
				value2, knowntx := value[positionTemp]
				if knowntx {
					var hasknowproxyaddress bool
					for _, knownproxyaddress := range value2.ProxyAddress {
						if knownproxyaddress == common.HexToAddress(proxyAddress) {
							hasknowproxyaddress = true
						}
					}
					if !hasknowproxyaddress {
						value2.ProxyAddress = append(value2.ProxyAddress, common.HexToAddress(proxyAddress))
					}
				} else {
					value[positionTemp] = &ImplPosition{ImplAddress: common.HexToAddress(line[linenum]), ProxyAddress: []common.Address{common.HexToAddress(proxyAddress)}}
				}
			} else {
				value = make(BlockProxy)
				value[positionTemp] = &ImplPosition{ImplAddress: common.HexToAddress(line[linenum]), ProxyAddress: []common.Address{common.HexToAddress(proxyAddress)}}
			}
			BlockImpl[uint64(blockNumTemp)] = value
		}
		ProxyByte[common.HexToAddress(proxyAddress)] = &implBytes
	}
}

func WriteBytecodeJSON() {
	if DirStoreJSONS == "" {
		panic("Empty Directory")
	}
	os.MkdirAll(DirStoreJSONS, 0755)
	for proxy, impls := range ProxyByte {
		proxyfilename := DirStoreJSONS + "/" + proxy.String() + ".json"
		data, _ := json.Marshal(impls)
		err := ioutil.WriteFile(proxyfilename, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func UpdateBytecodeJSON() {
	filepath := DirStoreJSONS + "/" + common.HexToAddress(DTraceFile).String() + ".json"
	file, _ := os.OpenFile(filepath, os.O_RDWR, 0666)
	file.Truncate(0)
	file.Seek(0, 0)
	defer file.Close()
	oriStorage := (*ProxyByte[ProxyAddress])[OriVersion].Storage
	modiStorage := (*ProxyByte[ProxyAddress])[ModiVersion]
	modiStorage.Storage = make(map[common.Hash]common.Hash)
	for storagekey, storagevalue := range oriStorage {
		modiStorage.Storage[storagekey] = storagevalue
	}
	oriProxyStorage := (*ProxyByte[ProxyAddress])[OriVersion].ProxyStorage
	modiProxyStorage := (*ProxyByte[ProxyAddress])[ModiVersion]
	modiProxyStorage.ProxyStorage = make(map[common.Hash]common.Hash)
	for storagekey, storagevalue := range oriProxyStorage {
		modiProxyStorage.ProxyStorage[storagekey] = storagevalue
	}
	(*ProxyByte[ProxyAddress])[OriVersion].Storage = make(map[common.Hash]common.Hash)
	(*ProxyByte[ProxyAddress])[OriVersion].ProxyStorage = make(map[common.Hash]common.Hash)
	updateWriter := bufio.NewWriter(file)
	for _, impls := range ProxyByte {
		data, _ := json.Marshal(impls)
		_, err := updateWriter.Write(data)
		if err != nil {
			panic(err)
		}
	}
	updateWriter.Flush()
}

func DirReadBytecodeJSONS() {
	files, _ := ioutil.ReadDir(DirStoreJSONS)
	for _, filename := range files {
		file := DirStoreJSONS + "/" + filename.Name()
		dataEncoded, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		implaccounts := make(BytecodeAccounts)
		json.Unmarshal(dataEncoded, &implaccounts)
		proxyAddress := strings.Split(filename.Name(), ".")[0]
		ProxyByte[common.HexToAddress(proxyAddress)] = &implaccounts
	}
}

func FileReadBytecodeJSONS() {
	files := strings.Split(DTraceFile, ",")
	for _, r_filename := range files {
		filename := common.HexToAddress(r_filename).String()
		file := DirStoreJSONS + "/" + filename + ".json"
		dataEncoded, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		implaccounts := make(BytecodeAccounts)
		json.Unmarshal(dataEncoded, &implaccounts)
		proxyAddress := strings.Split(filename, ".")[0]
		ProxyByte[common.HexToAddress(proxyAddress)] = &implaccounts
	}
	if (*ProxyByte[ProxyAddress])[OriVersion].ProxyStorage != nil && len((*ProxyByte[ProxyAddress])[OriVersion].ProxyStorage) != 0 {
		oriProxyStorage := (*ProxyByte[ProxyAddress])[OriVersion]
		modiProxyStorage := (*ProxyByte[ProxyAddress])[ModiVersion]
		for storagekey, storagevalue := range oriProxyStorage.ProxyStorage {
			modiProxyStorage.ProxyStorage[storagekey] = storagevalue
		}
	}
}

func PrintProxyByte() {
	for proxy, impls := range ProxyByte {
		bytecodes := *impls
		for impl, stage := range bytecodes {
			fmt.Println("init::proxyAddress: ", proxy, "impl: ", impl)
			for key, value := range stage.Storage {
				fmt.Println("storage::key: ", key, " value: ", value)
			}
			for key, value := range stage.ProxyStorage {
				fmt.Println("ProxyStorage::key: ", key, " value: ", value)
			}
		}
	}
}
