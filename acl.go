//go:generate ragel -Z -o acl_parser.go acl_parser.rl
//go:generate ragel -V -o acl_parser.dot acl_parser.rl
//go:generate dot -oacl_parser.png -Tpng acl_parser.dot

package archercl

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/op/go-logging"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const (
	// This constant identifies the key under which any build information
	// that has been written to the global variable BuildInfo is stored.
	BUILDINFO_KEY = "buildInfo"

	// A constant that can be used to initialize random number generators
	// for tests so that they produce predictable results. If this is not
	// set during the configuration cascade it will be set by the Load()
	// command at the end to the value of time.Now().UnixNano()
	//
	// Load will set this to the standard random number generator via a
	// rand.seed(...) call, so all you need to do for your tests is put
	// an int in the config file.
	//
	// While this is very useful for most situations, please don't rely on
	// it for anything security related.
	RANDOMSEED_KEY = "randomSeed"

	// If this key is set at to true at the root level, when Load() is
	// done it will log the result of the entire configuration cascade
	// at debug level.
	DUMPCONFIG_KEY = "dumpConfig"

	// If this key is set the root level, whenever the configuration
	// is dumped because DUMPCONFIG_KEY was also set it will be logged
	// with ASCII color codes to make it a little more readable. You
	// usually want this only for debugging and not for dumping into
	// log files.
	DUMPCOLOR_KEY = "dumpColor"

)

var alog = Logger("acl")

// This variable can be set from the init() function of an external module
// to specify build time configuration information that will be added to
// the config node obtained by calling Load(). This string should be a
// valid ACL string or the Load() call will fail.
//
// The idea this is supporting is that on a CI server you can have an
// external build script that has visibility into things like the build
// identifier, git hash, or whatever you want and it can write this
// data into a file like build.go which might look like the following:
//
//     package main
//
//     import "github.com/eyethereal/go-archercl"
//
//     func init() {
//         archercl.BuildInfo = `
//     build_num: 1234
//     git_hash: 3d1004026a5f857551a474a92493c0330b68dd37
//     build_server: coolci.org
//     `
//     }
//
// Such a file is relatively easy to generate from a shell script by
// simply echoing the various parts to it using environment variables.
// Assuming you have done this, your code can then access these variables
// which are effectively baked into your binary by doing something like
//
//     cfg.ChildAsString(archercl.BUILDINFO_KEY, "git_hash")
//
// Any data added to BuildInfo is added after everything else has been
// loaded and is rooted at BUILDINFO_KEY. The implication is that a
// configuration file could specify default values, but anything baked
// in to the executable by setting this value before calling Load() will
// overwrite whatever has been specified in the cascade of configuration
// sources.
var BuildInfo string = ""

// The Opts type is used to
type Opts struct {
	// Name used to search for default config files with. If not specified
	// it will default ot os.Args[0] which might be "go" if you aren't
	// running an installed application. If you are using default config
	// files you probably want to specify this.
	Name string

	// A prefix used when looking for environment variables to parse as
	// value definitions. If not specified it will default to the name
	// of the program.
	EnvPrefix string

	// A default configuration at the lowest level of precedence. It's
	// more common to set these defaults using DefaultText, which is
	// applied after this.
	Defaults *AclNode

	// If this flag is set a default logging configuration identical to what
	// would be configured by a call to ColoredLoggingToConsole() will be
	// added onto any provided Defaults. This makes it easy in simple programs
	// that just need a couple configuraiton values and a reasonable console
	// output to be setup by a single call to Load().
	//
	// This configuration fragment is added on top of anything provided by
	// Defaults but before DefaultText is parsed so you can easily override
	// it from there or anywhere else as you wish.
	AddColorConsoleLogging bool

	// A valid ACL string which defines a default configuration. This
	// will be loaded on top of anything specified in Defaults. In general
	// you probably want to use one or the other. This really exists
	// as a shorthand to writing simple default configuration in a string
	// literal when you have a fairly simple app that only has a few
	// configuration values.
	DefaultText string


	// If set the default files based on the Name will not be loaded. This
	// flag can be set via the command line parsing using -i, but if it
	// is set in the Opts struct during load it can not be re-enabled via
	// the command line.
	IgnoreDefaultFiles bool

	// If set the environment will not be scanned for values.
	IgnoreEnvironment bool

	// If set os.Args will not be parsed
	IgnoreCommandLine bool

	// Additional files to load. The ExtrasRequired flag determines if
	// the files must exist.
	ExtraFiles []string

	// If this flag is set and one of the files specified by the ExtraFiles
	// slice is not found an error will be thrown.
	ExtrasRequired bool

	// Extra logging backends to attach to the logging infrastructure in
	// addition to ones that have been configured through the configuration
	// file. Each backend has a name so that it's level and formatting
	// can be configured via the same configuration files. Because backends
	// can effectively only be created once this allows an application to create
	// a unique local backed, such as a UI widget, that will receive all log
	// messages and still preserve the ability to configure other backends
	// via the configuration infrastructure.
	ExtraBackends map[string]logging.Backend

	// If set the config will be dumped to the log after parsing is done
	// the same as if the DUMPCONFIG_KEY was set at the root level. Useful
	// for when even basic parsing isn't working...
	DumpConfig bool
}

