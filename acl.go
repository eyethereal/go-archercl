//go:generate ragel -Z -o acl_parser.go acl_parser.rl
//go:generate ragel -V -o acl_parser.dot acl_parser.rl
//go:generate dot -oacl_parser.png -Tpng acl_parser.dot

package config

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
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	NODE_BUILDINFO = "buildInfo"

	DEFAULT_BUILDINFOACL = NODE_BUILDINFO + ": \"Unspecified\""
)

var alog = Logger("acl")

// This variable will be set by the init() funciton in a build.go file that
// is written into the config package during building on the CI server. It
// is parsed at the end of configuration and is manually inserted into the
// resulting structure so that it won't/can't be overwritten by accident.
var BuildInfo string = ""

// LoadACLConfig is the standard method of loading and parsing a
// configuration.
func LoadACLConfig(name string, prefix string) *AclNode {

	cfg := NewAclNode()

	if name == "" {
		name = "archer"
	}

	// Parse the command line arguments
	ignoreDefaults, filenames, toParse := ParseCmdLine()

	if !ignoreDefaults {
		// Start with the default files
		_ = cfg.ParseFile("/etc/" + name + ".acl")

		usr, _ := user.Current()
		dir := usr.HomeDir
		_ = cfg.ParseFile(dir + "/." + name + ".acl")

		cfg.ParseFile("./" + name + ".acl")
	}

	for _, fname := range filenames {
		_ = cfg.ParseFile(fname)
	}

	// Environment variables
	if prefix == "" {
		prefix = name
	}
	prefix = strings.ToUpper(name)

	env := make([]string, 0)
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, prefix+"_") {
			env = append(env, v[7:])
		}
	}
	cfg.ParseEnviron(env)

	// Command line strings
	for ix, str := range toParse {
		location := &ParseLocation{
			Filename: fmt.Sprintf("CMDLINE(%d)", ix),
		}
		cfg.ParseString(str, location)
	}

	{
		// Manually add the build info node
		bi := NewAclNode()
		_ = bi.ParseString(BuildInfo, nil) // TODO - handle the error
		cfg.Children[NODE_BUILDINFO] = bi
	}

	// Setup random either using a seed from the config or the time. This ensure
	// that we can both be testable or can have reasonale pseudo-randomness
	seed := int64(cfg.ChildAsInt("randomSeed"))
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	logDelayed(logging.DEBUG, fmt.Sprintf("Random seed is %d", seed))
	rand.Seed(seed)

	SetLoggingConfig(cfg)

	if cfg.ChildAsBool("dumpConfig") {
		outputDelayedLog(alog)
		alog.Debug("Canonical config after all parsing:")
		if cfg.ChildAsBool("dumpColor") {
			alog.Debug(cfg.ColoredString())
		} else {
			alog.Debug(cfg.String())
		}
	}

	return cfg
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
}

func NewAclNode() (node *AclNode) {
	return &AclNode{
		Values:   make([]interface{}, 0),
		Children: make(map[string]*AclNode),
	}
}

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

		key := strings.ToLower(v[0])
		keys := strings.Split(key, "_")

		cNode := node
		for _, k := range keys {
			cNode := node.Children[k]
			if cNode == nil {
				cNode = NewAclNode()
				node.Children[k] = cNode
			}
		}

		cNode.Values = append(cNode.Values, v[1])
	}
}

func (node *AclNode) ChildAsInt(names ...string) int {
	cNode := node.Child(names...)
	return cNode.AsIntN(0)
}

func (node *AclNode) ChildAsFloat(names ...string) float64 {
	cNode := node.Child(names...)
	return cNode.AsFloatN(0)
}

func (node *AclNode) ChildAsString(names ...string) string {
	cNode := node.Child(names...)
	return cNode.AsStringN(0)
}

