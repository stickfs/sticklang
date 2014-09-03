package lexer

import "fmt"

type token struct {
	typ tokenType
	val string
}

func (l token) Valid() bool {
	return l.val != ""
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

type tokenChan chan token

type lexWorker struct {
	currentLine int
	out         tokenChan
}

func newWorker() *lexWorker {
	return &lexWorker{
		out: make(chan token),
	}
}

func (l *lexWorker) run() {
	close(l.out)
}
