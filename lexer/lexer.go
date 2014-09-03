
// line 1 "lexer.rl"
package lexer

import (
	"fmt"
)


// line 16 "lexer.rl"



// line 15 "lexer.go"
var _test_lexer_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 
}

var _test_lexer_key_offsets []byte = []byte{
	0, 0, 4, 6, 9, 14, 
}

var _test_lexer_trans_keys []byte = []byte{
	43, 45, 48, 57, 48, 57, 44, 48, 
	57, 10, 43, 45, 48, 57, 
}

var _test_lexer_single_lengths []byte = []byte{
	0, 2, 0, 1, 3, 0, 
}

var _test_lexer_range_lengths []byte = []byte{
	0, 1, 1, 1, 1, 0, 
}

var _test_lexer_index_offsets []byte = []byte{
	0, 0, 4, 6, 9, 14, 
}

var _test_lexer_indicies []byte = []byte{
	0, 0, 2, 1, 2, 1, 3, 2, 
	1, 4, 0, 0, 2, 1, 1, 
}

var _test_lexer_trans_targs []byte = []byte{
	2, 0, 3, 4, 5, 
}

var _test_lexer_trans_actions []byte = []byte{
	0, 0, 3, 5, 1, 
}

const test_lexer_start int = 1
const test_lexer_first_final int = 5
const test_lexer_error int = 0

const test_lexer_en_main int = 1


// line 19 "lexer.rl"

func lex(data string) tokenChan {
	worker := newWorker()

	cs, p, pe := 0, 0, len(data)

	
// line 69 "lexer.go"
	{
	cs = test_lexer_start
	}

// line 26 "lexer.rl"

	go func() {
		
// line 78 "lexer.go"
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_test_lexer_key_offsets[cs])
	_trans = int(_test_lexer_index_offsets[cs])

	_klen = int(_test_lexer_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _test_lexer_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _test_lexer_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_test_lexer_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _test_lexer_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _test_lexer_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_test_lexer_indicies[_trans])
	cs = int(_test_lexer_trans_targs[_trans])

	if _test_lexer_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_test_lexer_trans_actions[_trans])
	_nacts = uint(_test_lexer_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _test_lexer_actions[_acts-1] {
		case 0:
// line 10 "lexer.rl"

 worker.currentLine++ 
		case 1:
// line 11 "lexer.rl"

 fmt.Println("Integer")
		case 2:
// line 12 "lexer.rl"

 fmt.Println(",") 
// line 169 "lexer.go"
		}
	}

_again:
	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	_out: {}
	}

// line 29 "lexer.rl"
		close(worker.out)
	}()
	return worker.out
}
