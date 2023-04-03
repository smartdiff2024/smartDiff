package differential

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

var ImplBlockNum uint64
var ImplPositionNum int

type ImplMap struct {
	ImplAddress common.Address
	BlockNum    uint64
	PositionNum int
}

type ProxyInitStruct struct {
	Initblocknum  uint64
	BlockPosition int
	ImplAddresses []ImplMap
	LastestImpl   common.Address
	KnownBlock    uint64
}

var ProxyInits map[common.Address]*ProxyInitStruct

func ReadProxyList() {
	ProxyInits = make(map[common.Address]*ProxyInitStruct)
	if FileToReadModi == "" {
		panic("Emtpy File")
	}
	csvFile, _ := os.Open(FileToReadModi)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	reader.FieldsPerRecord = -1
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[0] == string("FALSE") {
			initblocknum, _ := strconv.Atoi(line[6])
			blockPosition, _ := strconv.Atoi(line[7])
			knownblock, _ := strconv.Atoi(line[9])
			ProxyInits[common.HexToAddress(line[5])] = &ProxyInitStruct{
				Initblocknum:  uint64(initblocknum),
				BlockPosition: blockPosition,
				KnownBlock:    uint64(knownblock),
			}
		}
	}
}

func WriteProxyList() {
	if FileToReadModi == "" {
		panic("Emtpy File")
	}
	csvFile, _ := os.Open(FileToReadModi)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	reader.FieldsPerRecord = -1

	File, _ := os.OpenFile(CSVfileToStoreResult, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	WriterCsv := csv.NewWriter(bufio.NewWriter(File))
	defer File.Close()

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[0] == string("FALSE") {
			proxyaccount := ProxyInits[common.HexToAddress(line[5])]
			if len(proxyaccount.ImplAddresses) == 0 {
				WriterCsv.Write(line)
			} else {
				writerString := []string{}
				writerString = append(writerString, line[0:8]...)
				for implLength := 0; implLength < len(ProxyInits[common.HexToAddress(line[5])].ImplAddresses); implLength++ {
					writerString = append(writerString, proxyaccount.ImplAddresses[implLength].ImplAddress.String())
					writerString = append(writerString, strconv.FormatUint(proxyaccount.ImplAddresses[implLength].BlockNum, 10))
					writerString = append(writerString, strconv.Itoa(proxyaccount.ImplAddresses[implLength].PositionNum))
					writerString = append(writerString, []string{"-1", "-1"}...)
				}
				writerString = append(writerString, line[8:]...)
				WriterCsv.Write(writerString)
			}
		} else {
			WriterCsv.Write(line)
		}
	}
	WriterCsv.Flush()
}