// Loads ArcherCL data from pontentially multiple locations and returns the root node.
// This is the main entry point of the package and it pulls together a standardized
// method of loading from a bunch of different potential sources. The defaults are
// reasonable so often you can just call this with a nil opts object.
//
// One of the most common options you will want to set is the Name of the
// application which identifies which default files should be looked for.
//
// Generally file system errors such as non-loadable default files are silently ignored
// because the most common cause is simple the files aren't being used. However, if a
// file does exist and an error was encountered during parsing, you probably want to
// know about that. Thus, any parsing errors will end the configuration process
// and be returned to the caller.
//
// Parsing errors will have a type of ParseLocation which provides further information
// about the exact error that was encountered.
//
// After everything else the last thing to be parsed into the config is a string
// from BuildInfo if set. See that global variable for more information.
//
// Once the configuration is setup, Load() will set the random number generator seed
// to a value from the key in the constant RANDOMSEED_KEY. See that constant for more.
//
// At the end of the configuration loading, the logging system from
// "github.com/op/go-logging" will be configured.  See the documentation for logging.go
// for example configuration values that can be used to setup all of the backends
// supported by that fairly robust package. Any additional logging backends, such as
// a native UI widget that wants to see all the log output, can be passed to the
// logging configuratino by setting the ExtraBackends variable.
//
// Often, a reasonable set of default logging options can be configured using the
// Default
func Load(opts *Opts) (*AclNode, error) {

	var err error

	// Get us a default options object
	if opts == nil {
		opts = &Opts{}
	}

	//fmt.Printf("Opts = %v", opts)

	programName := opts.Name
	if len(programName) == 0 {
		programName = os.Args[0]

		logDelayed(logging.DEBUG, "Program name = " + programName)
	}

	// Start with a base configuration passed in from the user if any
	cfg := opts.Defaults
	if cfg == nil {
		cfg = NewAclNode()
	}

	if opts.AddColorConsoleLogging {
		cfg.ParseString(COLOR_LOGGING_ACL, nil)
	}

	if len(opts.DefaultText) > 0 {
		location := &ParseLocation{
			Filename: "DefaultText",
		}
		err = cfg.ParseString(opts.DefaultText, location)
		if err != nil {
			return nil, err
		}
	}

	// Have to read the command line to see if we are going to ignore defaults or not
	ignoreDefaults := opts.IgnoreDefaultFiles
	filesToLoad := make([]string, 0)
	stringsToParse := make([]string, 0)

	// Parse the command line arguments
	if !opts.IgnoreCommandLine {
		clIgnore, clFilenames, clStrings := ParseCmdLine()

		// If options say ignore, then ignore, otherwise go with the command line
		if !ignoreDefaults {
			ignoreDefaults = clIgnore
		}

		filesToLoad = append(filesToLoad, clFilenames...)
		stringsToParse = append(stringsToParse, clStrings...)
	}


	// Unless we are ignoring them, load from the default values. If we get a
	// ParseLocation that means the file existed, and could be loaded, but was
	// whacky town. We want to let the caller know about that rather tha swallowing
	// these sorts of things.
	if !ignoreDefaults {
		// Start with the default files
		err = cfg.ParseFile("/etc/" + programName + ".acl")
		if pl, ok := err.(*ParseLocation); ok {
			return nil, pl
		}

		// In 1.12 this works, but not in 1.10
		// dir,err := os.UserHomeDir()

		// 1.10 style of getting the user's home dir
		current, err := user.Current()
		if err != nil {
			return nil, err
		}
		dir := current.HomeDir
		
		err = cfg.ParseFile(dir + "/." + programName + ".acl")
		if pl, ok := err.(*ParseLocation); ok {
			return nil, pl
		}

		err = cfg.ParseFile("./" + programName + ".acl")
		if pl, ok := err.(*ParseLocation); ok {
			return nil, pl
		}
	}

	// Load any files we found on the command like
	for _, fname := range filesToLoad {
		err = cfg.ParseFile(fname)
		if pl, ok := err.(*ParseLocation); ok {
			return nil, pl
		}
	}

	// Environment variables
	if !opts.IgnoreEnvironment {
		prefix := opts.EnvPrefix
		if prefix == "" {
			prefix = programName
		}

		env := make([]string, 0)
		for _, v := range os.Environ() {
			if strings.HasPrefix(v, prefix+"_") {
				env = append(env, v[len(prefix)+1:])
			}
		}
		cfg.ParseEnviron(env)
	}

	// Command line strings
	for ix, str := range stringsToParse {
		location := &ParseLocation{
			Filename: fmt.Sprintf("CMDLINE(%d)", ix),
		}
		cfg.ParseString(str, location)
	}

	// Possibly add some build info
	if len(BuildInfo) > 0 {
		bi := NewAclNode()
		_ = bi.ParseString(BuildInfo, nil) // TODO - handle the error
		//cfg.Children[BUILDINFO_KEY] = bi
		cfg.SetValAt(bi, BUILDINFO_KEY)
	}

	// Setup random either using a seed from the config or the time. This ensure
	// that we can both be testable or can have reasonale pseudo-randomness
	seed := int64(cfg.ChildAsInt(RANDOMSEED_KEY))
	if seed == 0 {
		seed = time.Now().UnixNano()
		cfg.SetValAt(seed, RANDOMSEED_KEY)
	}
	logDelayed(logging.DEBUG, fmt.Sprintf("Random seed is %d", seed))
	rand.Seed(seed)

	SetLoggingConfig(cfg)

	if cfg.ChildAsBool(DUMPCONFIG_KEY) || opts.DumpConfig {
		outputDelayedLog(alog)
		alog.Debug("Canonical config after all parsing:")
		if cfg.ChildAsBool(DUMPCOLOR_KEY) {
			alog.Debug(cfg.ColoredString())
		} else {
			alog.Debug(cfg.String())
		}
	}

	return cfg, nil
}

