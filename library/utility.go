package utility

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"
	"io"
	"bufio"
	"bytes"
	"encoding/json"
)

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

// Utf8BomClean
// Clean returns b with the 3 byte BOM stripped off the front if it is present.
// If the BOM is not present, then b is returned.
func Utf8BomClean(b []byte) []byte {
	if len(b) >= 3 && b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		return b[3:]
	}
	return b
}

// Utf8BomNewReader
// NewReader returns an io.Reader that will skip over initial UTF-8 byte order marks.
func Utf8BomNewReader(r io.Reader) io.Reader {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		// not enough bytes
		return buf
	}
	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		_, _ = buf.Discard(3)
	}
	return buf
}

// ParseTokenInt
// token으로 문자열을 분리하여 integer slice로 반환하며 공백은 모두 무시한다
//  example)
//	s = "1234, 567, 1  "
//	token = ","
//	return = [1234] [567] [1]
func ParseTokenInt(s string, token string) []int {
	list := make([]int, 0)
	s = strings.Trim(s, " ")
	split := strings.Split(s, token)
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
		portNum, err := strconv.Atoi(split[i])
		if err == nil {
			list = append(list, portNum)
		}
	}
	return list
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Uint16FromBytes(bytes []byte) uint16 {
	bits := binary.LittleEndian.Uint16(bytes)
	return bits
}

func Uint32FromBytes(bytes []byte) uint32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return bits
}

func Uint64FromBytes(bytes []byte) uint64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return bits
}

// GetPrettyJsonStr
// byte array 를 보기 좋게 정렬된 json string 으로 반환
func GetPrettyJsonStr(message []byte) (string, error) {
	var logBuff bytes.Buffer
	err := json.Indent(&logBuff, message, "", "   ")
	if err != nil {
		return "", err
	}
	logBuff.Write([]byte("\n"))
	return logBuff.String(), nil
}
