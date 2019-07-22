# Archer Configuration Langauge

The Archer Configuration Language is compatible with "JSON plus comments" 
and is meant to be [UCL](https://github.com/vstakhov/libucl) -ish more or less.
The major benefit over JSON is allowing for comments, a relaxed structure,
and some syntactic sugar here and there. The advantage of ACL over 
[HCL](https://github.com/hashicorp/hcl) from [HashiCorp](https://www.hashicorp.com)
is that HCL isn't actually documented well, and it uses a model where in your
go code you would define a go data structure, and then would parse a 
single file into that structure.

In contrast to HCL, ACL is based on the concept that you parse one or more
strings (which are probably loaded from files), into a configuration object
(a single `AclNode` root object) and then in your code you query for known
values from this common root. We could implement the HCL style "fill out this
data structure" functionality using an ACL root node, but have not yet 
done that as of yet. The main advantage of the pull approach implemented
by ACL is that it allows loading multiple files into one structure and then
consumption by different modules which don't know about each others 
configuration needs.

In ACL the basic expression is a key/value tuple. Keys are identifiers. They can be 
unqouted, in which they must start with an alpha or underscore character
and can only contain alphanumeric and underscore characters. If they are
quoted with single or double quotes they can contain anything. (There is
a case where keys may begin with a `!` character for a special meaning
which is discussed in more detail below.)

Values are written as ints, floats, strings, or hexadecimal ints. Booleans
are written as strings that can be interpretted by 
[strconv.ParseBool()](https://golang.org/pkg/strconv/#ParseBool),
which means there is a wide range of 1, t, T, TRUE, true, True type expressions 
that are valid.

Keys are separated from values by either `:` or `=` (this separator is optional
for object values, see below).

Comments can be C, C++, bash, or sql style

**A simple example**

	/* 
		Multi-line C comment
	*/
	port: 1234
	seed: 0xfeedface

	// A C++ style comment
	tcpEndpoint: 'foo.com:80'  // Also at the end of a line

	# Bash style
	name = "Bob"

	message = "Hello there bob" -- sql style comments for funsies
	"reply to" = me

Quotes are optional for alphanumeric string values, but mandatory if the 
values contain spaces or punctuation. Hence the need for quotes on the
`tcpEndpoint` above.

In addition to the simple key/value syntax, any value may be an array where
values are separated by whitespace and terminated by either a newline or
a `;` character. Enclosing `[` and `]` braces are optional, as
are separating commas within the braces. If arrays are going to span
more than one line, they must be enclosed in `[]`'s, otherwise the array 
terminates at the fist newline.

	an_array = one two three
	another_one = [80,81,82,83]
	a_third = [ "fred" "wilma" "barney" ]

As with UCL, arrays can be specified simply by specifying the same value
multiple times as follows

	# This 

	ports = 80
	ports = 90
	ports = 100

	# is equivalent to

	ports = 80 90 100

	# is equivalent to

	ports = [80, 90, 100];

*Note:* Because of this additive nature, if you parse the configuration
above you will get the following canonical result

	{
		"ports": [
			80,
			90,
			100,
			80,
			90,
			100,
			80,
			90,
			100,
		],
	}

This additive cascade is the default to keep the language inline with
the way UCL works. If this is not what you want, the key name can be
prefixed by a `!` to indicate that any previous values should be overwritten
starting with that instance of the key.

	# Build an array
	key = a
	key = b
	key = c

	# Right now, key = [a,b,c]

	# Then overwrite it
	!key = d
	key = e

	# Now, key = [d, e]

In practice this functionality is probably most useful when values are 
expected to be defined in multiple files. The default additive cascade is
necessary so that a cascaded file can be allowed to define a small portion
of an object without needing to redefine the other values. This reset features
allows the additive cascade to be the default while still allowing a later file 
to simply redefine something that was previously defined in an earlier file.

## Objects

The third form of a value is an object which itself contains keys and values.
As with JSON, objects are identified by wrapping a set of key/value definitions in 
enclosing curly `{` and `}` brackets. There is an implicit object at the top
level. Expressing a second unnamed object immediately inside a previous one is ignored,
thus the following three examples evaluate to the same thing:

**Implicit Only**

	one = 1
	two = 1

**Top Level Object**

	{
		one = 1
		two = 2
	}

**Duplicate Object**

	{
		{
			one = 1
		}
		two = 2
	}

Objects are values, with or without the `:` or `=` character between the name
and the definition. They can of course be embedded anywhere a value can appear,
including inside arrays.

	redis {
		server: "cloud.server.com"
		password: "abcdefgh"
		prefix: "cyril-production"
	}

	endpoints = [
		{ host: "foo.com", port: 80 }
		{ host: "bar.com", port: 800 }
	]

Objects are used to namespace key names as with the examples above. 

Multi-level objects can also be specified by simply listing key names as arrays
of identifiers, albeit without the `[]` braces and without commas. Quoting is fine.
This makes the most sense when one of the names is a type name, typically the
earlier name, followed by an instance name.

	server cyril {
		port : 9771
		hostname : "home.decidedly.com"
	}

	server ray port = 9772
	server ray "instance count" = 4

	server sterling port = 9772

Canonically the above is expressed as 

	{
		"server": {
			"cyril": {
				"hostname": "home.decidedly.com",
				"port": 9771,
			},
			"ray": {
				"instance count": 4,
				"port": 9772,
			},
			"sterling": {
				"port": 9772,
			},
		},
	}

## API

The base object of the API is the `AclNode` struct. The configuration file(s) is 
parsed into a tree of `AclNode`'s with a single root. Each `AclNode` either has values or
children, but not both. While the data structure technically allows both, nodes
produced by the parser will never contain both, and if such a structure is
encountered, the `Children` map will be ignored.

The common / easy / **what you should do** way of accessing values in the configuration
is to use one of the `ChildAsXXXX()` functions.

  * `node.ChildAsInt(names ...string) int`
  * `node.ChildAsFloat(names ...string) float64`
  * `node.ChildAsString(names ...string) string`
  * `node.ChildAsBool(names ...string) bool`

For instance with the last example configuration:

	hostname := cfg.ChildAsString("server", "cyril", "hostname")
	// = "home.decidedly.com"

	count := cfg.ChildAsint("server", "ray", "port")
	// = 4

A limited amount of type coercision will be performed, primarily attempting to convert
to and from `string` and the requested type. The `AsBool` methods in particular rely on
the use of [strconv.ParseBool()](https://golang.org/pkg/strconv/#ParseBool) to turn 
strings into booleans. The ACL file format does not include an explicit boolean 
syntax, but the API does include `AsBool` methods that parse strings because this is
such a common occurrence.

While those are the most common methods for simple use of the API, `AclNode` values
also have the following for more verbose usage. All methods operate safely on `nil` or
uninitialized / zero value `AclNode` pointers.

  * `node.Child(names ...string) *AclNode` - walks down the tree and finds a specific 
    child or returns `nil` if the child is not found. Used as the basis for the `ChildAsXXX()`
    functions.
  * `node.Len() int` - convenient and safe method for accessing `len(node.Values)`
  * `node.AsInt() int` - convenience for `AsIntN(0)`
  * `node.AsFloat() float64` - convenience for `AsFloatN(0)`
  * `node.AsString() string` - convenience for `AsStringN(0)`
  * `node.AsBool() bool` - convenience for `AsBoolN(0)`
  * `node.AsIntN(ix int) int` - attempt to coerce `node.Values[ix]` into an `int`
    returning the zero value if non-existent or if coercision fails
  * `node.AsFloatN(ix int) float64` - attempt to coerce `node.Values[ix]` into a `float64`
    returning the zero value if non-existent or if coercision fails
  * `node.AsStringN(ix int) float64` - attempt to coerce `node.Values[ix]` into a `string`
    returning the zero value if non-existent or if coercision fails
  * `node.AsBoolN(ix int) float64` - attempt to coerce `node.Values[ix]` into a `bool` 
    returning the zero value if non-existent or if coercision fails
  * `node.StringTo(writer *bufio.Writer, indentStr string, level int)` - Writes the value
    of the node to a `bufio.Writer` with the given indention level and indention string.
    This method recurses into AclNode's it encounters and is the mechanism by which the
    `String()` method creates as string.
  * `node.String() string` - Implementation of the Stringer interface. Calls `StringTo()`
    with an indention string of `\t` and a level of `0`. 

The `node.String()` method can be interpretted as producing a canonical form of the
configuration structure contained in memory. This structure includes ignoring the
children of nodes which have both `Values` and `Children`. It doesn't produce redundant
object definitions, and it does always explicitly list each object level with a
single key name instead of using arrays as key names. It will explictly include the
root level curly braces.

Thus, while canonical, the output produced by `node.String()` isn't really what you
want to use when writing a configuration file meant to be easily interpretted by
humans. Not only will it not contain any comments, the syntatic sugar of the language
is meant to make things more concise and readable. The `node.String()` method is 
meant for debugging and for test cases - which is why it alphabetizes the results so
that they stay the same from test run to test run.


## Test Files

There is an `acl_test.go` file which will run a non-exhaustive set of tests, primarily
oriented towards the parser. Because of the nature of trying to test configuration
files that can become complex quickly, there is a test utility meant for intepretation
by a human in `/sandbox/acl_t.go` 

The sandbox test has an array of here documents that are run through the parser and 
then dumped in canonical form. If you need to work on the parser, it's probably
going to be really useful for debugging / fixing / implementing more complex
syntax that might not have gotten tested otherwise.


## Building

To be able to run `go generate` in this module you need both `ragel` and `graphviz`
installed to build the parser. 

	brew install ragel graphviz

The `go generate` commands are at the top of acl.go. They turn `acl_parser.rl` into
`acl_parser.go` and also `acl_parser.png` which is a diagram of the state machine
the parser uses to interpret the input file. It's a little difficult to disect because
it lists characters by ascii value, but at least it's there if you wish to consult it.

The parse is actually a combination of a state machine created and executed by the
ragel code, with a couple of state flags and a stack of names that are maintained
by the action code. This dual nature might be a little confusing to work on, but 
seemed like the clearest implementation as opposed to embedding absolutely
everything into the ragel machine.

More information about ragel can be found at 
[The Ragel Website](http://www.colm.net/open-source/ragel/)

