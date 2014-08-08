package sticklang

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// heavily influenced by
// http://cuddle.googlecode.com/hg/talk/lex.html

type token struct {
	typ tokenType
	val string
}

func (l token) String() string {
	switch l.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return l.val
	}
	if len(l.val) == 0 {
		return string(l.typ)
	}
	toPrint := fmt.Sprintf("%q:%q", string(l.typ), l.val)
	if len(toPrint) > 20 {
		toPrint = fmt.Sprintf("%.10q...:%s", l.val)
	}
	return toPrint
}

type tokenType string

const (
	tokenAmpAmp              tokenType = "&&"
	tokenBackSlash                     = "/"
	tokenBackTick                      = "`"
	tokenCloseBrace                    = "}"
	tokenCloseBracket                  = "]"
	tokenCloseParen                    = ")"
	tokenColon                         = ":"
	tokenColonEquals                   = ":="
	tokenComma                         = ","
	tokenCp                            = "cp"
	tokenDef                           = "def"
	tokenDo                            = "do"
	tokenImport                        = "import"
	tokenDot                           = "."
	tokenDoubleQuote                   = "\""
	tokenElse                          = "else"
	tokenEnd                           = "end"
	tokenEOF                           = "<EOF>"
	tokenError                         = "<ERR>"
	tokenForwardSlash                  = "\\"
	tokenHyphenArrow                   = "->"
	tokenIf                            = "if"
	tokenMv                            = "mv"
	tokenNewline                       = "\n"
	tokenNumber                        = "<NUMBER>"
	tokenOpenBrace                     = "{"
	tokenOpenBracket                   = "["
	tokenOpenComment                   = "#"
	tokenOpenParen                     = "("
	tokenPathSegment                   = "<PATH>"
	tokenPeriod                        = "."
	tokenPipe                          = "|"
	tokenPipePipe                      = "||"
	tokenProgramOption                 = "<OPTION>"
	tokenRedirectDoubleRight           = ">>"
	tokenRedirectLeft                  = "<"
	tokenRedirectRight                 = ">"
	tokenRm                            = "rm"
	tokenString                        = "<STRING>"
	tokenTypeKw                        = "type"
	tokenUnless                        = "unless"
	tokenVariableReference             = "<VAR>"
)

const (
	alphabetic   string = "abcdefhgijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumeric        = "abcdefhgijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	whitespace          = " \t\r\n"
)

type stateFn func(*lexer) stateFn

// type stateStack struct {
// 	top *stateFn
// 	size int
// }

type lexer struct {
	name   string     // error reporting
	input  string     // the string being scanned
	start  int        // start position
	pos    int        // current position of input
	width  int        // width of last rune read
	tokens chan token // channel of scanned tokens

}

const eof rune = -1

func (l *lexer) run() {
	for state := lexStick; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) next() (theRune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	theRune, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return theRune
}

func (l *lexer) acceptWord(raw tokenType) bool {
	added := 0 // if we don't match the keyword, need to reduce pos by value read

	for _, char := range raw {
		r := l.next()
		added += l.width
		if r == eof || r != char {
			l.pos -= added
			return false
		}
	}
	return true
}

func isAlphaNumeric(r rune) bool {
	return strings.IndexRune(alphanumeric, r) >= 0
}

func isSpace(r rune) bool {
	return strings.IndexRune(whitespace, r) >= 0
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.input)-1
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l
}

func lexStick(l *lexer) stateFn {
	for {
		if l.acceptWord(tokenBackSlash) {
			l.emit(tokenBackSlash)
		}
		if l.acceptWord(tokenImport) {
			l.emit(tokenImport)
		}

		switch r := l.next(); {
		case r == eof:
			l.emit(tokenEOF)
			return nil
		case r == '\n':
			l.emit(tokenNewline)
		case isSpace(r):
			l.ignore()
		case r == '|':
			l.emit(tokenPipe)
		case r == '"':
			l.ignore()
			return lexDoubleQuote(l)
		case r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case r == '&':
			{
				if r = l.next(); r == '&' {
					l.emit(tokenAmpAmp)
					return lexStick
				} else {
					return l.errorf(fmt.Sprintf("Unexpected %q following after &", r))
				}
			}

		}
	}
}

func lexDoubleQuote(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '"':
			l.backup()
			l.emit(tokenString)
			l.next()
			l.ignore()
			return lexStick
		case r == eof:
			return l.errorf("Unexpected eof during string")
		}

	}
}

func lexEscape(ll *lexer, ret stateFn) stateFn {
	return func(l *lexer) stateFn {
		if !l.accept("\"n") {
			return l.errorf("%v is not a valid escapable char", l.peek())
		} else {
			return ret
		}
	}
}

func lexNumber(l *lexer) stateFn {
	l.accept("-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenNumber)
	return lexStick
}
