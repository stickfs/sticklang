package sticklang

import (
	"strings"
	"testing"
)

func checkLex(t *testing.T, l *lexer, expectedTypes []tokenType, expectedValues []string) {
	var nextToken token
	typesForPrinting := make([]string, len(expectedTypes))
	for i, arg := range expectedTypes {
		typesForPrinting[i] = string(arg)
	}
	for index, expectedType := range expectedTypes {
		nextToken = <-l.tokens

		if expectedType != nextToken.typ {
			t.Fatalf("The token types did not match for token %+v. Seen so far: [%v]", nextToken, strings.Join(typesForPrinting[0:index+1], ","))
		}
		if expectedValues[index] != nextToken.val {
			t.Fatalf("The token types did not match for token %+v. Seen so far: [%v]", nextToken, strings.Join(expectedValues[0:index+1], ","))
		}
	}
	nextToken = <-l.tokens
	if nextToken.typ != tokenEOF {
		t.Fatal("Expected the last token to be eof")
	}
}

func TestImportStatements(t *testing.T) {
	expectedTypes := []tokenType{tokenImport, tokenString}
	expectedValues := []string{"import", "github.com/anthonybishopric/mybin"}
	checkLex(t, lex("an import statement", "import \"github.com/anthonybishopric/mybin\""), expectedTypes, expectedValues)
}
