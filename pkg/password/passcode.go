package password

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

const (
	// keytext is string for aes cipher.
	// TODO: this should be stored in github secret. And I'm gonna embed it by `ldflogs` in build
	keytext = "asyqgie12798akljzmknm.ahkjkijl;k"
)

var (
	// DefaultPasscodePath is default file path which contains secret data
	DefaultPasscodePath = filepath.Join(os.Getenv("HOME"), ".matpw", "secret.json")
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

// Passcode represents passcode
type Passcode struct {
	Order     []int  `json:"-"`
	orderByte []byte `json:"-"`
	OrderCode []byte `json:"orderCode"`
	Path      string `json:"-"`

	rw     io.ReadWriter `json:"-"`
	cipher cipher.Block  `json:"-"`
}

// create 'aes' algorithm
func (p *Passcode) createCipher() error {
	c, err := aes.NewCipher([]byte(keytext))
	if err != nil {
		return err
	}
	p.cipher = c
	return nil
}

// Encript exec encript code
func (p *Passcode) Encript(w io.Writer) error {
	if p.Order == nil || reflect.DeepEqual(p.Order, []int{}) {
		return errors.New("p.Order is empty")
	}

	if err := p.createCipher(); err != nil {
		return err
	}

	orderByte := make([]byte, len(p.Order))
	for i, v := range p.Order {
		orderByte[i] = byte(v)
	}

	cfb := cipher.NewCFBEncrypter(p.cipher, commonIV)
	p.OrderCode = make([]byte, len(orderByte))
	cfb.XORKeyStream(p.OrderCode, orderByte)

	return json.NewEncoder(w).Encode(&p)
}

// Decript exec decript code
func (p *Passcode) Decript(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return err
	}
	if p.OrderCode == nil {
		return errors.New("p.OrderCode is not found")
	}

	if err := p.createCipher(); err != nil {
		return err
	}

	// 復号文字列
	cfbdec := cipher.NewCFBDecrypter(p.cipher, commonIV)
	orderByte := make([]byte, len(p.OrderCode))
	cfbdec.XORKeyStream(orderByte, p.OrderCode)

	order := make([]int, len(p.OrderCode))
	for i, v := range orderByte {
		order[i] = int(v)
	}
	p.Order = order
	return nil
}

func convertOrderIntToByte(order []int) []byte {
	buf := make([]byte, len(order))
	for i := range order {
		buf[i] = byte(order[i])
	}
	return buf
}
