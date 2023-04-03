package research

import (
	"bytes"
	"strconv"
)

// turn the code which is string into []byte
func TurnCodeIntoBytes(code string) []byte {
	var result []byte
	for i := 0; i < len(code); i += 2 {
		stringUnit := code[i : i+2]
		byteUnit, err := strconv.ParseUint(stringUnit, 16, 64)
		PanicOnError(err)
		result = append(result, byte(byteUnit))
	}
	return result
}

func TurnBytesIntoCode(code []byte) string {
	var buffer bytes.Buffer
	var byteUint uint64
	for i := 0; i < len(code); i += 1 {
		byteUint = uint64(code[i])
		// fmt.Println(byteUint)
		stringUint := strconv.FormatUint(byteUint, 16)
		if byteUint < 16 {
			buffer.WriteString("0")
		}
		// fmt.Println(stringUint)
		buffer.WriteString(stringUint)
	}
	// fmt.Println(buffer.String())
	result := buffer.String()
	return result
}