const (
	clStart = iota
	clConfigFilename
)

func ParseCmdLine() (ignore bool, filenames []string, toParse []string) {

	state := clStart

	ignore = false
	filenames = make([]string, 0)
	toParse = make([]string, 0)

	args := os.Args[1:]
	for _, arg := range args {
		switch state {
		case clStart:
			if len(arg) == 0 {
				continue
			}
			if arg[0] == '-' {
				arg = arg[1:]
				if arg[0] == '-' {
					arg = arg[1:]
				}

				if arg == "c" || arg == "config" {
					// The next arg is the filename
					state = clConfigFilename
					continue
				}

				if arg == "i" || arg == "ignore" {
					ignore = true
					continue
				}

				logDelayed(logging.ERROR, "Unrecognized command line argument '"+arg+"'")
			} else {
				// Store it as a string to parse after all files are loaded
				toParse = append(toParse, arg)
			}

		case clConfigFilename:
			if len(arg) == 0 {
				// Whatever
				continue
			}
			filenames = append(filenames, arg)
			state = clStart
		}
	}

	// fmt.Printf("Command line parsing ignore=%v  files=%v   toParse=%v", ignore, filenames, toParse)
	// panic("Stop")
	return
}

// An AclNode represents a node in a tree of configuration values that is
// rooted with a single AclNode. Each node either has Values or has Children,
// but never both. If len(Values) > 0 then Children should be ignored.
// Although the file format has syntactic sugar to specify multiple names
// for a node, each node only has a single name, which will be identical
// to the name it is stored under in it's parent. The root node has no name.
type AclNode struct {
	Values []interface{}

	Children map[string]*AclNode

	// These fields are related to the way things are parsed so that we can
	// reproduce the order they were written in
	OrderedChildNames []string
	IsMultiline       bool
	UsesEquals        bool
}

