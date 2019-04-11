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
	//"sort"
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
			return ""
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
// be used for real code because it hides errors.
func StringToACL(str string) *AclNode {
	node := NewAclNode()
	_ = node.ParseString(str, nil)
	return node
}
