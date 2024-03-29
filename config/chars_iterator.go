package config

import (
	"bufio"
	"bytes"
	"io/ioutil"
)

type charIterator struct {
	scanner       *bufio.Scanner
	currentLine   int
	lastCatchChar string
}

func newCharIterator(filename string) (*charIterator, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return newCharIteratorWithBytes(bs), nil
}

func newCharIteratorWithBytes(bs []byte) *charIterator {
	s := bufio.NewScanner(bytes.NewBuffer(bs))
	s.Split(bufio.ScanRunes)
	return &charIterator{scanner: s, currentLine: 1}
}

func (self *charIterator) nextFilter(filter Filter) (word string, line int, has bool) {
	previous := ""
	for {
		if word, line, has = self.next(); !has {
			return
		} else if filter(word, previous) {
			return
		} else {
			previous = word
		}
	}
}

//不包括最后一个
func (self *charIterator) nextTo(filter Filter, includeLast bool) (word string, line int, has bool) {
	lastChar := ""
	for {
		if lastChar, line, has = self.next(); !has {
			return
		} else if filter(lastChar, word) {
			if includeLast {
				word += lastChar
			} else {
				self.lastCatchChar = lastChar
			}
			return
		} else {
			word += lastChar
		}
	}
}

func (it *charIterator) next() (char string, line int, has bool) {
	if it.lastCatchChar != "" {
		char = it.lastCatchChar
		line = it.currentLine
		has = true
		it.lastCatchChar = ""
		return
	}

	if has = it.scanner.Scan(); !has {
		return
	}
	char = it.scanner.Text()
	line = it.currentLine
	if char == "\n" {
		it.currentLine += 1
	}
	if char == "\\" {
		nextChar := ""
		nextChar, line, has = it.next()
		char = char + nextChar
	}
	return
}
