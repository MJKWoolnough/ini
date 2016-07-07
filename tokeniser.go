package ini

import (
	"errors"

	"github.com/MJKWoolnough/parser"
)

const (
	sectionOpen  = '['
	sectionClose = ']'
)

const (
	tokenSection parser.TokenType = iota
	tokenName
	tokenValue
)

const (
	phraseSection parser.PhraseType = iota
	phraseNameValue
)

func (d *decoder) sectionName(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept(string(sectionOpen))
	t.Get()
	t.ExceptRun(string(sectionClose))
	data := t.Get()
	t.Accept(string(sectionClose))
	t.Get()
	return parser.Token{
		Type: tokenSection,
		Data: data,
	}, d.name
}

func (d *decoder) name(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Peek() == sectionOpen {
		return d.sectionName(t)
	}
	t.ExceptRun(string(d.NameValueDelim))
	data := t.Get()
	t.Accept(string(d.NameValueDelim))
	t.Get()
	return parser.Token{
		Type: tokenName,
		Data: data,
	}, d.value
}

func (d *decoder) value(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	var data string
	for {
		switch t.ExceptRun("\n\\") {
		case '\n':
			break
		case '\\':
			data += t.Get()
			t.Accept("\\")
			if d.AllowMultiline && t.Peek() == '\n' {
				data += "\n"
			} else {
				data += "\\"
			}
		}
	}
	data += t.Get()
	return parser.Token{
		Type: tokenValue,
		Data: data,
	}, d.name
}

func (d *decoder) section(p *parser.Parser) (parser.Phrase, parser.PhraseFunc) {
	if !p.Accept(tokenSection) {
		if p.Err != nil {
			p.Err = ErrUnexpectedError
		}
		return p.Error()
	}
	return parser.Phrase{
		Type: phraseSection,
		Data: p.Get(),
	}, d.nameValue
}

func (d *decoder) nameValue(p *parser.Parser) (parser.Phrase, parser.PhraseFunc) {
	if !p.Accept(tokenName) {
		return d.section(p)
	}
	if !p.Accept(tokenValue) {
		if p.Err != nil {
			p.Err = ErrUnexpectedError
		}
		return p.Error()
	}

	return parser.Phrase{
		Type: phraseNameValue,
		Data: p.Get(),
	}, d.nameValue
}

//Errors
var (
	ErrUnexpectedError = errors.New("unexpected error parsing INI file")
)
