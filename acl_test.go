package config

import (
	"testing"
)

func init() {
	ColoredLoggingToConsole()
}

func Test_Comment(t *testing.T) {

	src := `
// This is a comment
a = 1

/* And this is also a comment
that is multi line*/
b = 2

/* Yeah? * */ c = 3 /**/
`
	node := NewAclNode()
	err := node.ParseString(src, nil)
	if err != nil {
		t.Fatal(err)
	}

	if node.Child("a").AsInt() != 1 {
		t.Fatalf("a didn't parse properly")
	}
	if node.Child("b").AsInt() != 2 {
		t.Fatalf("b didn't parse properly")
	}
	if node.Child("c").AsInt() != 3 {
		t.Fatalf("v didn't parse properly")
	}
}

func Test_ValuesTypes(t *testing.T) {
	src := `
    
    // Value Types
    unquoted = def
    squoted = 's\'in"gle'
    dquoted = "dou'ble"
    integer1 = 1
    integer2 = 123
    integer3 = +123
    integer4 = -123
    float1 = 0.1
    float2 = +0.1
    float3 = -0.1

    // Arrays
    array unquoted = one two three
    array quoted = "the first" "the second" 
`

	node := NewAclNode()
	err := node.ParseString(src, nil)
	if err != nil {
		t.Fatal(err)
	}

	if node.Child("unquoted").AsString() != "def" {
		t.Fatal("Wrong value for unquoted")
	}
	if node.Child("squoted").AsString() != "s'in\"gle" {
		t.Fatal("Wrong value for squoted")
	}
	if node.Child("dquoted").AsString() != "dou'ble" {
		t.Fatal("Wrong value for unquoted")
	}
	if node.Child("integer1").AsInt() != 1 {
		t.Fatal("Wrong value for integer1")
	}
	if node.Child("integer2").AsInt() != 123 {
		t.Fatal("Wrong value for integer2")
	}
	if node.Child("integer3").AsInt() != 123 {
		t.Fatal("Wrong value for integer3")
	}
	if node.Child("integer4").AsInt() != -123 {
		t.Fatal("Wrong value for integer4")
	}
	if node.Child("float1").AsFloat() != 0.1 {
		t.Fatal("Wrong value for float1")
	}
	if node.Child("float2").AsFloat() != 0.1 {
		t.Fatal("Wrong value for float2")
	}
	if node.Child("float3").AsFloat() != -0.1 {
		t.Fatal("Wrong value for float3")
	}

	ac := node.Child("array")
	if ac == nil {
		t.Fatal("No child 'array'")
	}

	acn := ac.Child("unquoted")
	if acn == nil {
		t.Fatal("No child 'array unquoted'")
	}
	if acn.Len() != 3 {
		t.Fatal("Wrong length for array unquoted")
	}
	if acn.AsStringN(2) != "three" {
		t.Fatal("Wrong value for array unquoted 3")
	}

	acq := ac.Child("quoted")
	if acq == nil {
		t.Fatal("No child 'array quoted'")
	}
	if acq.Len() != 2 {
		t.Fatal("Wrong length for array quoted")
	}
	if acq.AsStringN(1) != "the second" {
		t.Fatal("Wrong value for array quoted 2")
	}

}

func Test_ValuesObject(t *testing.T) {
	src := `
    sub {
        one: 1
        second {
            two: 2
        }
    }
`

	node := NewAclNode()
	err := node.ParseString(src, nil)
	if err != nil {
		t.Fatal(err)
	}

	alog.Infof("Node is %s", node.String())

	obj := node.Child("sub")
	if obj == nil {
		t.Fatal("No sub object")
	}

	if obj.Child("one").AsInt() != 1 {
		t.Fatal("Wrong value for sub one")
	}

	second := obj.Child("second")
	if second == nil {
		t.Fatal("No second object")
	}

	if second.Child("two").AsInt() != 2 {
		t.Fatal("Wrong value for sub second two")
	}
}

func Test_StringOne(t *testing.T) {

	root := NewAclNode()
	n := NewAclNode()
	n.Values = append(n.Values, 1)
	n.Values = append(n.Values, 2.3)
	n.Values = append(n.Values, true)
	n.Values = append(n.Values, "abc")
	root.Children["one"] = n

	n = NewAclNode()
	n.Values = append(n.Values, "world")
	root.Children["hello"] = n

	str := root.String()

	should := "{\n\t\"hello\": \"world\",\n\t\"one\": [\n\t\t1,\n\t\t2.3,\n\t\ttrue,\n\t\t\"abc\",\n\t],\n}"
	if str != should {
		t.Fatal("Expected: \n" + should + "But got\n" + str)
	}
}
