
//line acl_parser.rl:1
package archercl

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "log"
)


//line acl_parser.go:15
var _ACLParser_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3, 
	1, 4, 1, 5, 1, 6, 1, 8, 
	1, 9, 1, 10, 1, 11, 1, 12, 
	1, 13, 1, 14, 1, 15, 1, 16, 
	1, 17, 1, 18, 1, 19, 1, 20, 
	1, 21, 1, 22, 1, 23, 1, 24, 
	1, 25, 1, 26, 1, 27, 1, 28, 
	1, 29, 1, 30, 1, 31, 1, 32, 
	1, 33, 1, 34, 1, 35, 1, 36, 
	1, 37, 1, 38, 2, 0, 7, 2, 
	0, 31, 
}

var _ACLParser_key_offsets []byte = []byte{
	0, 0, 3, 4, 7, 8, 10, 13, 
	16, 17, 20, 21, 23, 29, 54, 61, 
	64, 65, 68, 69, 71, 71, 98, 101, 
	102, 105, 107, 110, 112, 115, 117, 121, 
	127, 
}

var _ACLParser_trans_keys []byte = []byte{
	10, 34, 92, 10, 10, 39, 92, 10, 
	10, 42, 10, 42, 47, 10, 34, 92, 
	10, 10, 39, 92, 10, 48, 57, 48, 
	57, 65, 70, 97, 102, 10, 33, 34, 
	35, 39, 44, 45, 46, 47, 58, 59, 
	61, 91, 92, 93, 94, 96, 123, 125, 
	36, 64, 65, 122, 124, 126, 95, 48, 
	57, 65, 90, 97, 122, 10, 34, 92, 
	10, 10, 39, 92, 45, 42, 47, 10, 
	34, 35, 39, 43, 44, 45, 47, 48, 
	59, 91, 92, 93, 94, 96, 123, 125, 
	33, 46, 49, 57, 58, 64, 65, 122, 
	124, 126, 10, 34, 92, 10, 10, 39, 
	92, 48, 57, 46, 48, 57, 48, 57, 
	45, 48, 57, 42, 47, 46, 120, 48, 
	57, 48, 57, 65, 70, 97, 102, 95, 
	48, 57, 65, 90, 97, 122, 
}

var _ACLParser_single_lengths []byte = []byte{
	0, 3, 1, 3, 1, 2, 3, 3, 
	1, 3, 1, 0, 0, 19, 1, 3, 
	1, 3, 1, 2, 0, 17, 3, 1, 
	3, 0, 1, 0, 1, 2, 2, 0, 
	1, 
}

var _ACLParser_range_lengths []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 1, 3, 3, 3, 0, 
	0, 0, 0, 0, 0, 5, 0, 0, 
	0, 1, 1, 1, 1, 0, 1, 3, 
	3, 
}

var _ACLParser_index_offsets []byte = []byte{
	0, 0, 4, 6, 10, 12, 15, 19, 
	23, 25, 29, 31, 33, 37, 60, 65, 
	69, 71, 75, 77, 80, 81, 104, 108, 
	110, 114, 116, 119, 121, 124, 127, 131, 
	135, 
}

var _ACLParser_trans_targs []byte = []byte{
	1, 13, 2, 1, 1, 1, 3, 13, 
	4, 3, 3, 3, 5, 6, 5, 5, 
	6, 20, 5, 7, 21, 8, 7, 7, 
	7, 9, 21, 10, 9, 9, 9, 27, 
	21, 31, 31, 31, 21, 13, 14, 15, 
	16, 17, 13, 18, 13, 19, 13, 13, 
	13, 13, 13, 13, 13, 13, 13, 13, 
	13, 14, 13, 13, 14, 14, 14, 14, 
	13, 1, 13, 2, 1, 13, 16, 3, 
	13, 4, 3, 16, 13, 13, 16, 13, 
	0, 21, 22, 23, 24, 25, 21, 28, 
	29, 30, 21, 21, 21, 21, 21, 21, 
	21, 21, 21, 26, 21, 32, 21, 21, 
	7, 21, 8, 7, 21, 23, 9, 21, 
	10, 9, 26, 21, 11, 26, 21, 27, 
	21, 23, 26, 21, 21, 23, 21, 11, 
	12, 26, 21, 31, 31, 31, 21, 32, 
	32, 32, 32, 21, 13, 13, 13, 13, 
	21, 21, 21, 21, 21, 21, 13, 13, 
	13, 13, 13, 13, 21, 21, 21, 21, 
	21, 21, 21, 21, 21, 21, 21, 
}