func NewAclNode() (node *AclNode) {
	return &AclNode{
		Values:   make([]interface{}, 0),
		Children: make(map[string]*AclNode),

		OrderedChildNames: make([]string, 0),
		IsMultiline:       false,
		UsesEquals:        false,
	}
}

// Attempts to load and parse the named file. If syntax errors occuring during
// the parsing, the error will be of type ParseLocation. Any other type is
// indicative of a issue loading the file.
func (node *AclNode) ParseFile(filename string) error {
	// fmt.Printf("ParseFile(%v)\n", filename)
	data, err := ioutil.ReadFile(filename)

	// fmt.Printf(" len(data)=%d  err=%v\n", len(data), err)

	if err != nil {
		logDelayed(logging.NOTICE, "Could not open file "+filename)
		return err
	}

	location := &ParseLocation{
		Filename: filename,
	}
	err = node.ParseString(string(data), location)
	if err != nil {
		logDelayed(logging.ERROR, err.Error())
	}

	return err
}

func (node *AclNode) ParseEnviron(env []string) {
	for _, e := range env {
		v := strings.SplitN(e, "=", 2)

		key := v[0]
		keys := strings.Split(key, "_")

		node.SetValAt(v[1], keys...)
	}
}

func (node *AclNode) ForEachOrderedChild(fn func(string, *AclNode)) {
	if node == nil || node.OrderedChildNames == nil {
		return
	}

	for i := 0; i < len(node.OrderedChildNames); i++ {
		name := node.OrderedChildNames[i]
		child := node.Children[name]

		fn(name, child)
	}
}

func (node *AclNode) ChildAsInt(names ...string) int {
	cNode := node.Child(names...)
	return cNode.AsInt()
}

func (node *AclNode) ChildAsFloat(names ...string) float64 {
	cNode := node.Child(names...)
	return cNode.AsFloat()
}

func (node *AclNode) ChildAsString(names ...string) string {
	cNode := node.Child(names...)
	return cNode.AsString()
}

func (node *AclNode) ChildAsBool(names ...string) bool {
	cNode := node.Child(names...)
	return cNode.AsBool()
}

/////
func (node *AclNode) ChildAsIntList(names ...string) []int {
	cNode := node.Child(names...)
	if cNode == nil {
		return []int{}
	}
	out := make([]int, len(cNode.Values))
	for ix, v := range cNode.Values {
		out[ix] = valAsInt(v)
	}
	return out
}

func (node *AclNode) ChildAsByteList(names ...string) []byte {
	cNode := node.Child(names...)
	if cNode == nil {
		return []byte{}
	}
	out := make([]byte, len(cNode.Values))
	for ix, v := range cNode.Values {
		out[ix] = byte(valAsInt(v))
	}
	return out
}

func (node *AclNode) ChildAsFloatList(names ...string) []float64 {
	cNode := node.Child(names...)
	if cNode == nil {
		return []float64{}
	}
	out := make([]float64, len(cNode.Values))
	for ix, v := range cNode.Values {
		out[ix] = valAsFloat(v)
	}
	return out
}

func (node *AclNode) ChildAsStringList(names ...string) []string {
	cNode := node.Child(names...)
	if cNode == nil {
		return []string{}
	}
	out := make([]string, len(cNode.Values))
	for ix, v := range cNode.Values {
		s := v.(string)
		out[ix] = s
	}
	return out
}

func (node *AclNode) ChildAsBoolList(names ...string) []bool {
	cNode := node.Child(names...)
	if cNode == nil {
		return []bool{}
	}
	out := make([]bool, len(cNode.Values))
	for ix, v := range cNode.Values {
		out[ix] = valAsBool(v)
	}
	return out
}

/////
func (node *AclNode) ChildAsBytes(names ...string) []byte {
	str := node.Child(names...).AsStringN(0)
	if len(str) == 0 {
		return []byte{}
	}

	data, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return []byte{}
	}

	return data
}

