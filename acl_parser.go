
//line acl_parser.rl:1
package config

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)


//line acl_parser.go:14
var _ACLParser_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3, 
	1, 4, 1, 5, 1, 6, 1, 7, 
	1, 8, 1, 9, 1, 10, 1, 11, 
	1, 12, 1, 13, 1, 14, 1, 15, 
	1, 16, 1, 17, 1, 18, 1, 19, 
	1, 20, 1, 21, 1, 22, 1, 23, 
	1, 24, 1, 25, 1, 26, 1, 27, 
	1, 28, 1, 29, 1, 30, 1, 31, 
	1, 32, 1, 33, 1, 34, 1, 35, 
	2, 0, 6, 2, 0, 28, 
}

var _ACLParser_key_offsets []byte = []byte{
	0, 0, 3, 4, 7, 8, 10, 13, 
	16, 17, 20, 21, 23, 29, 53, 60, 
	63, 64, 67, 68, 70, 70, 97, 100, 
	101, 104, 106, 109, 111, 114, 116, 120, 
	126, 
}

var _ACLParser_trans_keys []byte = []byte{
	10, 34, 92, 10, 10, 39, 92, 10, 
	10, 42, 10, 42, 47, 10, 34, 92, 
	10, 10, 39, 92, 10, 48, 57, 48, 
	57, 65, 70, 97, 102, 10, 33, 34, 
	35, 39, 44, 45, 47, 58, 59, 61, 
	95, 123, 125, 36, 64, 65, 90, 91, 
	96, 97, 122, 124, 126, 95, 48, 57, 
	65, 90, 97, 122, 10, 34, 92, 10, 
	10, 39, 92, 45, 42, 47, 10, 34, 
	35, 39, 43, 44, 45, 47, 48, 59, 
	91, 92, 93, 94, 96, 123, 125, 33, 
	46, 49, 57, 58, 64, 65, 122, 124, 
	126, 10, 34, 92, 10, 10, 39, 92, 
	48, 57, 46, 48, 57, 48, 57, 45, 
	48, 57, 42, 47, 46, 120, 48, 57, 
	48, 57, 65, 70, 97, 102, 95, 48, 
	57, 65, 90, 97, 122, 
}

var _ACLParser_single_lengths []byte = []byte{
	0, 3, 1, 3, 1, 2, 3, 3, 
	1, 3, 1, 0, 0, 14, 1, 3, 
	1, 3, 1, 2, 0, 17, 3, 1, 
	3, 0, 1, 0, 1, 2, 2, 0, 
	1, 
}

var _ACLParser_range_lengths []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 1, 3, 5, 3, 0, 
	0, 0, 0, 0, 0, 5, 0, 0, 
	0, 1, 1, 1, 1, 0, 1, 3, 
	3, 
}

var _ACLParser_index_offsets []byte = []byte{
	0, 0, 4, 6, 10, 12, 15, 19, 
	23, 25, 29, 31, 33, 37, 57, 62, 
	66, 68, 72, 74, 77, 78, 101, 105, 
	107, 111, 113, 116, 118, 121, 124, 128, 
	132, 
}

var _ACLParser_trans_targs []byte = []byte{
	1, 13, 2, 1, 1, 1, 3, 13, 
	4, 3, 3, 3, 5, 6, 5, 5, 
	6, 20, 5, 7, 21, 8, 7, 7, 
	7, 9, 21, 10, 9, 9, 9, 27, 
	21, 31, 31, 31, 21, 13, 14, 15, 
	16, 17, 13, 18, 19, 13, 13, 13, 
	14, 13, 13, 13, 14, 13, 14, 13, 
	13, 14, 14, 14, 14, 13, 1, 13, 
	2, 1, 13, 16, 3, 13, 4, 3, 
	16, 13, 13, 16, 13, 0, 21, 22, 
	23, 24, 25, 21, 28, 29, 30, 21, 
	21, 21, 21, 21, 21, 21, 21, 21, 
	26, 21, 32, 21, 21, 7, 21, 8, 
	7, 21, 23, 9, 21, 10, 9, 26, 
	21, 11, 26, 21, 27, 21, 23, 26, 
	21, 21, 23, 21, 11, 12, 26, 21, 
	31, 31, 31, 21, 32, 32, 32, 32, 
	21, 13, 13, 13, 13, 21, 21, 21, 
	21, 21, 21, 13, 13, 13, 13, 13, 
	13, 21, 21, 21, 21, 21, 21, 21, 
	21, 21, 21, 21, 
}

