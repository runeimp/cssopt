package parser

import (
	"os"

	"github.com/runeimp/cssopt/parser/lexer"
	"github.com/runeimp/termlog"
)

var tlog = termlog.New()

type CssFile struct {
	imports         []string
	path            string
	processComplete bool
	source          []byte
}

type ParserCSS struct {
	bytes []byte
	path  string
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
	css = &CssFile{}
	tlog := termlog.New()

	if len(bytes) > 0 {
		p.bytes = bytes[0]
	} else {
		p.bytes, err = p.getFile()
	}

	lex := lexer.New()
	lex.Run(p.bytes)

	tlog.Info("parserCSS.Run() | len(lex.Tokens): %d", len(lex.Tokens))
	// tlog.Info("parserCSS.Run() | len(lex.GetTokens()): %d", len(lex.GetTokens()))

	// result.processComplete = true

	const infoFormat = "parserCSS.Run(%q) | comment: %-5t | %-9s |	 %s"
	for _, tok := range lex.Tokens {
		// tlog.Info("parserCSS.Run(%q) | i: %03d | tok: %s", path, i, tok)
		// tlog.Info("parserCSS.Run(%q) | tok: %s", path, tok)
		// msg := fmt.Sprintf("%17s: ", string(t.Type))

		switch tok.Type {
		case lexer.TokenCarriageReturn:
			tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, `\r`)
		case lexer.TokenLineFeed:
			tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, `\n`)
		default:
			tlog.Info(infoFormat, p.path, tok.InComment, tok.Type, string(tok.Value))
		}
	}

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