/////

func (node *AclNode) DefChildAsInt(def int, names ...string) int {
	cNode := node.Child(names...)
	out := cNode.AsInt()
	if out == 0 {
		return def
	}
	return out
}

func (node *AclNode) DefChildAsFloat(def float64, names ...string) float64 {
	cNode := node.Child(names...)
	out := cNode.AsFloat()
	if out == 0.0 {
		return def
	}
	return out
}

func (node *AclNode) DefChildAsString(def string, names ...string) string {
	cNode := node.Child(names...)
	out := cNode.AsString()
	if len(out) == 0 {
		return def
	}
	return out
}

func (node *AclNode) Child(names ...string) *AclNode {
	if node == nil || node.Children == nil {
		return nil
	}

	cNode := node
	for _, name := range names {
		if cNode != nil {
			if len(cNode.Values) > 0 || cNode.Children == nil {
				cNode = nil
			} else {
				cNode = cNode.Children[name]
			}
		}
	}

	return cNode
}

func (node *AclNode) Len() int {
	if node == nil || node.Values == nil {
		return 0
	}

	return len(node.Values)
}

func (node *AclNode) IsArray() bool {
	return node.Len() > 1
}

// By default return the LAST value as opposed to the first. We used to return the
// first but it seems more natural to use an "overwrite" behavior while preserving
// the option for people to get the first if they really want it.

// We have to protect the "last" section from nils
func (node *AclNode) AsInt() int {
	return node.AsIntN(-1)
}

func (node *AclNode) AsFloat() float64 {
	return node.AsFloatN(-1)
}

func (node *AclNode) AsString() string {
	return node.AsStringN(-1)
}

func (node *AclNode) AsBool() bool {
	return node.AsBoolN(-1)
}

// Get the first vaues
func (node *AclNode) FirstAsInt() int {
	return node.AsIntN(0)
}

func (node *AclNode) FirstAsFloat() float64 {
	return node.AsFloatN(0)
}

func (node *AclNode) FirstAsString() string {
	return node.AsStringN(0)
}

func (node *AclNode) FirstAsBool() bool {
	return node.AsBoolN(0)
}

// Get any on the values
func (node *AclNode) findValIx(ix int) int {
	if node == nil || node.Values == nil || len(node.Values) < 1 {
		// fmt.Printf("---- first\n")
		return -1
	}

	if ix < 0 {
		// Reverse it
		ix = len(node.Values) + ix
	}

	if ix < 0 || ix >= len(node.Values) {
		// fmt.Printf("---- second ix=%v\n", ix)
		return -1
	}

	return ix
}

func valAsInt(v interface{}) int {
	r, ok := v.(int64)
	if !ok {
		r, ok := v.(int32)
		if !ok {
			r, ok := v.(int)
			if !ok {
				return 0
			}
			return r
		}
		return int(r)
	}
	return int(r)
}

func (node *AclNode) AsIntN(ix int) int {
	ix = node.findValIx(ix)
	if ix < 0 {
		return 0
	}

	return valAsInt(node.Values[ix])
}

func valAsFloat(v interface{}) float64 {
	r, ok := v.(float64)
	if !ok {
		r, ok := v.(float32)
		if !ok {
			// BUT, it might be an int which we could cast as a float?
			// So use the result from valAsInt to try to get a numerical
			// value out of this thing
			return float64(valAsInt(v))
		}
		return float64(r)
	}
	return r
}

func (node *AclNode) AsFloatN(ix int) float64 {
	ix = node.findValIx(ix)
	if ix < 0 {
		return 0.0
	}

	return valAsFloat(node.Values[ix])
}

func valAsString(v interface{}) string {
	r, ok := v.(string)
	if !ok {
		r, ok := v.(fmt.Stringer)
		if !ok {
			// If we are here then it's almost certainly a literal so
			// let's just defer to how fmt.Sprintf thinks it should go
			return fmt.Sprintf("%v", v)
		}
		return r.String()
	}
	return r
}

