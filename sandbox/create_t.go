package main

import (
	"fmt"
	"github.com/eyethereal/go-archercl"
)

func main() {
	root := config.NewAclNode()

	root.SetValAt(10, "first")
	root.SetValAt("World", "Hello")
	root.SetValAt(2.3, "second", "deeper", "level")

	fmt.Printf("%s", root.ColoredString())
}
