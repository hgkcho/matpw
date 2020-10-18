package password

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	digits = "0123456789"
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower  = "abcdefghijklmnopqrstuvwxyz"
	// syms   = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	syms = "!#$%&()*+,-.<=>?@[]_{}"
	// CharAlpha is the class of letters
	CharAlpha = upper + lower
	// CharAlphaNum is the class of alpha-numeric characters
	CharAlphaNum = digits + upper + lower
	// CharAll is the class of all characters
	CharAll = digits + upper + lower + syms
)

// Password represents password
type Password struct {
	Title       string   `json:"title"`
	Account     string   `json:"account"`
	Descripiton string   `json:"description"`
	PasswordSet []string `json:"PasswordSet"`
}

// Create create password set
func (p *Password) Create() error {
	passLen := 5 * 5
	charLen := len(CharAll)
	var wg = new(sync.WaitGroup)
	var ch = make(chan string, passLen)

	for i := 0; i < passLen; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			char := generateTwoChar(charLen)
			ch <- string(char)
		}()
	}
	wg.Wait()
	close(ch)

	for v := range ch {
		p.PasswordSet = append(p.PasswordSet, v)
	}
	return nil

}

func generateTwoChar(max int) []byte {
	var ret []byte
	rand.Seed(time.Now().UnixNano())
	ret = append(ret, CharAll[rand.Intn(max)])
	ret = append(ret, CharAll[rand.Intn(max)])
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
			if k== 24 {
			  ret += fmt.Sprintln("   --------------------------")
			} else {
				ret += fmt.Sprintln("                                            ")
			}
		} else if k%5 == 0{
			ret += fmt.Sprintf("   | %s ", v)
		} else {
			ret += fmt.Sprintf("| %s ", v)
		}
	}
	return
}
