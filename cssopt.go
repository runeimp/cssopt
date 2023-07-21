package cssopt

import (
	"io/fs"
	"os"
	"time"

	"github.com/runeimp/cssopt/configuration"
	"github.com/runeimp/cssopt/parser"
	"github.com/runeimp/termlog"
)

const (
	AppVersion = "0.1.0-alpha+001"
	AppName    = "CSS Optimizer"
	AppLabel   = AppName + " v" + AppVersion
)

var (
	tlog      = termlog.New()
	optimizer *Optimizer
)

type Optimizer struct {
	cache     string
	config    *configuration.Config
	processed time.Time
	source    []string
}

func (opt *Optimizer) processFilePath(path string) (result *parser.CssFile, err error) {
	var fileBytes []byte

	tlog.Level = termlog.WarnLevel
	tlog.Info("optimizer.processFilePath() | path: %q", path)

	// fileBytes, err = os.ReadFile(path)
	// if err != nil {
	// 	tlog.Fatal(err)
	// 	return result, err
	// }

	// result = CssFile{
	// 	path:   path,
	// 	source: fileBytes,
	// }

	tlog.Debug("optimizer.processFilePath() | fileBytes: %q", string(fileBytes))

	proc := parser.NewCSS(path)
	proc.Config = opt.config
	result, err = proc.Run()
	if err != nil {
		return result, err
	}

	/*
		lex := lexer.New(fileBytes)
		lex.Lex()

		tlog.Info("optimizer.processFilePath() | len(lex.Tokens): %d", len(lex.Tokens))
		// tlog.Info("optimizer.processFilePath() | len(lex.GetTokens()): %d", len(lex.GetTokens()))

		// result.processComplete = true

		const infoFormat = "optimizer.processFilePath(%q) | comment: %-5t | %-9s |	 %s"
		for _, tok := range lex.Tokens {
			// tlog.Info("optimizer.processFilePath(%q) | i: %03d | tok: %s", path, i, tok)
			// tlog.Info("optimizer.processFilePath(%q) | tok: %s", path, tok)
			// msg := fmt.Sprintf("%17s: ", string(t.Type))

			switch tok.Type {
			case lexer.TokenCarriageReturn:
				tlog.Info(infoFormat, path, tok.InComment, tok.Type, `\r`)
			case lexer.TokenLineFeed:
				tlog.Info(infoFormat, path, tok.InComment, tok.Type, `\n`)
			default:
				tlog.Info(infoFormat, path, tok.InComment, tok.Type, string(tok.Value))
			}
		}
	*/
	// os.Stdout.Write(data.File)

	return result, err
}

func (opt *Optimizer) ProcessSliceOfStrings(cssList []string) (result string) {

	return result
}

func (opt *Optimizer) ProcessString(css string) (result string) {

	return result
}

func (opt *Optimizer) ProcessSliceOfBytes(cssBytes []byte) (result string) {

	return result
}

func (opt *Optimizer) ProcessSlicesOfBytes(css [][]byte) (result string) {

	return result
}

func (opt *Optimizer) ProcessPath(path string) (result string, err error) {
	var (
		// cssFile *parser.CssFile
		finfo fs.FileInfo
	)
	tlog.Info("optimizer.ProcessPath() | path: %q", path)

	finfo, err = pathInfo(path)
	if err != nil {
		tlog.Error(err.Error())
		return "", err
	}

	if finfo.IsDir() {
		// Loop over files

		return result, err
	}

	// cssFile, err = opt.processFilePath(path)
	_, err = opt.processFilePath(path)

	// for i, imp := range cssFile.imports {
	// 	tlog.Info("optimizer.ProcessPath() | %03d | imp: %q", i, imp)
	// 	// Go Routine those imports
	// }

	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	tlog.Fatal(err)
	// }
	// os.Stdout.Write(data.File)
	return result, err
}

func init() {
	// tlog = termlog.New()
	tlog.Level = termlog.InfoLevel
}

func GetOptimizer(conf *configuration.Config) *Optimizer {
	if optimizer == nil {
		optimizer = &Optimizer{
			config: conf,
		}
	}

	return optimizer
}

func pathInfo(path string) (fs.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		tlog.Error(err.Error())
		return nil, err
	}

	defer file.Close()

	finfo, err := file.Stat()
	if err != nil {
		tlog.Error(err.Error())
		return nil, err
	}

	return finfo, nil
}

//