func (node *AclNode) AsStringN(ix int) string {
	ix = node.findValIx(ix)
	if ix < 0 {
		// fmt.Printf("-- Default string for ix=%v node=%v\n", ix, node)
		// if node != nil {
		// 	if node.Values != nil {
		// 		fmt.Printf("-- len=%v\n", len(node.Values))
		// 	} else {
		// 		fmt.Printf("-- node.Values is nil\n")
		// 	}
		// } else {
		// 	fmt.Printf("-- nod is nil\n")
		// }
		return ""
	}
	// fmt.Printf("AsStringN(%v) v=%v\n", ix, node.Values[ix])

	return valAsString(node.Values[ix])
}

func valAsBool(v interface{}) bool {
	r, ok := v.(bool)
	if !ok {
		r, ok := v.(string)
		if !ok {
			r, ok := v.(fmt.Stringer)
			if !ok {
				return false
			}
			b, err := strconv.ParseBool(r.String())
			if err != nil {
				return false
			}
			return b
		}
		b, err := strconv.ParseBool(r)
		if err != nil {
			return false
		}
		return b
	}
	return r
}

func (node *AclNode) AsBoolN(ix int) bool {
	ix = node.findValIx(ix)
	if ix < 0 {
		return false
	}

	return valAsBool(node.Values[ix])
}

//////////////////////////////////////////////////////

func (node *AclNode) createChild(names ...string) (*AclNode, error) {
	if node == nil {
		return nil, fmt.Errorf("Can not create a child AclNode on nil")
	}

	if node.Children == nil {
		return nil, fmt.Errorf("AclNode Children map was nil")
	}

	cNode := node
	var nextNode *AclNode
	for ix, name := range names {
		if len(cNode.Values) > 0 {
			return nil, fmt.Errorf("AclNode %v had values", names[0:ix])
		}

		nextNode = cNode.Children[name]
		if nextNode == nil {
			nextNode = NewAclNode()
			cNode.Children[name] = nextNode
			cNode.OrderedChildNames = append(cNode.OrderedChildNames, name)
		}

		cNode = nextNode
	}

	return cNode, nil
}

func (node *AclNode) SetValAt(v interface{}, names ...string) error {
	target, err := node.createChild(names...)
	if err != nil {
		return err
	}

	// Clear it first if necessary
	if len(target.Values) > 0 {
		target.Values = target.Values[:0]
	}
	target.Values = append(target.Values, v)
	return nil
}

// TODO: This REALLY needs to be optimized, but for now if it will work...
func (node *AclNode) Duplicate() *AclNode {
	next := NewAclNode()

	if node == nil {
		// Always give you something....
		return next
	}

	str := node.String()
	next.ParseString(str, nil)
	// ignoring all sorts of error conditions...
	return next
}

//////////////////////////////////////////////////////

func (node *AclNode) valueTo(writer *bufio.Writer, indentStr string, level int, withColor bool, value interface{}) {

	switch v := value.(type) {
	case *AclNode:
		//alog.Info("Value is AclNode")
		v.StringTo(writer, indentStr, level, withColor)

	case string:
		if withColor {
			writer.WriteString(ansi.Black)
		}
		writer.WriteString(strconv.Quote(v))

	case int:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.Itoa(v))

	case int32:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.Itoa(int(v)))

	case int64:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.Itoa(int(v)))

	case float32:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 32))

	case float64:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.FormatFloat(v, 'g', -1, 64))

	case bool:
		if withColor {
			writer.WriteString(ansi.Red)
		}
		writer.WriteString(strconv.FormatBool(v))
	}
}

func (node *AclNode) maybeWriteNewline(writer *bufio.Writer, indentStr string, level int) bool {
	if node.IsMultiline {
		writer.WriteString("\n")
		for l := 0; l < level; l++ {
			writer.WriteString(indentStr)
		}
		return true
	}

	return false
}