var _ACLParser_trans_actions []byte = []byte{
	1, 47, 0, 0, 1, 0, 1, 47, 
	0, 0, 1, 0, 1, 0, 0, 1, 
	0, 3, 0, 1, 11, 0, 0, 1, 
	0, 1, 11, 0, 0, 1, 0, 0, 
	43, 0, 0, 0, 43, 76, 0, 9, 
	0, 9, 61, 0, 0, 49, 59, 49, 
	0, 51, 53, 63, 0, 63, 0, 63, 
	57, 0, 0, 0, 0, 65, 1, 47, 
	0, 0, 67, 0, 1, 47, 0, 0, 
	0, 69, 55, 0, 69, 0, 73, 9, 
	0, 9, 0, 21, 0, 0, 9, 13, 
	23, 29, 25, 29, 29, 15, 17, 29, 
	9, 29, 0, 29, 27, 1, 11, 0, 
	0, 39, 0, 1, 11, 0, 0, 9, 
	41, 0, 9, 33, 0, 35, 0, 9, 
	41, 19, 0, 41, 0, 0, 9, 33, 
	0, 0, 0, 37, 0, 0, 0, 0, 
	31, 71, 71, 71, 71, 45, 45, 45, 
	45, 43, 43, 65, 69, 67, 69, 69, 
	69, 41, 39, 41, 41, 33, 35, 41, 
	41, 33, 37, 31, 
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
	0, 141, 141, 141, 141, 0, 0, 145, 
	145, 145, 145, 147, 147, 0, 148, 153, 
	150, 153, 153, 153, 0, 0, 161, 155, 
	161, 161, 162, 159, 161, 161, 162, 163, 
	164, 
}

const ACLParser_start int = 13
const ACLParser_first_final int = 13
const ACLParser_error int = 0

const ACLParser_en_c_comment int = 5
const ACLParser_en_obj_value int = 21
const ACLParser_en_main int = 13


//line acl_parser.rl:13


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


    
//line acl_parser.go:366
	{
	cs = ACLParser_start
	top = 0
	ts = 0
	te = 0
	act = 0
	}

//line acl_parser.go:375
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

//line acl_parser.go:398
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
//line acl_parser.rl:224
 location.Line++ 
		case 1:
//line acl_parser.rl:228
top--; cs = stack[top]
goto _again

		case 4:
//line NONE:1
te = p+1

		case 5:
