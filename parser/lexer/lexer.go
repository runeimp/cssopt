package lexer

import (
	"fmt"

	"github.com/runeimp/termlog"
)

const (
	hexAsterisk       = 0x2A
	hexCarriageReturn = 0x0D
	hexEscape         = 0x5C
	hexLineFeed       = 0x0A
	hexSlash          = 0x2F
)

const (
	TokenCarriageReturn TokenType = "Carriage Return"
	// TokenComment        TokenType = "Comment"
	// TokenCommentText    TokenType = "Comment Text"
	// TokenCommentCR TokenType = "Comment Carriage Return"
	// TokenCommentLF TokenType = "Comment Line Feed"
	TokenLineFeed TokenType = "Line Feed"
	TokenText     TokenType = "Text"
	TokenUnknown  TokenType = "Token Unknown"
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
		inComment     bool
		inLineComment bool
		j             int
		lastChar      byte
		nextChar      byte
		srcLength     = len(lex.source)
		tok           = Token{Type: TokenText}
	)

	tlog.Info("lexer.Run() | len(lex.source): %d", len(lex.source))

	for i, b := range lex.source {
		tlog.Debug("lexer.Run() | %05d | 0x%02X", i, b)

		j = i + 1
		if j < srcLength {
			nextChar = lex.source[j]
		}

		switch b {
		case hexCarriageReturn:
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
		case hexLineFeed:
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
		case hexSlash:
			tlog.Debug("lexer.Run() | %05d | 0x%02X /    | inComment: %-5t | inLineComment: %-5t", i, b, inComment, inLineComment)

			if inComment {
				if lastChar == hexAsterisk {
					inComment = false

					tok.Value = append(tok.Value, b)
					lex.Tokens = append(lex.Tokens, tok)

					tlog.Debug("lexer.Run() | %05d | tokens: %05d | tok: %+v", i, len(lex.Tokens), tok)

					tok.Type = TokenUnknown
					tok.Value = nil
				}
			} else if !inComment {
				if nextChar == hexAsterisk {
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
