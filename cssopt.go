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
	var (
		files []*parser.CssFile
		res   *parser.CssFile
	)

	tlog.Level = termlog.WarnLevel
	tlog.Info("optimizer.processFilePath() | path: %q", path)

	// fileBytes, err = os.ReadFile(path)
	// if err != nil {
	// 	tlog.Fatal(err)
	// 	return result, err
	// }

	// tlog.Debug("optimizer.processFilePath() | fileBytes: %q", string(fileBytes))

	result, err = parseCssFiles(path, opt.config)
	if err != nil {
		return result, err
	}
	tlog.Warn("optimizer.processFilePath() | path: %q | res.GetBodyLength(): %d", path, result.GetBodyLength())
	// files = append(files, result)

	tlog.Warn("optimizer.processFilePath() | res: %s", result)
	imports := result.GetImports()
	tlog.Warn("optimizer.processFilePath() | len(imports): %d", len(imports))

	for k, v := range imports {
		// tlog.Warn("optimizer.processFilePath() | %q | import: %s", k, v.String())
		tlog.Info("optimizer.processFilePath() | %q | path: %q", k, v.GetPath())
		res, err = parseCssFiles(v.GetPath(), opt.config)
		if err != nil {
			return result, err
		}
		tlog.Info("optimizer.processFilePath() | %q | res.GetBodyLength(): %d", k, res.GetBodyLength())
		files = append(files, res)
		tlog.Info("optimizer.processFilePath() | %q | len(imports): %d", k, len(imports))
	}

	// tlog.Error("optimizer.processFilePath() | Results | len(result): %d", len(result))

	// Merge Imports?
	if opt.config.Imports {
		for i, file := range files {
			// tlog.Error("optimizer.processFilePath() | %d | file.GetPath(): %q\n%s", i, file.GetPath(), file.GetBody())
			tlog.Error("optimizer.processFilePath() | %d | file.GetPath(): %q", i, file.GetPath())

			for k, v := range result.GetImports() {
				if file.GetPath() == v.GetPath() {
					target := "@import " + k
					tlog.Warn("optimizer.processFilePath() | %q | %q", target, v.GetPath())
					result.ReplaceImport(target, file.GetBody())
				}
			}
		}
	}

	// tlog.Error("optimizer.processFilePath() | %q | %s", result[0].GetPath(), result[0].GetBody())

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
		cssFile *parser.CssFile
		finfo   fs.FileInfo
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

	cssFile, err = opt.processFilePath(path)
	// _, err = opt.processFilePath(path)

	result = cssFile.GetBody()

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
	tlog.Level = termlog.DebugLevel
}

func GetOptimizer(conf *configuration.Config) *Optimizer {
	if optimizer == nil {
		optimizer = &Optimizer{
			config: conf,
		}
	}

	return optimizer
}

func parseCssFiles(path string, config *configuration.Config) (result *parser.CssFile, err error) {
	tlog.Warn("cssopt.parseCssFiles() | path: %q", path)
	proc := parser.NewCSS(path)
	proc.Config = config
	result, err = proc.Run()
	if err != nil {
		return result, err
	}

	for _, imp := range proc.GetImports() {
		result.AddImport(imp)
	}

	return result, err
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