//line acl_parser.rl:262
te = p+1
{
                // fmt.Printf("Quoted literal %v\n", data[ts:te])
                err = quotedStringValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing quoted value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 6:
//line acl_parser.rl:305
te = p+1
{
                if !multiLineArray {
                    resetKey()
                    top--; cs = stack[top]
goto _again

                } // else, let it end the current value, but not the array
            }
		case 7:
//line acl_parser.rl:314
te = p+1
{
                p--
                top--; cs = stack[top]
goto _again

            }
		case 8:
//line acl_parser.rl:319
te = p+1
{
                resetKey()
                p--
                top--; cs = stack[top]
goto _again

            }
		case 9:
//line acl_parser.rl:326
te = p+1
{ { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 5; goto _again
 } }
		case 10:
//line acl_parser.rl:335
te = p+1
{
                // fmt.Printf("Comma\n")
            }
		case 11:
//line acl_parser.rl:340
te = p+1
{
                multiLineArray = true
            }
		case 12:
//line acl_parser.rl:345
te = p+1
{
                resetKey()
                top--; cs = stack[top]
goto _again

            }
		case 13:
//line acl_parser.rl:351
te = p+1
{
                // fmt.Printf("Value Whitespace '%v'\n", data[ts:te])
            }
		case 14:
//line acl_parser.rl:356
te = p+1
{
                location.Message = fmt.Sprintf("Syntax error while looking for a value: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 15:
//line acl_parser.rl:257
te = p
p--
{ 
                // fmt.Printf("Value Identifier %v\n", data[ts:te])
                stringValue(data[ts:te])
            }
		case 16:
//line acl_parser.rl:273
te = p
p--
{
                // fmt.Printf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 17:
//line acl_parser.rl:284
te = p
p--
{
                // fmt.Printf("Value Float %v\n", data[ts:te])
                err = floatValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing float value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 18:
//line acl_parser.rl:295
te = p
p--
{
                // fmt.Printf("Value Hex %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing hex value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 19:
//line acl_parser.rl:331
te = p
p--

		case 20:
//line acl_parser.rl:356
te = p
p--
{
                location.Message = fmt.Sprintf("Syntax error while looking for a value: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 21:
//line acl_parser.rl:273
p = (te) - 1
{
                // fmt.Printf("Value Integer %v\n", data[ts:te])
                err = integerValue(data[ts:te])
                if err != nil {
                    location.Message = fmt.Sprintf("Error parsing integer value: %v", err)
                    location.Col = findCol(data, ts)
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 22:
//line acl_parser.rl:356
p = (te) - 1
{
                location.Message = fmt.Sprintf("Syntax error while looking for a value: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 23:
//line acl_parser.rl:372
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
		case 24:
//line acl_parser.rl:388
te = p+1
{
                //fmt.Printf("Call obj_value '%v'\n", data[ts:te])
                multiLineArray = false
                { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 21; goto _again
 }
            }
		case 25:
//line acl_parser.rl:394
te = p+1
{
                descend()
            }
		case 26:
//line acl_parser.rl:398
te = p+1
{
                if !ascend() {
                    location.Message = "Syntax error unmatched } "
                    location.Col = findCol(data, ts)              
                    cs = (ACLParser_error)
goto _again

                }
            }
		case 27:
//line acl_parser.rl:407
te = p+1
{ { 
            if top >= len(stack)-1 {
                stack = append(stack, 0)
            }
        stack[top] = cs; top++; cs = 5; goto _again
 } }
		case 28:
//line acl_parser.rl:414
te = p+1
{
                // fmt.Printf("Other '%v'\n", data[ts:te])
            }
		case 29:
//line acl_parser.rl:418
te = p+1
{
                resetKey()
            }
		case 30:
//line acl_parser.rl:424
te = p+1
{

            }
		case 31:
//line acl_parser.rl:429
te = p+1
{
                location.Message = fmt.Sprintf("Syntax error while looking for a key name: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 32:
//line acl_parser.rl:367
te = p
p--
{ 
                //fmt.Printf("Key Identifier %v\n", data[ts:te])
                appendKey(data[ts:te])
            }
		case 33:
//line acl_parser.rl:410
te = p
p--

		case 34:
//line acl_parser.rl:429
te = p
p--
{
                location.Message = fmt.Sprintf("Syntax error while looking for a key name: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
		case 35:
//line acl_parser.rl:429
p = (te) - 1
{
                location.Message = fmt.Sprintf("Syntax error while looking for a key name: %v\n", data[ts:ts+1])
                location.Col = findCol(data, ts)              
                cs = (ACLParser_error)
goto _again

            }
//line acl_parser.go:786
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

//line acl_parser.go:800
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

//line acl_parser.rl:439


    if cs == ACLParser_error {
        if len(location.Message) == 0 {
            location.Message = "Configuration file syntax error"
        }
        fmt.Printf("Error! %v\n", location)
        return location
    }

    return nil
}
