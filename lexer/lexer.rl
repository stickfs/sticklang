package lexer

import (
	"fmt"
)

%%{
	machine test_lexer;

	newline = "\n" @{ worker.currentLine++ };
	integer = ('+'|'-')?[0-9]+ @{ fmt.Println("Integer")};
	comma = "," @{ fmt.Println(",") };

	main := (integer comma)+ newline;

}%%

%% write data;

func lex(data string) tokenChan {
	worker := newWorker()

	cs, p, pe := 0, 0, len(data)

	%% write init;

	go func() {
		%% write exec;
		close(worker.out)
	}()
	return worker.out
}
