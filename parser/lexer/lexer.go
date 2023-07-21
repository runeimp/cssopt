package lexer

import (
	"fmt"

	"github.com/runeimp/cssopt/parser/hex"
	"github.com/runeimp/termlog"
)

/*
const (

	hexAsterisk       = 0x2A
	hexAt             = 0x40
	hexCarriageReturn = 0x0D
	hexEscape         = 0x5C
	hexLineFeed       = 0x0A
	hexSlash          = 0x2F

)

const (

	hexSpace      = 0x20
	hexCHARc      = 0x63
	hexCHARd      = 0x62
	hexCHARf      = 0x64
	hexCHARi      = 0x69
	hexCHARk      = 0x6B
	hexCHARl      = 0x6C
	hexCHARm      = 0x6D
	hexCHARn      = 0x6E
	hexCHARo      = 0x6F
	hexCHARp      = 0x70
	hexCHARs      = 0x73
	hexCHARt      = 0x74
	hexCHARu      = 0x75
	hexParenRight = 0x29
	hexSemicolon  = 0x3B

)
*/
const (
	TokenCarriageReturn      TokenType = "Carriage Return"
	TokenAtRule              TokenType = "At-Rule"
	TokenAtCharSet           TokenType = "Character Set"
	TokenAtColorProfile      TokenType = "Color Profile"
	TokenAtContainer         TokenType = "Container"
	TokenAtCounterStyle      TokenType = "Counter Style"
	TokenAtDocument          TokenType = "Document"
	TokenAtFontFace          TokenType = "Font Face"
	TokenAtFontFeatureValues TokenType = "Font Feature Values"
	TokenAtFontPaletteValues TokenType = "Font Palette Values"
	TokenAtImport            TokenType = "Import"
	TokenAtKeyFrames         TokenType = "Key Frames"
	TokenAtLayer             TokenType = "Layer"
	TokenAtMedia             TokenType = "Media"
	TokenAtNameSpace         TokenType = "Name Space"
	TokenAtPage              TokenType = "Page"
	TokenAtProperty          TokenType = "Property"
	TokenAtSupports          TokenType = "Supports"
	TokenLineFeed            TokenType = "Line Feed"
	TokenText                TokenType = "Text"
	TokenUnknown             TokenType = "Token Unknown"
)

type TokenType string

type Token struct {
	InComment bool
	Type      TokenType
	Value     []byte
}

func (t Token) String() string {
	msg := fmt.Sprintf("%s: ", string(t.Type))
	switch t.Type {
	case TokenCarriageReturn:
		msg += `\r`
	case TokenLineFeed:
		msg += `\n`
	default:
		for _, v := range t.Value {
			msg += string(v)
		}
	}

	if t.InComment {
		msg += `(in comment)`
	}

	return msg
}

type Lexer struct {
	source []byte
	Tokens []Token
}

// func (lex Lexer) GetTokens() []Token {
// 	return lex.Tokens
// }

