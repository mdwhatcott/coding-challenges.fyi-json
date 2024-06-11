package lexing

import "bytes"

type TokenType string

const (
	TokenNull         TokenType = "<null>"
	TokenTrue         TokenType = "<true>"
	TokenFalse        TokenType = "<false>"
	TokenNegativeSign TokenType = "<->"
	TokenDecimalPoint TokenType = "<.>"
	TokenZero         TokenType = "<0>"
	TokenOne          TokenType = "<1>"
	TokenTwo          TokenType = "<2>"
	TokenThree        TokenType = "<3>"
	TokenFour         TokenType = "<4>"
	TokenFive         TokenType = "<5>"
	TokenSix          TokenType = "<6>"
	TokenSeven        TokenType = "<7>"
	TokenEight        TokenType = "<8>"
	TokenNine         TokenType = "<9>"
)

type Token struct {
	Type  TokenType
	Value []byte
}

type stateMethod func() stateMethod

type Lexer struct {
	input  []byte
	start  int // start position of this item.
	pos    int // current position in the input.
	width  int // width of last rune read from input.
	output chan Token
}

func New(input []byte) *Lexer { // TODO: accept io.Reader?
	return &Lexer{
		input:  input,
		output: make(chan Token),
	}
}

func (this *Lexer) Output() <-chan Token {
	return this.output
}

func (this *Lexer) Lex() {
	defer close(this.output)
	if len(this.input) == 0 {
		return
	}
	if isWhiteSpace(this.at(0)) {
		return
	}
	for state := this.lexValue; state != nil && this.pos < len(this.input); {
		state = state()
	}
}

func (this *Lexer) at(offset int) rune {
	return rune(this.input[this.pos+offset])
}
func (this *Lexer) emit(tokenType TokenType) {
	this.output <- Token{Type: tokenType, Value: this.input[this.start:this.pos]}
	this.start = this.pos
}

func (this *Lexer) lexValue() stateMethod {
	if bytes.HasPrefix(this.input, _null) {
		this.pos += len(_null)
		this.emit(TokenNull)
	} else if bytes.HasPrefix(this.input, _true) {
		this.pos += len(_true)
		this.emit(TokenTrue)
	} else if bytes.HasPrefix(this.input, _false) {
		this.pos += len(_false)
		this.emit(TokenFalse)
	} else if this.at(0) == negativeSign {
		return this.lexNumberFromNegativeSign
	} else if this.at(0) == _0 {
		return this.lexNumberFromZero
	} else if isDigit(this.at(0)) {
		return this.lexNumberFromDigit
	}
	return nil
}
func (this *Lexer) lexNumberFromNegativeSign() stateMethod {
	this.pos++
	this.emit(TokenNegativeSign)
	return this.lexNumberFromDigit
}
func (this *Lexer) lexNumberFromZero() stateMethod {
	this.pos++
	this.emit(TokenZero)
	return this.lexFraction
}
func (this *Lexer) lexNumberFromDigit() stateMethod {
	this.pos++
	this.emit(digitToken(rune(this.input[this.start])))
	this.emitDigits()
	return this.lexFraction
}
func (this *Lexer) lexFraction() stateMethod {
	if this.at(0) == '.' && isDigit(this.at(1)) {
		this.pos++
		this.emit(TokenDecimalPoint)
		this.emitDigits()
	}
	return nil
}
func (this *Lexer) emitDigits() {
	for isDigit(this.at(0)) {
		this.pos++
		this.emit(digitToken(rune(this.input[this.start])))
	}
}

func isWhiteSpace(r rune) bool {
	return r == ' ' // TODO: additional whitespace characters
}
func isDigit(r rune) bool {
	return digitToken(r) != ""
}

func digitToken(r rune) TokenType {
	switch r {
	case _0:
		return TokenZero
	case _1:
		return TokenOne
	case _2:
		return TokenTwo
	case _3:
		return TokenThree
	case _4:
		return TokenFour
	case _5:
		return TokenFive
	case _6:
		return TokenSix
	case _7:
		return TokenSeven
	case _8:
		return TokenEight
	case _9:
		return TokenNine
	default:
		return ""
	}
}

var (
	_null  = []byte("null")
	_true  = []byte("true")
	_false = []byte("false")
)

const (
	negativeSign = '-'

	_0 = '0'
	_1 = '1'
	_2 = '2'
	_3 = '3'
	_4 = '4'
	_5 = '5'
	_6 = '6'
	_7 = '7'
	_8 = '8'
	_9 = '9'
)
