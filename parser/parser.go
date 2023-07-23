package parser

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/runeimp/cssopt/configuration"
	"github.com/runeimp/cssopt/parser/lexer"
	"github.com/runeimp/termlog"
)

var (
	regexCSSimports = regexp.MustCompile(`@import +(?:url\()?["']?([0-9a-zA-Z_:/.?-]+)`)
	regexAtCharset  = regexp.MustCompile(`@charset ['"]?[a-zA-Z0-9-]+['"]?;?`)
	tlog            = termlog.New()
)

type CssFile struct {
	body            string
	imports         map[string]AtImport
	path            string
	processComplete bool
	source          []byte
}

func (cf *CssFile) AddImport(ai AtImport) {
	cf.imports[ai.name] = ai
}

func (cf *CssFile) GetBody() string {
	return cf.body
}

func (cf *CssFile) GetBodyLength() int {
	return len(cf.body)
}

func (cf *CssFile) GetImports() map[string]AtImport {
	return cf.imports
}

func (cf *CssFile) GetPath() string {
	return cf.path
}

func (cf *CssFile) ReplaceImport(key, value string) string {
	value = regexAtCharset.ReplaceAllString(value, "")
	cf.body = strings.Replace(cf.body, key, value, 1)
	return cf.body
}

func (cf *CssFile) String() string {
	dict := make(map[string]any)
	dict["path"] = cf.path
	dict["complete"] = cf.processComplete  // fmt.Sprintf("%t", cf.processComplete)
	dict["length_body"] = len(cf.body)     // fmt.Sprintf("%d", len(cf.body))
	dict["length_source"] = len(cf.source) // fmt.Sprintf("%d", len(cf.source))
	// dict["length_imports"] = len(cf.imports) // fmt.Sprintf("%d", len(cf.source))

	jsonData, err := json.MarshalIndent(dict, "", "\t")
	if err != nil {
		tlog.Fatal(err)
	}
	return string(jsonData)
}

type AtImport struct {
	name   string
	path   string
	source string
}

func (ai *AtImport) GetPath() string {
	return ai.path
}

func (ai *AtImport) String() string {
	m := make(map[string]string)
	m["name"] = ai.name
	m["path"] = ai.path
	m["source"] = ai.source
	bytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type ParserCSS struct {
	bytes   []byte
	Config  *configuration.Config
	imports map[string]AtImport
	path    string
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

func (p *ParserCSS) GetImports() map[string]AtImport {
	return p.imports
}

func (p *ParserCSS) Run(bytes ...[]byte) (css *CssFile, err error) {
	tlog.Debug("parserCSS.Run() | len(bytes): %d", len(bytes))
	var (
		firstPass  []lexer.Token
		secondPass []lexer.Token
	)

	css = &CssFile{
		imports: make(map[string]AtImport),
		path:    p.path,
	}

	if len(bytes) > 0 {
		p.bytes = bytes[0]
		css.source = bytes[0] // NOTE: probably don't need this
	} else {
		p.bytes, err = p.getFile()
		if err != nil {
			return nil, err
		}
	}
	tlog.Info("parserCSS.Run() | len(p.bytes): %d", len(p.bytes))

	lex := lexer.New()
	lex.Run(p.bytes)
	firstPass = lex.Tokens

	tlog.Debug("parserCSS.Run() | len(firstPass): %d", len(firstPass))
	tlog.Debug("parserCSS.Run() | p.Config.Comments: %#v", p.Config.Comments)
	tlog.Debug("parserCSS.Run() | p.Config.Newlines: %q", p.Config.Newlines)
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
			// case lexer.TokenAtCharSet:
			// 	if i == 0 {
			// 		secondPass = append(secondPass, tok)
			// 	}
			case lexer.TokenAtImport:
				matches := regexCSSimports.FindAllSubmatch(tok.Value, -1)
				for _, m := range matches {
					if len(m) > 1 && len(m[1]) > 0 {
						tlog.Debug("parserCSS.Run() | @import %q", m[1])
						tok.Value = m[1]
						secondPass = append(secondPass, tok)
					} else {
						tlog.Error("parserCSS.Run() | @import %q (%T)", m, m)
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
			impValue := string(tok.Value)
			tlog.Debug("parserCSS.Run() | import: %q", impValue)
			css.body += "@import " + impValue + ";"
			// css.imports = append(css.imports, p.path+"/"+impValue)
			// p.imports = append(p.imports, path.Dir(p.path)+"/"+impValue)
			/*

				BREATH

			*/
			imp := AtImport{
				path:   path.Join(path.Dir(p.path), impValue),
				name:   impValue,
				source: p.path,
			}
			p.imports[imp.name] = imp
			/*

				BREATH

			*/
		case lexer.TokenLineFeed:
			switch p.Config.Newlines {
			case "remove":
				// Adiós!
			case "merge":
				if lastToken.Type != lexer.TokenLineFeed && lastToken.InComment == false {
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

	// tlog.Warn("parserCSS.Run() | len(lex.Tokens): %d", len(lex.Tokens))
	// tlog.Warn("parserCSS.Run() | css.body:\n%s", css.body)

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

	for k, v := range p.imports {
		tlog.Debug("parserCSS.Run() | %q | imp: %s", k, v.String())
		// Go Routine those imports
	}

	return css, err
}

func NewCSS(path ...string) (proc *ParserCSS) {
	proc = &ParserCSS{
		imports: map[string]AtImport{},
	}
	if len(path) > 0 {
		proc.path = path[0]
	}
	return
}

func init() {
	tlog.Level = termlog.InfoLevel
}
