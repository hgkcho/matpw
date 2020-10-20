package password

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// Digits = "0123456789"
	Digits = "0123456789"
	// Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Lower = "abcdefghijklmnopqrstuvwxyz"
	Lower = "abcdefghijklmnopqrstuvwxyz"
	// Symbols = "!#$%&()*+,-.<=>?@[]_{}"
	Symbols = "!#$%&()*+,-.<=>?@[]_{}"
	// CharAlpha is the class of letters
	CharAlpha = Upper + Lower
	// CharAlphaNum is the class of alpha-numeric characters
	CharAlphaNum = Digits + Upper + Lower
	// CharAll is the class of all characters
	CharAll = Digits + Upper + Lower + Symbols
)

// Password represents password
type Password struct {
	ID           uuid.UUID `json:"ID"`
	Service      string    `json:"service"`
	Account      string    `json:"account"`
	Descripiton  string    `json:"description"`
	PasswordSet  []string  `json:"passwordSet"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Path         string    `json:"-"`
	UseUppercase bool      `json:"-"`
	UseDigit     bool      `json:"-"`
	UseSymbol    bool      `json:"-"`
}

// Create create passwordSet
func (p *Password) Create() error {
	passLen := 5 * 5
	useChar := Lower
	if p.UseUppercase {
		useChar += Upper
	}
	if p.UseDigit {
		useChar += Digits
	}
	if p.UseSymbol {
		useChar += Symbols
	}
	var wg = new(sync.WaitGroup)
	var ch = make(chan string, passLen)

	for i := 0; i < passLen; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			char := generateTwoChar(useChar)
			ch <- string(char)
		}()
	}
	wg.Wait()
	close(ch)

	for v := range ch {
		p.PasswordSet = append(p.PasswordSet, v)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	p.ID = id
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func generateTwoChar(useChar string) []byte {
	var ret []byte
	rand.Seed(time.Now().UnixNano())
	ret = append(ret, useChar[rand.Intn(len(useChar))])
	ret = append(ret, useChar[rand.Intn(len(useChar))])
	return ret
}

// Render render
func (p *Password) Render() (ret string) {
	for k, v := range p.PasswordSet {
		if k == 0 {
			ret += fmt.Sprintln("   --------------------------")
		}
		if k%5 == 4 {
			ret += fmt.Sprintf("| %s |\n", v)
			if k == 24 {
				ret += fmt.Sprintln("   --------------------------")
			} else {
				ret += fmt.Sprintln("                                            ")
			}
		} else if k%5 == 0 {
			ret += fmt.Sprintf("   | %s ", v)
		} else {
			ret += fmt.Sprintf("| %s ", v)
		}
	}
	return
}
