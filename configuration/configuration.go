package configuration

import (
	"os"
	"strings"

	"github.com/runeimp/termlog"
)

var tlog = termlog.New()

const (
	EnvCaching   = "CSSOPT_CACHING"
	EnvColors    = "CSSOPT_COLORS"
	EnvComments  = "CSSOPT_COMMENTS"
	EnvGzip      = "CSSOPT_GZIP"
	EnvImports   = "CSSOPT_IMPORTS"
	EnvNewlines  = "CSSOPT_NEWLINES"
	EnvSemicolon = "CSSOPT_SEMICOLON"
	EnvSpaces    = "CSSOPT_SPACES"
	EnvTabs      = "CSSOPT_TABS"
	EnvVars      = "CSSOPT_VARS"
)

type ProcessOption string

type CommentConf struct {
	All    bool
	Body   bool
	Header bool
	Legal  bool
}

const (
	CommentsAll     ProcessOption = "all"     // remote all comments
	CommentsBody    ProcessOption = "body"    // remove body comments
	CommentsHeader  ProcessOption = "header"  // remove header comments
	CommentsLegal   ProcessOption = "legal"   // remove legal comments
	CommentsInvalid ProcessOption = "invalid" // comment setting is invalid
	CommentsNone    ProcessOption = "none"    // remove no comments
)

const (
	NewlineMerge  ProcessOption = "merge"
	NewlineRemove ProcessOption = "remove"
	NewlineCRLF   ProcessOption = "windows"
	NewlineLF     ProcessOption = "posix"
	NewlineNone   ProcessOption = "" // remove no newlines
)

type Config struct {
	Caching       bool          // Enable caching
	Colors        bool          // Optimize color values
	Comments      *CommentConf  // Remove comments
	Gzip          bool          // Enable GZip
	HeaderComment bool          // Save header comments
	Imports       bool          // Merge imports
	LegalComment  bool          // Save legal comments
	Newlines      ProcessOption // "merge", "remove", or "" do nothing
	Semicolon     bool          // Remove last semicolon
	Spaces        bool          // Remove extra spaces
	Tabs          bool          // Remove tabs
	Vars          bool          // Compile vars
}

// func (c *Config) Comments(args ...ProcessOption) []ProcessOption {
// 	if len(args) > 0 {
// 		c.comments = args
// 	}
// 	return c.comments
// }

func New() *Config {
	conf := &Config{
		Caching:   envBool(EnvCaching),
		Colors:    envBool(EnvColors),
		Comments:  envCommentOption(),
		Gzip:      envBool(EnvGzip),
		Imports:   envBool(EnvImports),
		Newlines:  envNewlineOption(),
		Semicolon: envBool(EnvSemicolon),
		Spaces:    envBool(EnvSpaces),
		Tabs:      envBool(EnvTabs),
		Vars:      envBool(EnvVars),
	}
	return conf
}

func envBool(env string) bool {
	e := strings.TrimSpace(strings.ToLower(os.Getenv(env)))
	if len(e) > 0 && e == "true" {
		return true
	}
	return false
}

func envCommentOption() (result *CommentConf) {
	result = &CommentConf{}

	for _, v := range strings.Split(os.Getenv(EnvComments), ",") {
		check := strings.ToLower(strings.TrimSpace(v))
		switch check {
		case "all":
			result.All = true
		case "body":
			result.Body = true
		case "header":
			result.Header = true
		case "legal":
			result.Legal = true
		case "none", "": // doesn't catch here for some reason
			result.All = false
			result.Body = false
			result.Header = false
			result.Legal = false
			return result
		default:
			tlog.Warn("comment option %02X is invalid (%d)", v, len(v))
		}
	}

	return result
}

func envNewlineOption() (result ProcessOption) {
	e := strings.ToLower(strings.TrimSpace(os.Getenv(EnvNewlines)))
	switch e {
	case "merge":
		result = NewlineMerge
	case "remove":
		result = NewlineRemove
	case "windows":
		result = NewlineCRLF
	case "posix":
		result = NewlineLF
	case "":
		result = NewlineNone
	default:
		// result = "invalid"
		tlog.Warn("newline option %q is invalid", os.Getenv(EnvNewlines))
	}
	return result
}
