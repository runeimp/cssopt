package parser

import (
	// "fmt"
	"os"
	"regexp"

	"github.com/runeimp/cssopt/configuration"
	"github.com/runeimp/cssopt/parser/lexer"
	"github.com/runeimp/termlog"
)

var (
	regexCSSimports = regexp.MustCompile(`@import +(?:url\()?["']?([0-9a-zA-Z_:/.?-]+)`)
	tlog            = termlog.New()
)

type CssFile struct {
	body            string
	imports         []string
	path            string
	processComplete bool
	source          []byte
}

type ParserCSS struct {
	bytes  []byte
	Config *configuration.Config
	path   string
}

func (p *ParserCSS) getFile() ([]byte, error) {
	var err error

	p.bytes, err = os.ReadFile(p.path)
	if err != nil {
		tlog.Fatal(err)
		return []byte{}, err
	}

	return p.bytes, nil
}

func (p *ParserCSS) Run(bytes ...[]byte) (css *CssFile, err error) {
	var (
		firstPass  []lexer.Token
		imports    = []string{}
		secondPass []lexer.Token
	)

	css = &CssFile{
		path: p.path,
	}

	tlog := termlog.New()
	tlog.Level = termlog.WarnLevel

	if len(bytes) > 0 {
		p.bytes = bytes[0]
	} else {
		p.bytes, err = p.getFile()
	}

	lex := lexer.New()
	lex.Run(p.bytes)
	firstPass = lex.Tokens

	tlog.Info("parserCSS.Run() | len(firstPass): %d", len(firstPass))
	tlog.Info("parserCSS.Run() | p.Config.Comments: %#v", p.Config.Comments)
	tlog.Info("parserCSS.Run() | p.Config.Newlines: %q", p.Config.Newlines)
	// tlog.Info("parserCSS.Run() | len(lex.GetTokens()): %d", len(lex.GetTokens()))

	var lastToken lexer.Token

	for _, tok := range firstPass {
		tlog.Debug("parserCSS.Run() | tok.InComment: %-5t | tok.Type: %q", tok.InComment, tok.Type)

		switch {
		case tok.InComment:
			if p.Config.Comments.All { // NOTE: temporary code
				// css.body += string(tok.Value)
			} else {
				// tlog.Info("parserCSS.Run() | tok.InComment: %-5t | tok.Type: %s", tok.InComment, tok.Type)
				// css.body += string(tok.Value)
				secondPass = append(secondPass, tok)
			}
		default:
			switch tok.Type {
			case lexer.TokenAtImport:
				if p.Config.Imports {
					imports = append(imports, string(tok.Value))
				}
				matches := regexCSSimports.FindAllSubmatch(tok.Value, -1)
				for _, m := range matches {
					if len(m) > 1 && len(m[1]) > 0 {
						// tlog.Warn("parserCSS.Run() | @import %q", m[1])
						tok.Value = m[1]
						secondPass = append(secondPass, tok)
					} else {
						// tlog.Error("parserCSS.Run() | @import %q (%T)", m, m)
					}
				}
				// match0 := string(matches[0])
				// match1 := string(matches[1])
				// tlog.Warn("parserCSS.Run() | @import | rule: %t | match0: %q | match1: %q | %q", p.Config.Imports, match0, match1, imports)
			case lexer.TokenLineFeed:
				// if lastToken.Type == lexer.TokenLineFeed || lastToken.InComment && p.Config.Comments {
				// tlog.Info("parserCSS.Run() | tok.Type: %v | lastToken.Type: %v | skip: %t", tok.Type, lastToken.Type, (lastToken.Type == lexer.TokenLineFeed))
				switch p.Config.Newlines {
				case configuration.NewlineRemove:
					// Skip: remove all
				case configuration.NewlineNone:
					// remove no newlines
					secondPass = append(secondPass, tok)
				case configuration.NewlineMerge:
					// Merge: keep the first in a series of newlines
					if lastToken.Type != lexer.TokenLineFeed {
						secondPass = append(secondPass, tok)
					}
				case configuration.NewlineLF:
					secondPass = append(secondPass, tok)
				case configuration.NewlineCRLF:
					if lastToken.Type != lexer.TokenCarriageReturn {
						// Add a CR token ... ?
					}
				default:
					tlog.Warn("misconfiguration %q is an invalid option", p.Config.Newlines)
				}
			default:
				// css.body += string(tok.Value)
				secondPass = append(secondPass, tok)
			}
		}
		lastToken = tok
	}

	lastToken.Type = lexer.TokenUnknown
	lastToken.Value = nil

	for _, tok := range secondPass {
		switch tok.Type {
		case lexer.TokenAtImport:
			css.body += "@import " + string(tok.Value) + ";"
		case lexer.TokenLineFeed:
			switch p.Config.Newlines {
			case "remove":
				// Adiós!
			case "merge":
				if lastToken.Type != lexer.TokenLineFeed {
					css.body += string(tok.Value)
				}
			default:
				css.body += string(tok.Value)
			}
		default:
			css.body += string(tok.Value)
		}
		lastToken = tok
	}

	tlog.Warn("parserCSS.Run() | len(lex.Tokens): %d", len(lex.Tokens))
	tlog.Warn("parserCSS.Run() | css.body:\n%s", css.body)

	// result.processComplete = true

	// const infoFormat = "parserCSS.Run(%q) | comment: %-5t | %-9s | %s"
	// for _, tok := range lex.Tokens {
	// 	// tlog.Info("parserCSS.Run(%q) | i: %03d | tok: %s", path, i, tok)
	// 	// tlog.Info("parserCSS.Run(%q) | tok: %s", path, tok)
	// 	// msg := fmt.Sprintf("%17s: ", string(t.Type))

	// 	switch tok.Type {
	// 	case lexer.TokenCarriageReturn:
	// 		tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, `\r`)
	// 	case lexer.TokenLineFeed:
	// 		tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, `\n`)
	// 	default:
	// 		tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, string(tok.Value))
	// 	}
	// }

	// tlog.Info("CSS File:")
	// for _, tok := range lex.Tokens {
	// 	fmt.Print(string(tok.Value))
	// }

	for i, imp := range css.imports {
		tlog.Info("parserCSS.Run() | %03d | imp: %q", i, imp)
		// Go Routine those imports
	}

	return css, err
}

func NewCSS(path ...string) (proc *ParserCSS) {
	proc = &ParserCSS{}
	if len(path) > 0 {
		proc.path = path[0]
	}
	return
}
