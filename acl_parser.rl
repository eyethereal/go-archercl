package archercl

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "log"
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

// A key value context defines a place in the node tree against which identifiers
// are collected to create a key path naming a second deeper location in the tree.
// When values are encountered they are they added to this second deeper location,
// and when the parsing of values completes we are able to return to the original
// level of context.
type kvCtx struct {

    // The base of this context against which the keyPath is to be judged
    node *AclNode

    // A collection of keys which identifies a deeper location into the tree
    keyPath []string

    // Whether we are inside of an array scope or not. This allows us to understand
    // if a second word is a value or a key
    inArray bool

    // Whether an equals sign was used to move to value mode or not. This gets
    // annotated onto the value mode when values are attached
    usesEqual bool
}

func (root *AclNode) ParseString(data string, location *ParseLocation) (err error) {
    return root.ParseStringWithLogger(data, location, nil)
}

func (root *AclNode) ParseStringWithLogger(data string, location *ParseLocation, logger *log.Logger) (err error) {

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

    ///// End of required ragel FSM variables------

    lprintf := func(fmt string, v ...interface{}) {
        if logger == nil {
            return
        }
        logger.Printf(fmt, v...)
    }

    // In addition to token parsing, we need to maintain our own context in a stack
    // as we move into and out of object and array collections of values.
    ctxStack := make([]kvCtx, 0, 10)

    // We'll seed the ctxStack with a root node after that function is defined...

    // Resets the key at the end of a value, but does not change the context.
    // This is how we return to the current context after adding a value that
    // was named for a possibly multiple levels deep object which implied multiple
    // intermediate nodes.
    resetKey := func() {
        // Because the context stack is real structs not pointers we have to reference
        // it directly by name not via a copy or else we don't modify what we think we
        // are modifying
        ctxStack[len(ctxStack)-1].keyPath = ctxStack[len(ctxStack)-1].keyPath[:0]
        ctxStack[len(ctxStack)-1].usesEqual = false
    }

    appendKey := func(k string) error {
        ctx := ctxStack[len(ctxStack)-1]

        if ctx.inArray {        
            // Keys are not allowed in array contexts, so something has gone terribly wrong.
            // Most likely we have somehow returned to key mode when we meant to stay in
            // value mode somehow.
            return fmt.Errorf("Keys are not allowed in an array scope")
        }

        ctxStack[len(ctxStack)-1].keyPath = append(ctxStack[len(ctxStack)-1].keyPath, k)
        return nil
    }

    appendQuotedKey := func(q string) error {
        q = singlesToDoubles(q)
        k, err := strconv.Unquote(q)
        if err != nil {
            return errors.New("Can not parse "+q+" : "+err.Error())
        }

        return appendKey(k)
    }


    // Find the currently named target, creating it if necessary
    findCurrentTarget := func() *AclNode {
        ctx := ctxStack[len(ctxStack)-1]

        // Start with the current node
        target := ctx.node

        lprintf("findCurrentTarget ctx=%v\n", ctx)        

        // Need to follow the path (if any) down, possibly creating nodes
        // as we go along
        for ix:=0; ix<len(ctx.keyPath); ix++ {
            name := ctx.keyPath[ix]

            var next *AclNode
            if name[0] == '!' {
                // Gotta nuke any existing things, so we do
                // that by purposely not looking up the node
                name = name[1:]
            } else {
                // Try to get an existing node (which might fail)
                next = target.Children[name]
            }

            if next == nil {
                // Oh hey, it's new (or a replacement in the reset case)
                next = NewAclNode()
                target.Children[name] = next
                target.OrderedChildNames = append(target.OrderedChildNames, name)
            }

            target = next
        }

        return target        
    }

    // Attach a value at the currently named location in the context
    attachValue := func(v interface{}) {
        target := findCurrentTarget()
        target.Values = append(target.Values, v)

        // Record whether this used an equals sign or not
        target.UsesEquals = ctxStack[len(ctxStack)-1].usesEqual
    }

    stringValue := func(v string) {
        // fmt.Printf("stringValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        attachValue(v)
    }

    quotedStringValue := func(q string) error {
        // fmt.Printf("quotedStringValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        q = singlesToDoubles(q)
        v, err := strconv.Unquote(q)
        if err != nil {
            return errors.New("Can not parse "+q+" : "+err.Error())
        }

        attachValue(v)
        return nil
    }

    integerValue := func(v string) error {
        // fmt.Printf("integerValue\n\tkeyNameStack=%v\n\tnodeStack=%v\n", keyNameStack, nodeStack)
        // fmt.Printf("\tvalue=%v\n", value)

        i, err := strconv.ParseInt(v, 0, 0)
        if err != nil {
            return err
        }
        attachValue(i)
        return nil
    }

    floatValue := func(v string) error {
        // ensureNode()

        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return err
        }
        attachValue(f)
        return nil
    }

    pushContext := func(next *AclNode) {        
        ctxStack = append(ctxStack, kvCtx{
            node: next,
            keyPath: make([]string,0),
            })

        // We need to propogate uses equals into the next context unless it is the root element
        // that was just pushed
        if len(ctxStack) > 1 {
            ctxStack[len(ctxStack)-1].usesEqual = ctxStack[len(ctxStack)-2].usesEqual
        }
    }

    popContext := func() error {
        // Can't pop beyond the first node
        if len(ctxStack) <= 1 {
            return fmt.Errorf("Mismatched scope braces. Trying to pop the root scope.")
        }

        ctxStack = ctxStack[:len(ctxStack)-1]
        return nil
    }

    startObject := func() {
        ctx := ctxStack[len(ctxStack)-1]

        if ctx.inArray {
            // We need to make a new node, which will become our context, and instead
            // of being a child, it will be a value in the current context
            next := NewAclNode()
            ctx.node.Values = append(ctx.node.Values, next)

            // Now move into that new context
            pushContext(next)

        } else {
            // Just push our current target as the new context
            pushContext(findCurrentTarget())
        }
    }

    endObject := func() error {
        ctx := ctxStack[len(ctxStack)-1]

        if ctx.inArray {
            return fmt.Errorf("Mismatched brackets. Object end found while not inside of an object scope.")
        }

        // All the logic for where values get attached happens on the descent so
        // all we have to do is pop the context and we're back in business
        return popContext()
    }

    // Start a new array which might be a sub array or not
    startArray := func() {
        ctx := ctxStack[len(ctxStack)-1]

        if ctx.inArray {
            // It is a sub-array which means we push a new context which will collect
            // the values into a new "node", but that node's value array will actually
            // be placed into the value array of the current node
            shadow := NewAclNode()
            // subArrayValues := shadow.Values
            attachValue(shadow)

            // The new shadow node isn't connected to anything, so keys attached
            // to it aren't going to be interesting at all. It's just values
            pushContext(shadow)
        } else {
            // It's not a sub-array, so we are really just referring to the current
            // target level object. However, we go ahead and push a new context which
            // refers to the current target's value array, instead of a new sub-array.

            pushContext(findCurrentTarget())
        }

        // Mark this new context as being an array target not a map target
        ctxStack[len(ctxStack)-1].inArray = true
    }

    endArray := func() error {
        ctx := ctxStack[len(ctxStack)-1]

        if !ctx.inArray {
            return fmt.Errorf("Mismatched brackets. Array end found while not inside of an array scope.")
        }

        // Pop the existing context
        return popContext()
    }

    // Used to determine what mode is appropriate after an array or object scope has ended
    inArray := func() bool {
        return ctxStack[len(ctxStack)-1].inArray
    }


    // We start with a single root, which is node on which ParseString is called
    pushContext(root)


    // Writing ragel comments with #// means they keep syntax highlighting working in sublime
    %%{
        #// Keep the fcall stack large enough. Don't bother ever shrinking it 
        prepush {
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        }

        newline = '\n' @{ 
            location.Line++ 
            // The entire stack becomes multiline
            for i:=0; i<len(ctxStack); i++ {
                ctx := &ctxStack[i]
                ctx.node.IsMultiline = true

                // Now we also have to walk down the keyPath in this context to find
                // intermediate nodes which are not necessarily in the NEXT context on
                // the stack and make sure they all get set to multiline also. This is
                // necessary because at the time the nodes were created we may not
                // have known that multiline was going to be a thing
                nodeCursor := ctx.node
                for kix:=0; nodeCursor != nil && kix<len(ctx.keyPath); kix++ {
                    key := ctx.keyPath[kix]
                    nodeCursor.IsMultiline = true
                    nodeCursor = nodeCursor.Children[key]
                }
            }
        };
        any_count_line = any | newline;

        #// Consume a C style comment
        c_comment := any_count_line* :>> '*/' @{fret;};

        cpp_comment = "//" [^\n]*;
        bash_comment = "#" [^\n]*;
        sql_comment = "--" [^\n]*;
        one_line_comment = ( cpp_comment | bash_comment | sql_comment );

        ignorable_ws = any_count_line - 0x21..0x7e ;

        #// For unqouted identifiers
        alnum_u = alnum | '_';
        alpha_u = alpha | '_';

        identifier = alpha_u alnum_u*;

        notable_identifier = ("!" | alpha_u) alnum_u*;

        sliteralChar = [^'\\] | newline | ( '\\' . any_count_line );
        # '
        sliteral = "'" . sliteralChar* . "'";

        dliteralChar = [^"\\] | newline | ( '\\' . any_count_line );
        # "
        dliteral = '"' . dliteralChar* . '"';

        integer = ( "+" | "-" )? digit+;
        float = integer "." digit+;
        hex_integer = '0x' xdigit+;

        value_mode := |*
            # Identifiers that are not in quotes
            identifier => { 
                lprintf("Value Identifier %v\n", data[ts:te])
                stringValue(data[ts:te])
            };

            (sliteral | dliteral) {
                lprintf("Quoted literal %v\n", data[ts:te])
                err = quotedStringValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing quoted value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Integers
            integer {
                lprintf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Float
            float {
                lprintf("Value Float %v\n", data[ts:te])
                err = floatValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing float value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Hex
            hex_integer {
                lprintf("Value Hex %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing hex value: %v", err)
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }
            };

            # Since semicolons end value mode and return to key they only make
            # sense if not in an array scope
            ';' {
                if inArray() {
                    location.Message = "Invalid ';' found while in an array context"
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;
                }

                // Since we aren't in an array this allows a return to key mode with
                // a reset key in the same way that a newline not in an array does
                resetKey()
                fret;
            };

            # Newlines are value separators, but whether they change mode or not
            # depends on whether we are in an array or object context
            newline {
                lprintf("Value newline. inArray()=%v\n", inArray())
                if !inArray() {
                    lprintf("  return to key mode\n")
                    // Behave the same as a semicolon does if in an object context.
                    resetKey()
                    fret;
                }

                // Because we are in an array we stay in value mode. There are no
                // keys allowed in array contexts so nothing to reset.
            };

            # We need to descend into a new object scope so that means we
            # go back into key mode by returning
            '{' {
                startObject();
                fret;
            };

            # Object ends pop a context and then if the new context is an array, which
            # basically means it was probably a multi-line array, then we need to stay
            # in value mode. Otherwise return to key mode.
            '}' {                
                err := endObject()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    fgoto *ACLParser_error;
                }

                if !inArray() {
                    // Back to key mode
                    lprintf("Object end, return to key mode")
                    resetKey()
                    fret;
                }
                // Since we are in an array stay in value mode so that we
                // keep adding values into the named target
            };

            # Arrays introduce a new context and we stay in value mode
            '[' {
                startArray()
            };

            #// The end of our array (unless we weren't in one...). The mode we should
            # be in depends on the new context after the array has ended.
            ']' {
                err := endArray()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    fgoto *ACLParser_error;
                }

                if !inArray() {
                    // Back to key mode
                    resetKey()
                    fret;
                }                
            };


            # If a c style comment start, just consume all of it
            ("/*") => { fcall c_comment; };

            # If a c++, sql, or bash style comment starts consume it.
            #// Because these are single line it's easier so we
            #// don't go to the other machine. 
            one_line_comment;

            #// In arrays commas are syntatic sugar. Inside objects they end
            #// the value and those require a mode change.
            ',' {
                if !inArray() {
                    // Not sugar, return to key mode
                    lprintf("Comma in object, return to key mode\n")
                    resetKey()
                    fret;
                } else {
                    lprintf("Comma in array, ignore as sugar\n")
                }
            };

            # Whitespace is ok, it ends values. This includes newlines, but those
            # are also handled earlier in case they caused a mode change.
            ignorable_ws {
                lprintf("Value Whitespace '%v'\n", data[ts:te])
            };

            # Everything else is an error
            any {
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a value.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                fgoto *ACLParser_error;
            };            

        *|;

        # The main mode corresponds to searching for a key. We are either going to collect 
        # identifiers into a key path and then switch into value mode, or we are going to
        # descend or ascend into a different key scope, but remain in key searching mode
        # with a different context.
        # Because arrays are a different type of value, those are handled in value mode
        main := |*

            # Identifiers that are not in quotes
            notable_identifier { 
                lprintf("Key Identifier %v\n", data[ts:te])
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

            #// Key/Value separators are mandatory because they are a major mode change for us
            [:=] {
                ctxStack[len(ctxStack)-1].usesEqual = (data[ts] == '=')

                lprintf("Change to value_mode '%v' usesEqual=%v\n", data[ts:te], ctxStack[len(ctxStack)-1].usesEqual)
                fcall value_mode;
            };

            # Object starts push a new context but we then want to stay in key mode
            '{' {
                startObject()
            };

            # Object ends pop a context and then if the new context is an array, which
            # basically means it was probably a multi-line array, then we need to be
            # back in value mode
            '}' {                
                err := endObject()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    fgoto *ACLParser_error;
                }

                if inArray() {
                    fcall value_mode;
                }

                // Do have to reset the key though
                resetKey()
            };

            # Arrays introduce a new context and move us into value mode
            '[' {
                startArray()
                fcall value_mode;
            };

            #// It shouldn't happen that we end an array while in key mode
            ']' {
                location.Message = "Array scope ended while in key mode indicates a parser state error."
                location.Col = findCol(data, ts)
                fgoto *ACLParser_error;
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

            #// OLD COMMENT: Commas at the top level are ignorable. Without this you get an error
            #// trying to parse what is generated by .toString() - so bad!!!
            #// Commas and periods are key separators. We recognize them as unique tokens but don't
            #// need to do anything else with them
            ('.' | ',') {

            };

            #// Duplicate semicolons aren't interesting. Semicolons after we have collected
            #// some keys probably indicate a user mistake though so it seems worthwhile to
            #// report it. Really it's more of a warning, but we don't have that distinction.
            #// The alternate approach would be to just reset the key paths throwing things
            #// away, but it seems wiser to not allow this condition.
            ';' {
                if len(ctxStack[len(ctxStack)-1].keyPath) > 0 {
                    location.Message = "Key names found without a value."
                    location.Col = findCol(data, ts)
                    fgoto *ACLParser_error;                    
                }
            };

            # Everything else is an error
            any {
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a key.", data[ts:ts+1])
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
