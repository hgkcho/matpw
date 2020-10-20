package password

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	OrderByte []byte `json:"-"`
	OrderCode []byte `json:"orderCode"`
	Path      string `json:"-"`

	cipher cipher.Block `json:"-"`
}

func NewPasscode() *Passcode {
	return &Passcode{
		Path: DefaultPasscodePath,
	}
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
func (p *Passcode) Encript() error {
	// f, err := os.Create(p.Path)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	if err := p.createCipher(); err != nil {
		return err
	}

	p.convertOrderIntToByte()

	cfb := cipher.NewCFBEncrypter(p.cipher, commonIV)
	p.OrderCode = make([]byte, len(p.OrderByte))
	cfb.XORKeyStream(p.OrderCode, p.OrderByte)
	// json.NewEncoder(f).Encode(&p)

	// buf, err := json.MarshalIndent(p, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// _, err = f.Write(buf)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Decript exec encript code
func (p *Passcode) Decript() error {
	f, err := os.Open(p.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	json.NewDecoder(f).Decode(&p)

	if err = p.createCipher(); err != nil {
		return err
	}

	// 復号文字列
	cfbdec := cipher.NewCFBDecrypter(p.cipher, commonIV)
	p.OrderByte = make([]byte, 4)
	cfbdec.XORKeyStream(p.OrderByte, p.OrderCode)
	fmt.Printf("%x=>%v\n", p.OrderCode, p.OrderByte)
	return nil
}

func (p *Passcode) convertOrderIntToByte() {
	buf := make([]byte, len(p.Order))
	for i := range p.Order {
		buf[i] = byte(p.Order[i])
	}
	p.OrderByte = buf
}
