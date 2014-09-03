package lexer

import (
	"strings"
	"testing"
)

func checkLex(t *testing.T, tokens tokenChan, expectedTypes []tokenType, expectedValues []string) {
	var nextToken token
	typesForPrinting := make([]string, len(expectedTypes))
	for i, arg := range expectedTypes {
		typesForPrinting[i] = string(arg)
	}
	for index, expectedType := range expectedTypes {
		nextToken = <-tokens
		if !nextToken.Valid() {
			t.Fatalf("The token stream was closed prematurely. Seen so far: [%v]", strings.Join(typesForPrinting[0:index+1], ","))
		}

		if expectedType != nextToken.typ {
			t.Fatalf("The token types did not match for token %+v. Seen so far: [%v]", nextToken, strings.Join(typesForPrinting[0:index+1], ","))
		}
		if expectedValues[index] != nextToken.val {
			t.Fatalf("The token types did not match for token %+v. Seen so far: [%v]", nextToken, strings.Join(expectedValues[0:index+1], ","))
		}
	}
	nextToken = <-tokens
	if nextToken.typ != tokenEOF {
		t.Fatal("Expected the last token to be eof")
	}
}

func TestImportStatements(t *testing.T) {
	expectedTypes := []tokenType{tokenImport, tokenString}
	expectedValues := []string{"import", "github.com/anthonybishopric/mybin"}
	checkLex(t, lex("import \"github.com/anthonybishopric/mybin\""), expectedTypes, expectedValues)
}

func TestImportStatementWithoutSpaceIsErr(t *testing.T) {
	// expectedTypes := []tokenType{tokenError}
	// tokenChannel := lex("import\"github.com/anthonybishorpic/mybin\"")

}
