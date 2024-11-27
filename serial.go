package main

import (
	"encoding/hex"
	"strings"
	"sync/atomic"

	"go.bug.st/serial"
)

type serialPort struct {
	port        serial.Port
	displayMode int
	inputMode   int
	isRun       atomic.Bool
}

const (
	HEX = iota
	RAW
)

func openSerial(portName string, mode *serial.Mode) (*serialPort, error) {
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	out := &serialPort{
		port: port,
	}
	out.isRun.Store(true)
	return out, nil
}

func (s *serialPort) IsRun() bool {
	return s.isRun.Load()
}

func (s *serialPort) HandleRead(f func(data []byte)) {
	tmp := make([]byte, 512)
	for s.isRun.Load() {
		n, err := s.port.Read(tmp)
		if err != nil {
			return
		}
		switch s.displayMode {
		case HEX:
			dst := make([]byte, hex.EncodedLen(n))
			hex.Encode(dst, tmp[:n])
		case RAW:
			f(tmp[:n])
		}
	}
}

func (s *serialPort) WriteString(str string) error {
	if !s.isRun.Load() {
		return nil
	}
	var err error
	var data []byte
	switch s.inputMode {
	case HEX:
		str = strings.ReplaceAll(str, " ", "")
		data, err = hex.DecodeString(str)
	case RAW:
		data = []byte(str)
	}
	if err != nil {
		return err
	}
	_, err = s.port.Write(data)
	return err
}

func (s *serialPort) Close() error {
	s.isRun.Store(false)
	return s.port.Close()
}
