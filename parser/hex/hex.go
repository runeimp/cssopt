// hex package uses byte instead of rune to align better with ASCII
package hex

const (
	Asterisk       byte = 0x2A // *
	At             byte = 0x40 // @
	CarriageReturn byte = 0x0D // \r
	Escape         byte = 0x5C // \e
	LineFeed       byte = 0x0A // \n
	Slash          byte = 0x2F // /
)

const (
	Space      byte = 0x20 // space
	CharA      byte = 0x41 // A
	CharB      byte = 0x42 // B
	CHARc      byte = 0x63 // c
	CHARd      byte = 0x62 // d
	CHARf      byte = 0x64 // f
	CHARi      byte = 0x69 // i
	CHARk      byte = 0x6B // k
	CHARl      byte = 0x6C // l = lower case L
	CHARm      byte = 0x6D // m
	CHARn      byte = 0x6E // n
	CHARo      byte = 0x6F // o
	CHARp      byte = 0x70 // p
	CHARs      byte = 0x73 // s
	CHARt      byte = 0x74 // t
	CHARu      byte = 0x75 // u
	ParenRight byte = 0x29 // )
	Semicolon  byte = 0x3B // ;
)