var _ACLParser_trans_actions []byte = []byte{
	1, 47, 0, 0, 1, 0, 1, 47, 
	0, 0, 1, 0, 1, 0, 0, 1, 
	0, 3, 0, 1, 11, 0, 0, 1, 
	0, 1, 11, 0, 0, 1, 0, 0, 
	43, 0, 0, 0, 43, 80, 0, 9, 
	0, 9, 63, 0, 63, 0, 49, 65, 
	49, 55, 67, 57, 67, 67, 51, 53, 
	67, 0, 67, 61, 0, 0, 0, 0, 
	69, 1, 47, 0, 0, 71, 0, 1, 
	47, 0, 0, 0, 73, 59, 0, 73, 
	0, 77, 9, 0, 9, 0, 25, 0, 
	0, 9, 13, 19, 29, 21, 29, 29, 
	15, 17, 29, 9, 29, 0, 29, 27, 
	1, 11, 0, 0, 39, 0, 1, 11, 
	0, 0, 9, 41, 0, 9, 33, 0, 
	35, 0, 9, 41, 23, 0, 41, 0, 
	0, 9, 33, 0, 0, 0, 37, 0, 
	0, 0, 0, 31, 75, 75, 75, 75, 
	45, 45, 45, 45, 43, 43, 69, 73, 
	71, 73, 73, 73, 41, 39, 41, 41, 
	33, 35, 41, 41, 33, 37, 31, 
}

var _ACLParser_to_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 5, 0, 0, 
	0, 0, 0, 0, 0, 5, 0, 0, 
	0, 0, 0, 0, 0, 5, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 
}

var _ACLParser_from_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 7, 0, 0, 
	0, 0, 0, 0, 0, 7, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 
}

var _ACLParser_eof_trans []byte = []byte{
	0, 144, 144, 144, 144, 0, 0, 148, 
	148, 148, 148, 150, 150, 0, 151, 156, 
	153, 156, 156, 156, 0, 0, 164, 158, 
	164, 164, 165, 162, 164, 164, 165, 166, 
	167, 
}

const ACLParser_start int = 13
const ACLParser_first_final int = 13
const ACLParser_error int = 0

const ACLParser_en_c_comment int = 5
const ACLParser_en_value_mode int = 21
const ACLParser_en_main int = 13


//line acl_parser.rl:14


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
    
//line acl_parser.go:475
	{
	cs = ACLParser_start
	top = 0
	ts = 0
	te = 0
	act = 0
	}

