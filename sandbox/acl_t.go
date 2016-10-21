package main

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/mgutz/ansi"
	"strconv"
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
}

func main() {

	for n, sample := range samples {
		fmt.Printf(ansi.Color("\n==========================================\n", "black+b"))
		fmt.Printf(ansi.Color("Test #%d:\n", "black+b"), n+1)
		fmt.Printf(ansi.Color("%s\n", "blue"), sample)
		node := config.NewAclNode()
		err := node.ParseString(sample, nil)
		if err != nil {
			fmt.Printf("%s: %s\n", ansi.Color("ERROR:", "red+b"), ansi.Color(err.Error(), "red"))
			continue
		}

		fmt.Printf(ansi.Color("%s\n", "green"), node.String())
		fmt.Printf(ansi.Color("%s\n", "magenta"), strconv.Quote(node.String()))
	}
}
