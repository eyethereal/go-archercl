package config

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)

%%{
   machine ACLParser;
   write data; 
}%%

func findCol(data string, ts int) int {
    ix := strings.LastIndex(data[:ts], "\n")
    if ix == -1 {
        return 0;
    }
    return ts - ix
}

// This function is necessary because the go unqoute routine will only
// allow a single character inside single quotes, which is not what js will do
func singlesToDoubles(q string) string {
    if q[0] == '"' {
        return q
    }

    q = strings.Replace(q, "\\'", "'", -1)
    q = strings.Replace(q, "\"", "\\\"", -1)

    return "\"" + q[1:len(q)-1] + "\""
}

func (node *AclNode) ParseString(data string, location *ParseLocation) (err error) {

    // Handle the user not supplying anything as a base reference for the parse location
    if location == nil {
        location = new(ParseLocation)
    }

    // These are the required variables for the ragel FSM code. data is also required
    // but is an input parameter
    cs, p, pe := 0, 0, len(data)
    eof := pe // Marks this as the last (i.e. only) block of data in the file
    ts, te, act := 0, 0, 0
    _ = act

    stack := make([]int, 0)
    top := 0

    nodeStack := make([]*AclNode, 1)
    nodeStack[0] = node
    keyNameStack := make([][]string, 1)
    keyNameStack[0] = make([]string, 0)
    knsp := 0

    firstValue := true
    multiLineArray := false

    var value *AclNode

    resetKey := func() {
        firstValue = true
        multiLineArray = false
        keyNameStack[knsp] = keyNameStack[knsp][:0]
        value = nil
    }

    appendKey := func(k string) {
        keyNameStack[knsp] = append(keyNameStack[knsp], k)
    }

    appendQuotedKey := func(q string) error {
        q = singlesToDoubles(q)
        k, err := strconv.Unquote(q)
        if err != nil {
            return errors.New("Can not parse "+q+" : "+err.Error())
        }
        keyNameStack[knsp] = append(keyNameStack[knsp], k)
        return nil
    }

    ensureNode := func() {
        if value != nil {
            return 
        }

        // Take the current node and make sure the path from current to
        // the keyName exists
        // Don't consider keyName == 0 because that shouldn't have happened
        if len(keyNameStack[knsp]) == 0 {
            panic("keyName has no elements")
        }

        // Get / build all intermediate nodes
        cNode := nodeStack[len(nodeStack)-1]
        var nextNode *AclNode
        for ix:=0; ix < len(keyNameStack[knsp]); ix ++ {
            name := keyNameStack[knsp][ix]
            if name[0] == '!' {
                name = name[1:]
                if firstValue {
                    nextNode = nil
                    firstValue = false
                } else {
                    nextNode = cNode.Children[name]
                }
            } else {
                nextNode = cNode.Children[name]
            }

            if nextNode == nil {
                // Gotta make it and add it into the list of childrenes
                nextNode = NewAclNode()
                cNode.Children[name] = nextNode
            }
            cNode = nextNode
        }

        // Set the "value" node to the cNode
        value = cNode
    }

    stringValue := func(v string) {
        // fmt.Printf("stringValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        ensureNode()
        value.Values = append(value.Values, v)
    }

    quotedStringValue := func(q string) error {
        // fmt.Printf("quotedStringValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        ensureNode()
        q = singlesToDoubles(q)
        v, err := strconv.Unquote(q)
        if err != nil {
            return errors.New("Can not parse "+q+" : "+err.Error())
        }

        value.Values = append(value.Values, v)
        return nil
    }

    integerValue := func(v string) error {
        // fmt.Printf("integerValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)
        ensureNode()

        i, err := strconv.ParseInt(v, 0, 0)
        if err != nil {
            return err
        }
        value.Values = append(value.Values, i)
        return nil
    }

    floatValue := func(v string) error {
        ensureNode()

        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return err
        }
        value.Values = append(value.Values, f)
        return nil
    }

    descend := func() {
        // fmt.Printf("Descend\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        if len(keyNameStack[knsp]) == 0 {
            // It's not really a descent. Since no key is defined, we just
            // collapse this into the current node. However, because we need
            // to be able to ascend, we have to push still.
            nodeStack = append(nodeStack, nodeStack[knsp])
            keyNameStack = append(keyNameStack, keyNameStack[knsp])
            knsp++
            return
        }

        // Make sure the currently named node exists. This will also
        // set value to the currently named node
        ensureNode()

        // Because we have descended we must push
        nodeStack = append(nodeStack, value)
        keyNameStack = append(keyNameStack, make([]string, 0))
        value = nil
        knsp++
    }

    ascend := func() bool {
        // fmt.Printf("Ascend\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        // Can't pop beyond the first node
        if len(nodeStack) <= 1 {
            return false
        }

        // Ok, just a pop
        nodeStack = nodeStack[0:len(nodeStack)-1]
        keyNameStack = keyNameStack[0:knsp]
        knsp--
        resetKey()
        return true
    }


    %%{
        # Keep the fcall stack large enough. Don't bother ever shrinking it
        prepush {
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        }

        newline = '\n' @{ location.Line++ };
        any_count_line = any | newline;

        # Consume a C style comment
        c_comment := any_count_line* :>> '*/' @{fret;};

        cpp_comment = "//" [^\n]*;
        bash_comment = "#" [^\n]*;
        sql_comment = "--" [^\n]*;
        one_line_comment = ( cpp_comment | bash_comment | sql_comment );

        ignorable_ws = any_count_line - 0x21..0x7e ;

        # For unqouted identifiers
        alnum_u = alnum | '_';
        alpha_u = alpha | '_';

        identifier = alpha_u alnum_u*;

        notable_identifier = ("!" | alpha_u) alnum_u*;

        sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
        sliteral = "'" . sliteralChar* . "'";

        dliteralChar = [^"\\] | newline | ( '\\' . any_count_line );
        dliteral = '"' . dliteralChar* . '"';

        integer = ( "+" | "-" )? digit+;
        float = integer "." digit+;
        hex_integer = '0x' xdigit+;

        obj_value := |*
            # Identifiers that are not in quotes
            identifier => { 
                // fmt.Printf("Value Identifier %v\n", data[ts:te])
                stringValue(data[ts:te])
            };

            (sliteral | dliteral) {
                // fmt.Printf("Quoted literal %v\n", data[ts:te])
                err = quotedStringValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing quoted value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Integers
            integer {
                // fmt.Printf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Float
            float {
                // fmt.Printf("Value Float %v\n", data[ts:te])
                err = floatValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing float value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Hex
            hex_integer {
                // fmt.Printf("Value Hex %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing hex value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            (";" | newline) {
                if !multiLineArray {
                    resetKey()
                    fret;
                } // else, let it end the current value, but not the array
            };

            # Handle this at the higher level. This rule is what allows a 
            # : or = after a key name but before an object name definition starts
            "{" {
                p--
                fret;
            };

            "}" {
                resetKey()
                p--
                fret;
            };

            # If a c style comment start, just consume all of it
            ("/*") => { fcall c_comment; };

            # If a c++, sql, or bash style comment starts consume it.
            # Because these are single line it's easier so we
            # don't go to the other machine. 
            one_line_comment;

            # Commas imply arrays, but that's not actually relevant to us. We
            # allow multiple values to be arrays
            ',' {
                // fmt.Printf("Comma\n")
            };

            # An [  needs to move this to a multi-line mode. Commas can still be optional.
            '[' {
                multiLineArray = true
            };

            # End the array (if we are in one), but generally just return
            ']' {
                resetKey()
                fret;
            };

            # Whitespace is ok
            ignorable_ws {
                // fmt.Printf("Value Whitespace '%v'\n", data[ts:te])
            };

            # Everything else is an error
            any {
                location.Message = fmt.Sprintf("Syntax error while looking for a value: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                fgoto *ACLParser_error;
            };            

        *|;

        main := |*

            # Identifiers that are not in quotes
            notable_identifier { 
                //fmt.Printf("Key Identifier %v\n", data[ts:te])
                appendKey(data[ts:te])
            };

            (sliteral | dliteral) {
                k := data[ts:te]
                if len(k) < 3 {
                    location.Message = "Key names may not be empty"
                    location.Col = findCol(data,ts)
                    fgoto *ACLParser_error;
                } else {
                    err = appendQuotedKey(k)
                    if err != nil {
                        location.Message = fmt.Sprintf("Error parsing quoted key name: %v", err)
                        location.Col = findCol(data, ts)
                        fgoto *ACLParser_error;
                    }
                }
            };

            [:=] {
                //fmt.Printf("Call obj_value '%v'\n", data[ts:te])
                multiLineArray = false
                fcall obj_value;
            };

            "{" {
                descend()
            };

            "}" {
                if !ascend() {
                    location.Message = "Syntax error unmatched } "
                    location.Col = findCol(data, ts)              
                    fgoto *ACLParser_error;
                }
            };

            # If a c style comment start, just consume all of it
            ("/*") => { fcall c_comment; };

            # Also consume single line comments, C++, bash, or sql style
            one_line_comment;

            # Anything that is whitespacey or control codey just
            # gets ignored (and potentially counted)
            ignorable_ws {
                // fmt.Printf("Other '%v'\n", data[ts:te])
            };

            ";" {
                resetKey()
            };

            # Commas at the top level are ignorable. Without this you get an error
            # trying to parse what is generated by .toString() - so bad!!!
            "," {

            };

            # Everything else is an error
            any {
                location.Message = fmt.Sprintf("Syntax error while looking for a key name: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                fgoto *ACLParser_error;
            };
        *|;


        write init;
        write exec;
    }%%

    if cs == ACLParser_error {
        if len(location.Message) == 0 {
            location.Message = "Configuration file syntax error"
        }
        fmt.Printf("Error! %v\n", location)
        return location
    }

    return nil
}