//line acl_parser.go:484
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_acts = int(_ACLParser_from_state_actions[cs])
	_nacts = uint(_ACLParser_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		 _acts++
		switch _ACLParser_actions[_acts - 1] {
		case 3:
//line NONE:1
ts = p

//line acl_parser.go:507
		}
	}

	_keys = int(_ACLParser_key_offsets[cs])
	_trans = int(_ACLParser_index_offsets[cs])

	_klen = int(_ACLParser_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _ACLParser_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _ACLParser_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_ACLParser_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _ACLParser_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _ACLParser_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
_eof_trans:
	cs = int(_ACLParser_trans_targs[_trans])

	if _ACLParser_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_ACLParser_trans_actions[_trans])
	_nacts = uint(_ACLParser_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _ACLParser_actions[_acts-1] {
		case 0:
//line acl_parser.rl:332
 
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
        
		case 1:
//line acl_parser.rl:355
top--; cs = stack[top]
goto _again

		case 4:
//line NONE:1
te = p+1

		case 5:
//line acl_parser.rl:391
te = p+1
{
                lprintf("Quoted literal %v\n", data[ts:te])
                err = quotedStringValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing quoted value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 6:
//line acl_parser.rl:436
te = p+1
{
                if inArray() {
                    location.Message = "Invalid ';' found while in an array context"
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }

                // Since we aren't in an array this allows a return to key mode with
                // a reset key in the same way that a newline not in an array does
                resetKey()
                top--; cs = stack[top]
goto _again

            }
		case 7:
//line acl_parser.rl:451
te = p+1
{
                lprintf("Value newline. inArray()=%v\n", inArray())
                if !inArray() {
                    lprintf("  return to key mode\n")
                    // Behave the same as a semicolon does if in an object context.
                    resetKey()
                    top--; cs = stack[top]
goto _again

                }

                // Because we are in an array we stay in value mode. There are no
                // keys allowed in array contexts so nothing to reset.
            }
		case 8:
//line acl_parser.rl:466
te = p+1
{
                startObject();
                top--; cs = stack[top]
goto _again

            }
		case 9:
//line acl_parser.rl:474
te = p+1
{                
                err := endObject()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    cs = (ACLParser_error)
goto _again

                }

                if !inArray() {
                    // Back to key mode
                    lprintf("Object end, return to key mode")
                    resetKey()
                    top--; cs = stack[top]
goto _again

                }
                // Since we are in an array stay in value mode so that we
                // keep adding values into the named target
            }
		case 10:
//line acl_parser.rl:493
te = p+1
{
                startArray()
            }
		case 11:
//line acl_parser.rl:499
te = p+1
{
                err := endArray()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    cs = (ACLParser_error)
goto _again

                }

                if !inArray() {
                    // Back to key mode
                    resetKey()
                    top--; cs = stack[top]
goto _again

                }                
            }
		case 12:
//line acl_parser.rl:516
te = p+1
{ { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 5; goto _again
 } }
		case 13:
//line acl_parser.rl:525
te = p+1
{
                if !inArray() {
                    // Not sugar, return to key mode
                    lprintf("Comma in object, return to key mode\n")
                    resetKey()
                    top--; cs = stack[top]
goto _again

                } else {
                    lprintf("Comma in array, ignore as sugar\n")
                }
            }
		case 14:
//line acl_parser.rl:538
te = p+1
{
                lprintf("Value Whitespace '%v'\n", data[ts:te])
            }
		case 15:
//line acl_parser.rl:543
te = p+1
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a value.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 16:
//line acl_parser.rl:386
te = p
p--
{ 
                lprintf("Value Identifier %v\n", data[ts:te])
                stringValue(data[ts:te])
            }
		case 17:
//line acl_parser.rl:402
te = p
p--
{
                lprintf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 18:
//line acl_parser.rl:413
te = p
p--
{
                lprintf("Value Float %v\n", data[ts:te])
                err = floatValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing float value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 19:
//line acl_parser.rl:424
te = p
p--
{
                lprintf("Value Hex %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing hex value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 20:
//line acl_parser.rl:521
te = p
p--

		case 21:
//line acl_parser.rl:543
te = p
p--
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a value.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 22:
//line acl_parser.rl:402
p = (te) - 1
{
                lprintf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 23:
//line acl_parser.rl:543
p = (te) - 1
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a value.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 24:
//line acl_parser.rl:564
te = p+1
{
                k := data[ts:te]
                if len(k) < 3 {
                    location.Message = "Key names may not be empty"
                    location.Col = findCol(data,ts)
                    cs = (ACLParser_error)
goto _again

                } else {
                    err = appendQuotedKey(k)
                    if err != nil {
                        location.Message = fmt.Sprintf("Error parsing quoted key name: %v", err)
                        location.Col = findCol(data, ts)
                        cs = (ACLParser_error)
goto _again

                    }
                }
            }
		case 25:
//line acl_parser.rl:581
te = p+1
{
                ctxStack[len(ctxStack)-1].usesEqual = (data[ts] == '=')

                lprintf("Change to value_mode '%v' usesEqual=%v\n", data[ts:te], ctxStack[len(ctxStack)-1].usesEqual)
                { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 21; goto _again
 }
            }
		case 26:
//line acl_parser.rl:589
te = p+1
{
                startObject()
            }
		case 27:
//line acl_parser.rl:596
te = p+1
{                
                err := endObject()
                if err != nil {
                    location.Message = err.Error()
                    location.Col = findCol(data, ts)              
                    cs = (ACLParser_error)
goto _again

                }

                if inArray() {
                    { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 21; goto _again
 }
                }

                // Do have to reset the key though
                resetKey()
            }
		case 28:
//line acl_parser.rl:613
te = p+1
{
                startArray()
                { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 21; goto _again
 }
            }
		case 29:
//line acl_parser.rl:619
te = p+1
{
                location.Message = "Array scope ended while in key mode indicates a parser state error."
                location.Col = findCol(data, ts)
                cs = (ACLParser_error)
goto _again

            }
		case 30:
//line acl_parser.rl:626
te = p+1
{ { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 5; goto _again
 } }
		case 31:
//line acl_parser.rl:633
te = p+1
{
                // fmt.Printf("Other '%v'\n", data[ts:te])
            }
		case 32:
//line acl_parser.rl:641
te = p+1
{

            }
		case 33:
//line acl_parser.rl:650
te = p+1
{
                if len(ctxStack[len(ctxStack)-1].keyPath) > 0 {
                    location.Message = "Key names found without a value."
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again
                    
                }
            }
		case 34:
//line acl_parser.rl:659
te = p+1
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a key.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 35:
//line acl_parser.rl:559
te = p
p--
{ 
                lprintf("Key Identifier %v\n", data[ts:te])
                appendKey(data[ts:te])
            }
		case 36:
//line acl_parser.rl:629
te = p
p--

		case 37:
//line acl_parser.rl:659
te = p
p--
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a key.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 38:
//line acl_parser.rl:659
p = (te) - 1
{
                location.Message = fmt.Sprintf("Syntax error. Invalid character '%v' while looking for a key.", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
//line acl_parser.go:1016
		}
	}

_again:
	_acts = int(_ACLParser_to_state_actions[cs])
	_nacts = uint(_ACLParser_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _ACLParser_actions[_acts-1] {
		case 2:
//line NONE:1
ts = 0

//line acl_parser.go:1030
		}
	}

	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	if p == eof {
		if _ACLParser_eof_trans[cs] > 0 {
			_trans = int(_ACLParser_eof_trans[cs] - 1)
			goto _eof_trans
		}
	}

	_out: {}
	}

//line acl_parser.rl:669


    if cs == ACLParser_error {
        if len(location.Message) == 0 {
            location.Message = "Configuration file syntax error"
        }
        fmt.Printf("Error! %v\n", location)
        return location
    }

    return nil
}
