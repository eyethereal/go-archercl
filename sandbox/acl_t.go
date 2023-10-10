package main

import (
	"bytes"
	"fmt"
	"github.com/eyethereal/go-archercl"
	"github.com/mgutz/ansi"
	"log"
	// "strconv"
)

var samples = []string{
	`
port: 1234
seed: 0xfeedface
tcpEndpoint: 'foo.com:80'
name = "Bob"
message = "Hello there bob"
"reply to" = me
`,

	`
/* 
    Multi-line C comment
*/
port: 1234

// A C++ style comment
tcpEndpoint: 'foo.com:80'  // Also at the end of a line

# Bash style
name = "Bob"

message = "Hello there bob" -- sql style comments for funsies
"reply to" = me
`,

	`
an_array = one two three
another_one = [80,81,82,83]
a_third = [ "fred" "wilma" "barney" ]
`,

	`
# This 

ports = 80
ports = 90
ports = 100

# is equivalent to

ports = 80 90 100

# is equivalent to

ports = [80, 90, 100];

`,

	`
# Build an array
key = a
key = b
key = c

# Right now key = [a,b,c]

# Then overwrite it
!key = d
key = e

# key is now [d, e]
`,

	`
one = 1
two = 1
`,

	`
{
    one = 1
    two = 2
}
`,

	`
{
    {
        one = 1
    }
    two = 2
}
`,

	`
server cyril {
    port : 9771
    hostname : "home.decidedly.com"
}

server ray port = 9772
server ray "instance count" = 4

server sterling port = 9772
`,

	`
{
"server": {
    "cyril": {
        "hostname": "home.decidedly.com"
        "port": 9771
        }
    "ray": { "instance count" = 4, "port" = 9772, }
    "sterling": { "port" = 9772, }
    }
}
`,

	`
link: {
    tcp: {
        listen: {
            portMin: 9917
            portMax: 9917
        }
    }
}
`,
	`
{
    "link": {
        "tcp": {
            "listen": {
                "portMax": 9917,
                "portMin": 9917,
            },
        },
    },
}
`,

	`
{
    "home": {
        "buildInfo": "Unspecified",
    },
    "webapp": {
        "buildInfo": "Unspecified",
    },
    fred : 1
},
`,

	`
    endpoints = { host: "foo.com", port: 800}
    endpoints = { host: "bar.com", port: 900}    
`,

	`
    endpoints = [
        { host: "foo.com", port: 800}
        { host: "bar.com", port: 900}
    ]
`,

	`
paths = [
    [ one two three ]
    [ 1 2 3 ]
    [ [ a b ] c ]
]
`,

	`
a = [
    {
        names: [fred barney]
    }
    {
        names: [jim larry]
    }
    cows
    [dogs cats]
]
`,
}

func main() {

	for n, sample := range samples {
		// Run a single test...
		if n != 14 {
			//continue
		}

		fmt.Printf(ansi.Color("\n==========================================\n", "black+b"))
		fmt.Printf(ansi.Color("Test #%d:\n", "black+b"), n)
		fmt.Printf(ansi.Color("%s\n", "blue"), sample)
		node := archercl.NewAclNode()

		buffer := &bytes.Buffer{}
		logger := log.New(buffer, "", 0)
		//err := node.ParseString(sample, nil)
		err := node.ParseStringWithLogger(sample, nil, logger)
		//fmt.Printf("%s\n", buffer.String())

		if err != nil {
			fmt.Printf("%s: %s\n", ansi.Color("ERROR:", "red+b"), ansi.Color(err.Error(), "red"))

			continue
		}

		// That may have printed a bunch of debugging info, so maybe we want to redo the visual
		// comparison here...
		// fmt.Printf(ansi.Color("%s\n", "blue"), sample)

		firstString := node.String()
		fmt.Printf(ansi.Color("%s\n", "green"), node.ColoredString())
		// fmt.Printf(ansi.Color("%s\n", "magenta"), strconv.Quote(firstString))

		// Now re-parse the result to make sure our canonical representation is re-parsable
		node2 := archercl.NewAclNode()
		err = node2.ParseString(firstString, nil)
		if err != nil {
			fmt.Printf("%s: %s\n", ansi.Color("ERROR:", "red+b"), ansi.Color(err.Error(), "red"))
			continue
		}

		secondString := node2.String()
		if firstString != secondString {
			fmt.Printf("%s: %s\n", ansi.Color("ERROR:", "red+b"), ansi.Color("Roundtrip parsing failed. Second result:\n", "red"))
			fmt.Printf(ansi.Color("%s\n", "magenta"), secondString)
			continue
		}
		fmt.Printf(ansi.Color("%s\n", "green"), "\n-- Roundtrip parsing ok\n")

		// Use this to limit the number of tests run while debugging
		// if n == 1 {
		// 	break
		// }
	}
}