// StringTo recursively prints the value of this node and any children that
// it has. When child nodes only have a single child it will collapse
// those nodes. (maybe??)
func (node *AclNode) StringTo(writer *bufio.Writer, indentStr string, level int, withColor bool) {

	//alog.Info("StringTo len(Values)=%d, len(Children)=%d", len(node.Values), len(node.Children))
	if len(node.Values) > 0 {
		// Print the values ignoring the children
		if len(node.Values) == 1 {
			node.valueTo(writer, indentStr, level, withColor, node.Values[0])
		} else {
			if withColor {
				writer.WriteString(ansi.Cyan)
			}
			writer.WriteString("[")
			if !node.maybeWriteNewline(writer, indentStr, level+1) {
				writer.WriteString(" ")
			}

			lastIx := len(node.Values) - 1
			for ix, value := range node.Values {
				isLast := ix == lastIx

				node.valueTo(writer, indentStr, level+1, withColor, value)
				if withColor {
					writer.WriteString(ansi.Cyan)
				}
				iLevel := level + 1
				if isLast {
					iLevel--
				} else {
					writer.WriteString(",")
				}
				if !node.maybeWriteNewline(writer, indentStr, iLevel) {
					writer.WriteString(" ")
				}
			}
			if withColor {
				writer.WriteString(ansi.Cyan)
			}
			writer.WriteString("]")
		}
	} else {
		// It is a map node with children instead of values
		if withColor {
			writer.WriteString(ansi.Magenta)
		}
		writer.WriteString("{")
		if !node.maybeWriteNewline(writer, indentStr, level+1) {
			writer.WriteString(" ")
		}

		// // We go to the extra trouble to sort the keys here so that
		// // the output is predictable, which aids in testing.
		// keys := make([]string, 0, len(node.Children))
		// for name, _ := range node.Children {
		// 	keys = append(keys, name)
		// }
		// sort.Strings(keys)

		// Instead of sorting (which we maybe want to make a flag) we will
		// use the natural order that the keys were originally in
		keys := node.OrderedChildNames

		last := len(keys) - 1
		for ix, name := range keys {
			isLast := ix == last

			if withColor {
				writer.WriteString(ansi.Black)
			}
			obj := node.Children[name]

			if withColor {
				writer.WriteString(ansi.Blue)
			}
			writer.WriteString(strconv.Quote(name))

			if withColor {
				writer.WriteString(ansi.Magenta)
			}
			if obj.UsesEquals {
				writer.WriteString(" = ")
			} else {
				writer.WriteString(": ")
			}
			if withColor {
				writer.WriteString(ansi.Black)
			}
			node.valueTo(writer, indentStr, level+1, withColor, obj)

			if withColor {
				writer.WriteString(ansi.Cyan)
			}

			iLevel := level + 1
			if isLast {
				iLevel--
			} else {
				if withColor {
					writer.WriteString(ansi.Magenta)
				}
				writer.WriteString(",")
			}

			if !node.maybeWriteNewline(writer, indentStr, iLevel) {
				writer.WriteString(" ")
			}
		}
		if withColor {
			writer.WriteString(ansi.Magenta)
		}
		writer.WriteString("}")
	}

	if withColor {
		writer.WriteString(ansi.Reset)
	}
	writer.Flush()
}

func (node *AclNode) String() string {
	var buf bytes.Buffer

	node.StringTo(bufio.NewWriter(&buf), "\t", 0, false)

	return buf.String()
}

func (node *AclNode) ColoredString() string {
	var buf bytes.Buffer

	node.StringTo(bufio.NewWriter(&buf), "\t", 0, true)

	return buf.String()
}

func (node *AclNode) PrettyVersion() string {
	bi := node.Child(BUILDINFO_KEY)

	if bi.ChildAsBool("travis", "server") {
		return fmt.Sprintf("T#%s (%d)%s %s",
			bi.ChildAsString("travis", "build", "number"),
			bi.ChildAsInt("git", "commits"),
			bi.ChildAsString("git", "short"),
			bi.ChildAsString("time", "local"))

	} else {
		return fmt.Sprintf("**LOCAL BUILD** (%d)%s %s",
			bi.ChildAsInt("git", "commits"),
			bi.ChildAsString("git", "short"),
			bi.ChildAsString("time", "local"))
	}
}

type ParseLocation struct {
	Filename string
	Line     int
	Col      int

	Message string
}

func (l *ParseLocation) Error() string {
	if l == nil {
		return ""
	}

	return fmt.Sprintf("%v:%d:%d: %v", l.Filename, l.Line, l.Col, l.Message)
}

// StringToACL is a function for immediately parsing simple configs into
// an AclNode tree. It is useful for test case writing, but probably should not
// be used for real code because it hides errors.
func StringToACL(str string) *AclNode {
	node := NewAclNode()
	_ = node.ParseString(str, nil)
	return node
}
