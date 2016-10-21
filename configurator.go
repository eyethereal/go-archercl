// /*
//    The config package provides a way for us to load process configuration data from a cascade
//    of files and finally the command line. For now, it is a flat key-value structure, but the
//    thought is that it should support hierarchical data at some future point.

//    The Group is the basic node structure. Right now, there is a single root Group which contains
//    all loaded values. Trying to get sub groups will always return nil, mostly because the loading
//    routines don't currently know how to put data they find into sub-groups.

//    Given a filename FILENAME the files that will be loaded in order are:

//        /etc/FILENAME
//        ~/.FILENAME
//        ./FILENAME

//    Missing files are silently ignored. Values in later files overwrite values established in earlier
//    files.

//    The syntax of the files is very simple. Comment lines begin with a # as the first non whitespace
//    character on a line and take the entire line. There is no support for inline comments
//    coming after value definitions.

//    Names consist of whitespace trimmed characters to the left of an = sign and values are the
//    whitespace trimmed portion of the line to the right of the = sign. Lines without an = sign are
//    ignored, although they will probably print a warning when read.

//    An example config file might be something like this

//        hostname = fred.eyethereal.com
//        port=1234

//        # This is a good option to set
//        useSSL=true

//    After all 3 possible files are loaded, any options specified on the command line are added to
//    the configuration. Option names are specified with an initial - followed immediately by the
//    name of the option. The next argument is taken as the value for that name. Arguments which
//    do not immediately follow a name, as indicated by a -, are ignored.

//    Building on the previous example, one might do something like

//        go run server.go -port $PORT -useSSL false

//    All values are stored internally as strings and converted using the standard strconv functions
//    to other value types when read. There is no quoting or support for special formatting outside
//    of what will be handled by strconv natively.

//    The reading functions have two forms. One which will return the zero value for the type and
//    one which will return a specified default value of an error is returned by strconv during
//    conversion.

// */
package config

// import (
// 	"bufio"
// 	"encoding/base64"
// 	"fmt"
// 	"github.com/eyethereal/archer/security"
// 	"github.com/op/go-logging"
// 	"log"
// 	"os"
// 	"os/user"
// 	"strconv"
// 	"strings"
// )

// const _DEFAULT_LISTCAPACITY int = 4

// // A Group is the basic tree node for the configuration file. Each group may
// // have both direct named values as well as sub-groups. The namespaces for
// // values and sub-groups do not overlap. That is, a group may have both a
// // "foo" value and a "foo" sub-group.
// type Group struct {
// 	vals   map[string]string
// 	groups map[string]*Group
// }

// // SubGroup returns a Group which is a sub group of the target with the given
// // name. Use this to walk down the tree into more and more specific branches.
// func (g *Group) SubGroup(name string) *Group {
// 	return g.groups[name]
// }

// // Int parses the value stored under the specified name as an integer.
// // If the conversion fails (i.e. if the value is empty) it will return 0.
// func (g *Group) Int(name string) int {
// 	v, _ := strconv.ParseInt(g.vals[name], 0, 0)
// 	return int(v)
// }

// // IntDef is the same as Int but if no value is found, the specified default
// // value from the second parameter is returned instead of 0.
// func (g *Group) IntDef(name string, def int) int {
// 	v, e := strconv.ParseInt(g.vals[name], 0, 0)
// 	if e != nil {
// 		return def
// 	}
// 	return int(v)
// }

// // Float32 parses the value stored under the specifed name as a Float32.
// // If the conversion fails (i.e. if the value is empty) it will return the
// // float zero value.
// func (g *Group) Float32(name string) float32 {
// 	v, _ := strconv.ParseFloat(g.vals[name], 32)
// 	return float32(v)
// }

// // Float32Def is the same as Float32 but if no value is found, the specified
// // default value from the second parameter is returned.
// func (g *Group) Float32Def(name string, def float32) float32 {
// 	v, e := strconv.ParseFloat(g.vals[name], 32)
// 	if e != nil {
// 		return def
// 	}
// 	return float32(v)
// }

// // Bool uses strconv.ParseBool to turn the stored value into a boolean. This
// // will return false if the value is not present or is not otherwise understood
// // as a bool. Because of the use of strconv, accepted values are
// // 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Other values
// // also return false.
// func (g *Group) Bool(name string) bool {
// 	v, _ := strconv.ParseBool(g.vals[name])
// 	return bool(v)
// }

// // BoolDef is like Bool, but true may be specified as the value to return if
// // the named value is not present or is not parsable.
// func (g *Group) BoolDef(name string, def bool) bool {
// 	v, e := strconv.ParseBool(g.vals[name])
// 	if e != nil {
// 		return def
// 	}
// 	return bool(v)
// }

// // String returns the named value directly or a 0 length string if the value
// // has not been specified in the loaded configuration.
// func (g *Group) String(name string) string {
// 	return g.vals[name]
// }