func (node *AclNode) ChildAsBool(names ...string) bool {
	cNode := node.Child(names...)
	return cNode.AsBoolN(0)
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

func (node *AclNode) DefChildAsInt(def int, names ...string) int {
	cNode := node.Child(names...)
	out := cNode.AsIntN(0)
	if out == 0 {
		return def
	}
	return out
}

func (node *AclNode) DefChildAsFloat(def float64, names ...string) float64 {
	cNode := node.Child(names...)
	out := cNode.AsFloatN(0)
	if out == 0.0 {
		return def
	}
	return out
}

func (node *AclNode) DefChildAsString(def string, names ...string) string {
	cNode := node.Child(names...)
	out := cNode.AsStringN(0)
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

func (node *AclNode) AsInt() int {
	return node.AsIntN(0)
}

func (node *AclNode) AsFloat() float64 {
	return node.AsFloatN(0)
}

func (node *AclNode) AsString() string {
	return node.AsStringN(0)
}

func (node *AclNode) AsBool() bool {
	return node.AsBoolN(0)
}

func (node *AclNode) AsIntN(ix int) int {
	if node == nil || node.Values == nil || len(node.Values) < 1 {
		return 0
	}

	if ix >= len(node.Values) {
		return 0
	}

	r, ok := node.Values[ix].(int64)
	if !ok {
		r, ok := node.Values[ix].(int32)
		if !ok {
			r, ok := node.Values[ix].(int)
			if !ok {
				return 0
			}
			return r
		}
		return int(r)
	}
	return int(r)
}

func (node *AclNode) AsFloatN(ix int) float64 {
	if node == nil || node.Values == nil || len(node.Values) < 1 {
		return 0
	}

	if ix >= len(node.Values) {
		return 0.0
	}

	r, ok := node.Values[ix].(float64)
	if !ok {
		r, ok := node.Values[ix].(float32)
		if !ok {
			return 0.0
		}
		return float64(r)
	}
	return r
}

func (node *AclNode) AsStringN(ix int) string {
	if node == nil || node.Values == nil || len(node.Values) < 1 {
		return ""
	}

	if ix >= len(node.Values) {
		return ""
	}

	r, ok := node.Values[ix].(string)
	if !ok {
		r, ok := node.Values[ix].(fmt.Stringer)
		if !ok {
			return ""
		}
		return r.String()
	}
	return r
}

func (node *AclNode) AsBoolN(ix int) bool {
	if node == nil || node.Values == nil || len(node.Values) < 1 {
		return false
	}

	if ix >= len(node.Values) {
		return false
	}

	r, ok := node.Values[ix].(bool)
	if !ok {
		r, ok := node.Values[ix].(string)
		if !ok {
			r, ok := node.Values[ix].(fmt.Stringer)
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
			writer.WriteString("[\n")
			for _, value := range node.Values {
				if withColor {
					writer.WriteString(ansi.Black)
				}
				for l := 0; l <= level; l++ {
					writer.WriteString(indentStr)
				}
				node.valueTo(writer, indentStr, level+1, withColor, value)
				if withColor {
					writer.WriteString(ansi.Cyan)
				}
				writer.WriteString(",\n")
			}
			for l := 0; l < level; l++ {
				writer.WriteString(indentStr)
			}
			writer.WriteString("]")
		}
	} else {
		// It is a map node with children instead of values
		if withColor {
			writer.WriteString(ansi.Magenta)
		}
		writer.WriteString("{\n")

		// We go to the extra trouble to sort the keys here so that
		// the output is predictable, which aids in testing.
		keys := make([]string, 0, len(node.Children))
		for name, _ := range node.Children {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		for _, name := range keys {
			if withColor {
				writer.WriteString(ansi.Black)
			}
			obj := node.Children[name]
			for l := 0; l <= level; l++ {
				writer.WriteString(indentStr)
			}

			if withColor {
				writer.WriteString(ansi.Blue)
			}
			writer.WriteString(strconv.Quote(name))

			if withColor {
				writer.WriteString(ansi.Cyan)
			}
			writer.WriteString(": ")
			if withColor {
				writer.WriteString(ansi.Black)
			}
			node.valueTo(writer, indentStr, level+1, withColor, obj)

			if withColor {
				writer.WriteString(ansi.Cyan)
			}
			writer.WriteString(",\n")
		}
		if withColor {
			writer.WriteString(ansi.Black)
		}
		for l := 0; l < level; l++ {
			writer.WriteString(indentStr)
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
	bi := node.Child(NODE_BUILDINFO)

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
// be used for real code
func StringToACL(str string) *AclNode {
	node := NewAclNode()
	_ = node.ParseString(str, nil)
	return node
}
