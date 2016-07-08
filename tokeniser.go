package ini

import (
	"errors"
	"io"

	"github.com/MJKWoolnough/parser"
)

const (
	sectionOpen  = '['
	sectionClose = ']'
	commentA     = '#'
	commentB     = ';'
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
	switch t.ExceptRun(string(sectionClose) + "\n") {
	case sectionClose:
	case -1:
		d.Err = io.ErrUnexpectedEOF
		return t.Error()
	default:
		d.Err = ErrInvalidName
		return t.Error()
	}
	data := t.Get()
	t.Accept(string(sectionClose))
	t.Get()
	return parser.Token{
		Type: tokenSection,
		Data: data,
	}, d.name
}

func (d *decoder) name(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun("\n")
	t.Get()
	switch t.Peek() {
	case -1:
		return t.Done()
	case sectionOpen:
		return d.sectionName(t)
	case commentA, commentB:
		return d.comment(t)
	}
	c := t.ExceptRun(string(d.NameValueDelim) + "\n")
	data := t.Get()
	switch c {
	case d.NameValueDelim:
		t.Accept(string(d.NameValueDelim))
		t.Get()
		return parser.Token{
			Type: tokenName,
			Data: data,
		}, d.value
	case -1:
		t.Err = io.ErrUnexpectedEOF
	default:
		t.Err = ErrInvalidName
	}
	return t.Error()
}

func (d *decoder) value(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	data := ""
	next := d.name
Loop:
	for {
		c := t.ExceptRun("\n\\")
		switch c {
		case -1:
			next = (*parser.Tokeniser).Done
			break Loop
		case '\n':
			break Loop
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
	}, next
}

func (d *decoder) comment(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.ExceptRun("\n")
	t.Get()
	return d.name(t)
}

func (d *decoder) section(p *parser.Parser) (parser.Phrase, parser.PhraseFunc) {
	if !p.Accept(tokenSection) {
		if p.Err == nil {
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
		if p.Err == nil {
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
	ErrInvalidName     = errors.New("invalid name")
)
