package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func PipeWR(reader1, reader2 io.Reader, writer1, writer2 io.Writer) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(writer1, reader2)
	}()

	go func() {
		defer wg.Done()
		io.Copy(writer2, reader1)
	}()
	wg.Wait()
	logrus.Info("pipe closed")
}

func Pipe(con1 *net.Conn, con2 *net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(*con1, *con2)
		if err != nil {
			logrus.Error(err)
		}
		if con1 != nil {
			(*con1).Close()
		}
		if con2 != nil {
			(*con2).Close()
		}

		logrus.Debug("io copy 1 closed")
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(*con2, *con1)
		if err != nil {
			logrus.Error(err)
		}
		if con1 != nil {
			(*con1).Close()
		}
		if con2 != nil {
			(*con2).Close()
		}
		logrus.Debug("io copy 2 closed")
	}()
	wg.Wait()
	logrus.Info("pipe closed")
}

func ToNetConn(wsconn *websocket.Conn) *net.Conn {
	return &[]net.Conn{
		wsconn,
	}[0]
}

func HashString(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func DebugOn() bool {
	str := os.Getenv("SSHX_DEBUG")
	if str == "" {
		return false
	}
	lowStr := strings.ToLower(str)
	if lowStr == "1" || lowStr == "true" || lowStr == "yes" {
		return true
	}
	return false
}