func (lex *Lexer) Run(src []byte) {
	lex.source = src

	var (
		inAtRule  bool
		inComment bool
		j         int
		lastChar  byte
		nextChar  byte
		srcLength = len(lex.source)
		tok       = Token{Type: TokenText}
	)

	tlog.Info("lexer.Run() | len(lex.source): %d", len(lex.source))

	for i, b := range lex.source {
		tlog.Debug("lexer.Run() | %05d | 0x%02X", i, b)

		j = i + 1
		if j < srcLength {
			nextChar = lex.source[j]
		}

		switch b {
		case hex.At:
			// tlog.Info("lexer.Run() | %05d | 0x%02X | At-Rule", i, b)
			switch nextChar {
			case hex.CHARc:
				part := string(lex.source[i : i+4])
				tlog.Debug("lexer.Run() | At-Rule | %05d | 0x%02X | part: %q", i, b, part)
				switch part {
				case "@cha":
					// tok.Type = TokenAtCharSet
					tok.Type = TokenAtRule
				case "@col":
					// tok.Type = TokenAtColorProfile
					tok.Type = TokenAtRule
				case "@con":
					// tok.Type = TokenAtContainer
					tok.Type = TokenAtRule
				case "@cou":
					// tok.Type = TokenAtCounterStyle
					tok.Type = TokenAtRule
				}
			case hex.CHARd:
				tok.Type = TokenAtDocument
			case hex.CHARf:
				part := string(lex.source[i : i+12])
				tlog.Debug("lexer.Run() | At-Rule | %05d | 0x%02X | part: %q", i, b, part)
				switch part {
				case "@font-face":
					// tok.Type = TokenAtFontFace
					tok.Type = TokenAtRule
				case "@font-feat":
					// tok.Type = TokenAtFontFeatureValues
					tok.Type = TokenAtRule
				case "@font-pale":
					// tok.Type = TokenAtFontPaletteValues
					tok.Type = TokenAtRule
				}
			case hex.CHARi:
				part := string(lex.source[i : i+7])
				if part == "@import" {
					tok.Type = TokenAtImport
					// inAtRule = true
				}
			case hex.CHARk:
				part := string(lex.source[i : i+10])
				if part == "@keyframes" {
					// tok.Type = TokenAtKeyFrames
					tok.Type = TokenAtRule
					// inAtRule = true
				}
			case hex.CHARl:
				part := string(lex.source[i : i+6])
				if part == "@layer" {
					// tok.Type = TokenAtLayer
					tok.Type = TokenAtRule
				}
			case hex.CHARm:
				part := string(lex.source[i : i+6])
				if part == "@media" {
					// tok.Type = TokenAtMedia
					tok.Type = TokenAtRule
				}
			case hex.CHARn:
				part := string(lex.source[i : i+10])
				if part == "@namespace" {
					// tok.Type = TokenAtNameSpace
					tok.Type = TokenAtRule
				}
			case hex.CHARp:
				part := string(lex.source[i : i+5])
				switch part {
				case "@page":
					// tok.Type = TokenAtPage
					tok.Type = TokenAtRule
				case "@prop":
					// tok.Type = TokenAtProperty
					tok.Type = TokenAtRule
				}
			case hex.CHARs:
				part := string(lex.source[i : i+9])
				if part == "@supports" {
					// tok.Type = TokenAtSupports
					tok.Type = TokenAtRule
				}
			default:
				if nextChar != hex.Space {
					tok.Type = TokenAtRule
				}
			}
			inAtRule = true
			tok.Value = append(tok.Value, b)
		case hex.CarriageReturn:
			tlog.Debug("lexer.Run() | %05d | 0x%02X (CR) | inComment: %-5t", i, b, inComment)
			if len(tok.Value) > 0 {
				lex.Tokens = append(lex.Tokens, tok)
				tok.Value = nil
			}
			tok.Type = TokenCarriageReturn
			tok.Value = append(tok.Value, b)
			// if inComment {
			// 	tok.Type = TokenCommentCR
			// }
			lex.Tokens = append(lex.Tokens, tok)
			tok.Type = TokenText
			// if inComment {
			// 	// tok.Type = TokenCommentText
			// 	tok.Type = TokenComment
			// }
			tok.InComment = inComment
			tok.Value = nil
			// IGNORE: most basic minification as CRLF is not required for a newline in most new operating systems including Windows 10+
		case hex.LineFeed:
			tlog.Debug("lexer.Run() | %05d | 0x%02X (CR) | inComment: %-5t", i, b, inComment)
			if len(tok.Value) > 0 {
				lex.Tokens = append(lex.Tokens, tok)
				tok.Value = nil
			}
			tok.Type = TokenLineFeed
			tok.Value = append(tok.Value, b)
			// if inComment {
			// 	tok.Type = TokenCommentLF
			// }
			lex.Tokens = append(lex.Tokens, tok)
			tok.Type = TokenText
			// if inComment {
			// 	// tok.Type = TokenCommentText
			// 	tok.Type = TokenComment
			// }
			tok.InComment = inComment
			tok.Value = nil
		case hex.Semicolon:
			if !inComment {
				tok.Value = append(tok.Value, b)
				lex.Tokens = append(lex.Tokens, tok)
				tok.Type = TokenText
				tok.Value = nil
			}
		case hex.Slash:
			tlog.Debug("lexer.Run() | %05d | 0x%02X /    | inComment: %-5t | inAtRule: %-5t", i, b, inComment, inAtRule)

			if inComment {
				if lastChar == hex.Asterisk {
					inComment = false

					tok.Value = append(tok.Value, b)
					lex.Tokens = append(lex.Tokens, tok)

					tlog.Debug("lexer.Run() | %05d | tokens: %05d | tok: %+v", i, len(lex.Tokens), tok)

					tok.Type = TokenUnknown
					tok.Value = nil
				}
			} else if !inComment {
				if nextChar == hex.Asterisk {
					if len(tok.Value) > 0 {
						tlog.Debug("lexer.Run() | %05d | tokens: %05d | tok: %s (transitioning to comment block)", i, len(lex.Tokens), tok)
						lex.Tokens = append(lex.Tokens, tok)
						tok.Value = nil
					}
					inComment = true
					tok.Type = TokenText
					// tok.Type = TokenComment
				}
				tok.Value = append(tok.Value, b)
			}
			tok.InComment = inComment

			// tlog.Info("lexer.Run() | %05d | tokens: %05d | tok: %s", i, len(lex.Tokens), tok)

			// tok.Value = append(tok.Value, b)
		default:
			tok.InComment = inComment
			tok.Value = append(tok.Value, b)
		}

		lastChar = b
	}
	tlog.Info("lexer.Run() | len(lex.Tokens): %d", len(lex.Tokens))
}

func New() *Lexer {
	return &Lexer{
		Tokens: []Token{},
	}
}

var (
	tlog *termlog.Logger
)

func init() {
	tlog = termlog.New()
	tlog.Level = termlog.InfoLevel
}