// // StringDef is like String except the value specified by the second paramter
// // is returned if nothing is found in the loaded configuration.
// func (g *Group) StringDef(name string, def string) string {
// 	v := g.vals[name]
// 	if len(v) == 0 {
// 		return def
// 	}
// 	return v
// }

// // StringList returns a list of strings that are stored in value keys which
// // all begin with the same prefix name. The specified prefix name is appended
// // with integers begining at 0 until a name is hit for which no value is
// // specified. Thus, to load 4 integers in an array they would be listed in
// // the config file as
// //
// //      prefix.0 = one
// //      prefix.1 = two
// //      prefix.2 = three
// //      prefix.3 = four
// //
// // That would return the slice
// //
// //      []string{"one", "two", "three", "four"}
// //
// func (g *Group) StringList(name string) []string {

// 	list := make([]string, 0, _DEFAULT_LISTCAPACITY)

// 	for i := 0; ; i++ {

// 		elementName := fmt.Sprintf("%v.%v", name, i)

// 		element := g.String(elementName)

// 		if len(element) == 0 {
// 			break
// 		}

// 		list = append(list, element)
// 	}

// 	return list
// }

// // Bytes returns a slice of bytes. This is done by finding a string and then decoding
// // it using the URLEncoding base64. Padding is optional (The ='s signs can be stripped for brevity)
// func (g *Group) Bytes(name string) []byte {
// 	s := g.vals[name]
// 	if len(s) == 0 {
// 		return []byte{}
// 	}

// 	s = security.RestoreBase64Padding(s)
// 	data, err := base64.URLEncoding.DecodeString(s)
// 	if err != nil {
// 		return []byte{}
// 	}

// 	return data
// }

// // Parses a single line from a config file, ignoring comments and so forth
// func parseLine(root *Group, line string) {
// 	line = strings.TrimSpace(line)

// 	if len(line) == 0 {
// 		return
// 	}

// 	if line[0] == '#' {
// 		// A comment
// 		// log.Println("Ignoring comment line: ", line)
// 		return
// 	}

// 	ix := strings.Index(line, "=")
// 	if ix == -1 {
// 		log.Println("Config line had no = in it: ", line)
// 		return
// 	}

// 	lhs := strings.TrimSpace(line[:ix])
// 	rhs := strings.TrimSpace(line[ix+1:])

// 	// log.Println(lhs,"=", rhs)

// 	// Is lhs in a group or something???

// 	// TODO: Groups....

// 	root.vals[lhs] = rhs
// }

// // Parses a text file with the given name. If the name is not found this
// // function returns silently. The data found in the configuration file
// // overwrites any other data currently in the config. It won't delete
// // previously stored values, but previously stored values could be set to
// // the empty string or another zero value by overridding them in a conifg
// // file with a lower precedence.
// //
// // The results are added to the passed in Group, so that multiple config
// // files can be added one on top the other.
// func parseFile(root *Group, name string) {

// 	file, err := os.Open(name)
// 	if err != nil {
// 		logDelayed(logging.NOTICE, "Could not open file '"+name+"'")
// 		return
// 	}

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		parseLine(root, scanner.Text())
// 	}
// 	// Don't care about errors

// 	logDelayed(logging.INFO, "Parsed file '"+name+"'")
// }

// // Parses variables from the command line into the given Group
// func parseCmdLine(root *Group) {
// 	var key string

// 	for _, arg := range os.Args {

// 		if arg[0] == '-' {
// 			key = arg[1:]
// 		} else {
// 			if len(key) != 0 {
// 				root.vals[key] = arg

// 				logDelayed(logging.INFO, "Cmdline specified '"+key+"'")
// 				key = ""
// 			}
// 		}
// 	}
// }

// // Attempts to load 3 different variations of the configuration name given and
// // then adds any command line defined values. If the name is empty it will default
// // to `archer.cfg`
// func LoadConfig(name string) *Group {

// 	// First read system values
// 	root := &Group{
// 		make(map[string]string),
// 		make(map[string]*Group),
// 	}

// 	if len(name) == 0 {
// 		name = "archer.cfg"
// 	}

// 	parseFile(root, "/etc/"+name)

// 	usr, _ := user.Current()
// 	dir := usr.HomeDir

// 	parseFile(root, dir+"/."+name)
// 	parseFile(root, "./"+name)

// 	parseCmdLine(root)

// 	SetLoggingConfig(root)

// 	if root.Bool("dumpConfig") {
// 		lgr := logging.MustGetLogger("config")
// 		outputDelayedLog(lgr)
// 		lgr.Debug("Full parsed configuration:")
// 		for k, v := range root.vals {
// 			lgr.Debug("\t%s = %s", k, v)
// 		}
// 	}

// 	return root
// }

// // This is a convenience method for testing situations in particular where
// // you wish to create a Group easily from a map[string]string instance
// func MakeGroup(kvs map[string]string) *Group {
// 	return &Group{vals: kvs}
// }
