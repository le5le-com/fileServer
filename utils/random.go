package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/dchest/captcha"
)

var numChars = []byte("0123456789")

// GetRandString 获取随机字符串
func GetRandString(strLen uint8) (string, error) {
	b := make([]byte, strLen)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", errors.New("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

// GetRandCode 获取一个数据code
func GetRandCode(strLen uint8) string {
	b := captcha.RandomDigits(int(strLen))
	for i, c := range b {
		b[i] = numChars[c]
	}
	return string(b)
}

func readmachineID() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

var objectIDCounter uint32
var machineID = readmachineID()

// GetGUID 获取GUID
func GetGUID() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)

	b[4] = machineID[0]
	b[5] = machineID[1]
	b[6] = machineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return hex.EncodeToString(b[:])
}
