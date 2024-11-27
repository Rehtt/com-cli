package main

import (
	"errors"
	"strconv"

	"go.bug.st/serial"
)

type settings interface {
	Name() string
	Key() string
	Get() string
	Value() any
	Set(s string) error
	Trigger() (data []string, showMsg bool, err error)
}

type patity string

func NewParity() *patity {
	p := patity(serial.NoParity)
	return &p
}

func (p *patity) Name() string {
	return "奇偶校验"
}

func (p *patity) Key() string {
	return "parity"
}

func (p *patity) Get() string {
	return string(*p)
}

func (p *patity) Value() any {
	switch *p {
	case "无校验":
		return serial.NoParity
	case "奇校验":
		return serial.OddParity
	case "偶校验":
		return serial.EvenParity
	case "1 校验位":
		return serial.MarkParity
	case "0 校验位":
		return serial.SpaceParity
	}
	return serial.NoParity
}

func (p *patity) Trigger() ([]string, bool, error) {
	return []string{"无校验", "奇校验", "偶校验", "1 校验位", "0 校验位"}, true, nil
}

func (p *patity) Set(str string) error {
	*p = patity(str)
	return nil
}

type dataBits int

func NewDataBits() *dataBits {
	d := dataBits(8)
	return &d
}

func (d *dataBits) Name() string {
	return "字符大小"
}

func (d *dataBits) Key() string {
	return "data_bits"
}

func (d *dataBits) Get() string {
	return strconv.Itoa(int(*d))
}

func (d *dataBits) Value() any {
	return int(*d)
}

func (d *dataBits) Trigger() ([]string, bool, error) {
	return []string{"5", "6", "7", "8"}, true, nil
}

func (d *dataBits) Set(str string) error {
	rate, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	*d = dataBits(rate)
	return nil
}

type stopBits string

func NewStopBits() *stopBits {
	s := stopBits("1")
	return &s
}

func (s *stopBits) Name() string {
	return "停止位"
}

func (s *stopBits) Key() string {
	return "stop_bits"
}

func (s *stopBits) Get() string {
	return string(*s)
}

func (s *stopBits) Value() any {
	switch *s {
	case "1":
		return serial.OneStopBit
	case "1.5":
		return serial.OnePointFiveStopBits
	case "2":
		return serial.TwoStopBits
	}
	return serial.OneStopBit
}

func (s *stopBits) Trigger() ([]string, bool, error) {
	return []string{"1", "1.5", "2"}, true, nil
}

func (s *stopBits) Set(str string) error {
	*s = stopBits(str)
	return nil
}

type baudRate int

func NewBaudRate() *baudRate {
	b := baudRate(9600)
	return &b
}

func (b *baudRate) Name() string {
	return "波特率"
}

func (b *baudRate) Key() string {
	return "baud_rate"
}

func (b *baudRate) Get() string {
	return strconv.Itoa(int(*b))
}

func (b *baudRate) Value() any {
	baudRate, _ := strconv.Atoi(string(*b))
	return baudRate
}

func (b *baudRate) Trigger() ([]string, bool, error) {
	return []string{"4800", "9600", "14400", "19200", "28800", "38400", "57600", "115200"}, true, nil
}

func (b *baudRate) Set(s string) error {
	rate, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*b = baudRate(rate)
	return nil
}

type port string

func NewPort() *port {
	var p port
	return &p
}

func (p *port) Name() string {
	return "端口"
}

func (p *port) Key() string {
	return "port"
}

func (p *port) Get() string {
	return string(*p)
}

func (p *port) Value() any {
	return string(*p)
}

func (p *port) Trigger() ([]string, bool, error) {
	list, err := serial.GetPortsList()
	if err != nil {
		return nil, true, err
	}
	if len(list) == 0 {
		return nil, true, errors.New("无可用端口")
	}
	return list, true, err
}

func (p *port) Set(s string) error {
	*p = port(s)
	return nil
}

type displayMode string

func NewDisplayMode() *displayMode {
	m := displayMode("RAW")
	return &m
}

func (m *displayMode) Name() string {
	return "显示模式"
}

func (m *displayMode) Key() string {
	return "display_mode"
}

func (m *displayMode) Get() string {
	return string(*m)
}

func (m *displayMode) Value() any {
	switch *m {
	case "HEX":
		return HEX
	case "RAW":
		return RAW
	}
	return RAW
}

func (m *displayMode) Trigger() ([]string, bool, error) {
	return []string{"HEX", "RAW"}, true, nil
}

func (m *displayMode) Set(s string) error {
	switch s {
	case "HEX":
		*m = displayMode(s)
	case "RAW":
		*m = displayMode(s)
	default:
		*m = displayMode("RAW")
	}
	return nil
}

type inputMode string

func NewInputMode() *inputMode {
	m := inputMode("RAW")
	return &m
}

func (m *inputMode) Name() string {
	return "输入模式"
}

func (m *inputMode) Key() string {
	return "input_mode"
}

func (m *inputMode) Get() string {
	return string(*m)
}

func (m *inputMode) Value() any {
	switch *m {
	case "HEX":
		return HEX
	case "RAW":
		return RAW
	}
	return RAW
}

func (m *inputMode) Trigger() ([]string, bool, error) {
	return []string{"HEX", "RAW"}, true, nil
}

func (m *inputMode) Set(s string) error {
	switch s {
	case "HEX":
		*m = inputMode(s)
	case "RAW":
		*m = inputMode(s)
	default:
		*m = inputMode("RAW")
	}
	return nil
}
