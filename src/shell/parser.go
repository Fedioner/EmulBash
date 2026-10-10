package shell

import (
	"errors"
	"os"
	"strings"
	"unicode"
)

var (
	ErrUnclosedQuote   = errors.New("syntax error: unclosed quote")
	ErrBadSubstitution = errors.New("syntax error: bad substitution")
)

type parser struct {
	src    []rune
	pos    int
	word   strings.Builder
	inWord bool
	words  []string
}

func Parse(line string) ([]string, error) {
	p := &parser{src: []rune(line)}
	for p.pos < len(p.src) {
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	p.endWord()
	return p.words, nil
}

func (p *parser) next() error {
	c := p.src[p.pos]
	switch c {
	case ' ', '\t':
		p.endWord()
		p.pos++
	case '\'':
		return p.singleQuoted()
	case '"':
		return p.doubleQuoted()
	case '$':
		return p.variable()
	case '\\':
		p.escaped()
	default:
		p.add(string(c))
		p.pos++
	}
	return nil
}

func (p *parser) add(s string) {
	if s != "" {
		p.inWord = true
	}
	p.word.WriteString(s)
}

func (p *parser) endWord() {
	if p.inWord {
		p.words = append(p.words, p.word.String())
	}
	p.word.Reset()
	p.inWord = false
}

func (p *parser) escaped() {
	p.pos++
	if p.pos < len(p.src) {
		p.add(string(p.src[p.pos]))
		p.pos++
	}
}

func (p *parser) singleQuoted() error {
	end := p.find('\'', p.pos+1)
	if end < 0 {
		return ErrUnclosedQuote
	}
	p.inWord = true
	p.add(string(p.src[p.pos+1 : end]))
	p.pos = end + 1
	return nil
}

func (p *parser) doubleQuoted() error {
	p.inWord = true
	p.pos++
	for p.pos < len(p.src) {
		switch p.src[p.pos] {
		case '"':
			p.pos++
			return nil
		case '$':
			if err := p.variable(); err != nil {
				return err
			}
		case '\\':
			p.escaped()
		default:
			p.add(string(p.src[p.pos]))
			p.pos++
		}
	}
	return ErrUnclosedQuote
}

func (p *parser) variable() error {
	p.pos++
	if p.pos < len(p.src) && p.src[p.pos] == '{' {
		return p.bracedVariable()
	}
	start := p.pos
	for p.pos < len(p.src) && isNameChar(p.src[p.pos]) {
		p.pos++
	}
	if start == p.pos {
		p.add("$")
		return nil
	}
	p.add(os.Getenv(string(p.src[start:p.pos])))
	return nil
}

func (p *parser) bracedVariable() error {
	end := p.find('}', p.pos)
	if end < 0 {
		return ErrBadSubstitution
	}
	name := string(p.src[p.pos+1 : end])
	if !validName(name) {
		return ErrBadSubstitution
	}
	p.add(os.Getenv(name))
	p.pos = end + 1
	return nil
}

func (p *parser) find(c rune, from int) int {
	for i := from; i < len(p.src); i++ {
		if p.src[i] == c {
			return i
		}
	}
	return -1
}

func isNameChar(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)
}

func validName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !isNameChar(c) {
			return false
		}
	}
	return true
}
