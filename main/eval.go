package main

import (
	"strconv"
	"fmt"
	"math"
	"strings"
)

type ValueType int
const (
	TypeNull ValueType = iota
	TypeBool
	TypeNumber
	TypeString
	TypeArray
	TypeMap
)

type ValueNull struct {}
type ValueBool bool
type ValueNumber float64
type ValueString string
type ValueArray []Value
type ValueMap map[string]Value

type Value interface{
	StackValue
	getPath([]TreePathSegment) Value
	typ() ValueType
	clone() Value
	withAssignment([]Value, Value) Value
	castToBool() ValueBool
	castToNumber() ValueNumber
	castToString() ValueString
	castToArray() ValueArray
	castToMap() ValueMap
	add(Value) Value
	sub(Value) Value
	mul(Value) Value
	div(Value) Value
	index(Value) Value
	equals(Value) ValueBool
}

func castToType(v Value, t ValueType) Value {
	switch t {
		case TypeNull:
			return ValueNull{}
		case TypeBool:
			return v.castToBool()
		case TypeNumber:
			return v.castToNumber()
		case TypeString:
			return v.castToString()
		case TypeArray:
			return v.castToArray()
		case TypeMap:
			return v.castToMap()
		default:
			panic("Unknown value type")
	}
}

func (v ValueNull) getPath(path []TreePathSegment) Value {
	if len(path) != 0 {
		panic("Tried to index null")
	}
	return v
}
func (v ValueNull) typ() ValueType {
	return TypeNull
}
func (v ValueNull) clone() Value {
	return v
}
func (v ValueNull) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	res := make(ValueMap)
	res[string(path[0].castToString())] = ValueNull{}.withAssignment(path[1:], value)
	return res
}
func (v ValueNull) castToBool() ValueBool {
	return false
}
func (v ValueNull) castToNumber() ValueNumber {
	return 0
}
func (v ValueNull) castToString() ValueString {
	return ""
}
func (v ValueNull) castToArray() ValueArray {
	var res []Value
	return res
}
func (v ValueNull) castToMap() ValueMap {
	return make(map[string]Value)
}
func (v ValueNull) add(w Value) Value {
	return w
}
func (v ValueNull) sub(w Value) Value {
	typ := w.typ()
	if typ == TypeNull {
		return ValueNull{}
	}
	return castToType(v, typ).sub(w)
}
func (v ValueNull) mul(w Value) Value {
	typ := w.typ()
	if typ == TypeNull {
		return ValueNull{}
	}
	return castToType(v, typ).mul(w)
}
func (v ValueNull) div(w Value) Value {
	typ := w.typ()
	if typ == TypeNull {
		return ValueNull{}
	}
	return castToType(v, typ).div(w)
}
func (v ValueNull) index(w Value) Value {
	return ValueNull {}
}
func (v ValueNull) equals(w Value) ValueBool {
	typ := w.typ()
	if typ == TypeNull {
		return true
	}
	return castToType(v, typ).equals(w)
}

func (v ValueBool) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	res := make(ValueMap)
	res[string(path[0].castToString())] = ValueNull{}.withAssignment(path[1:], value)
	return res
}
func (v ValueBool) getPath(path []TreePathSegment) Value {
	if len(path) != 0 {
		panic("Tried to index bool")
	}
	return v
}
func (v ValueBool) typ() ValueType {
	return TypeBool
}
func (v ValueBool) clone() Value {
	return v
}
func (v ValueBool) castToBool() ValueBool {
	return v
}
func (v ValueBool) castToNumber() ValueNumber {
	if v {
		return 1
	} else {
		return 0
	}
}
func (v ValueBool) castToString() ValueString {
	if v {
		return "true"
	} else {
		return "false"
	}
}
func (v ValueBool) castToArray() ValueArray {
	var res []Value
	if v {
		res = append(res, ValueNull{})
	}
	return res
}
func (v ValueBool) castToMap() ValueMap {
	res := make(map[string]Value)
	if v {
		res[""] = ValueNull{}
	}
	return res
}
func (v ValueBool) add(w Value) Value {
	return v || w.castToBool()
}
func (v ValueBool) sub(w Value) Value {
	rhs := w.castToBool()
	return (v || rhs) && !(v && rhs)
}
func (v ValueBool) mul(w Value) Value {
	return v && w.castToBool()
}
func (v ValueBool) div(w Value) Value {
	rhs := w.castToBool()
	return (v && rhs) || !(v || rhs)
}
func (v ValueBool) index(w Value) Value {
	return v
}
func (v ValueBool) equals(w Value) ValueBool {
	rhs := w.castToBool()
	return v == rhs
}

func (v ValueNumber) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	res := make(ValueMap)
	res[string(path[0].castToString())] = ValueNull{}.withAssignment(path[1:], value)
	return res
}
func (v ValueNumber) getPath(path []TreePathSegment) Value {
	if len(path) != 0 {
		panic("Tried to index number")
	}
	return v
}
func (v ValueNumber) typ() ValueType {
	return TypeNumber
}
func (v ValueNumber) clone() Value {
	return v
}
func (v ValueNumber) castToBool() ValueBool {
	return v != 0
}
func (v ValueNumber) castToNumber() ValueNumber {
	return v
}
func (v ValueNumber) castToString() ValueString {
	return ValueString(strconv.FormatFloat(float64(v), 'g', 10, 64))
}
func (v ValueNumber) castToArray() ValueArray {
	res := make([]Value, int(math.Round(float64(v))))
	for i := range res {
		res[i] = ValueNull {}
	}
	return res
}
func (v ValueNumber) castToMap() ValueMap {
	res := make(map[string]Value)
	res[string(v.castToString())] = ValueNull {}
	return res
}
func (v ValueNumber) add(w Value) Value {
	return v + w.castToNumber()
}
func (v ValueNumber) sub(w Value) Value {
	return v - w.castToNumber()
}
func (v ValueNumber) mul(w Value) Value {
	return v * w.castToNumber()
}
func (v ValueNumber) div(w Value) Value {
	return v / w.castToNumber()
}
func (v ValueNumber) index(w Value) Value {
	return v
}
func (v ValueNumber) equals(w Value) ValueBool {
	rhs := w.castToNumber()
	return v == rhs
}

func (v ValueString) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	if len(path) > 1 {
		panic("Cannot index string twice")
	}
	index := int(math.Round(float64(path[0].castToNumber())))
	var builder strings.Builder
	reader := strings.NewReader(string(v))
	for i := 0; i < index; i += 1 {
		r, _, err := reader.ReadRune()
		if err != nil {
			panic("Error assigning to string index")
		}
		builder.WriteRune(r)
	}
	_, _, err := reader.ReadRune()
	if err != nil {
		panic("Error assigning to string index (2)")
	}
	builder.WriteString(string(value.castToString()))
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		builder.WriteRune(r)
	}
	return ValueString(builder.String())
}
func (v ValueString) getPath(path []TreePathSegment) Value {
	if len(path) != 0 {
		panic("Tried to index string")
	}
	return v
}
func (v ValueString) typ() ValueType {
	return TypeString
}
func (v ValueString) clone() Value {
	return v
}
func (v ValueString) castToBool() ValueBool {
	if v == "" || v == "false" {
		return false
	} else {
		return true
	}
}
func (v ValueString) castToNumber() ValueNumber {
	num, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0
	}
	return ValueNumber(num)
}
func (v ValueString) castToString() ValueString {
	return v
}
func (v ValueString) castToArray() ValueArray {
	fields := strings.Fields(string(v))
	res := make([]Value, len(fields))
	for i, field := range fields {
		res[i] = ValueString(field)
	}
	return res
}
func (v ValueString) castToMap() ValueMap {
	res := make(map[string]Value)
	res[string(v)] = ValueNull {}
	return res
}
func (v ValueString) add(w Value) Value {
	return v + w.castToString()
}
func (v ValueString) sub(w Value) Value {
	// TODO
	panic("Cannot subtract strings yet")
}
func (v ValueString) mul(w Value) Value {
	num := int(math.Round(float64(w.castToNumber())))
	var builder strings.Builder
	for i := 0; i < num; i += 1 {
		builder.WriteString(string(v))
	}
	return ValueString(builder.String())
}
func (v ValueString) div(w Value) Value {
	// TODO
	panic("Cannot divide strings yet")
}
func (v ValueString) index(w Value) Value {
	index := int(math.Round(float64(w.castToNumber())))
	// TODO: use proper strings functions here
	return ValueString(v[index])
}
func (v ValueString) equals(w Value) ValueBool {
	rhs := w.castToString()
	return v == rhs
}

func (v ValueArray) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	index := int(math.Round(float64(path[0].castToNumber())))
	res := v.clone().(ValueArray)
	res[index] = res[index].withAssignment(path[1:], value)
	return res
}
func (v ValueArray) getPath(path []TreePathSegment) Value {
	if len(path) == 0 {
		return v
	}
	switch path[0].(type) {
		case int:
			return v[path[0].(int)].getPath(path[1:])
		default:
			panic("Tried to index array with string")
	}
}
func (v ValueArray) typ() ValueType {
	return TypeArray
}
func (v ValueArray) clone() Value {
	var res []Value
	for _, el := range v {
		res = append(res, el.clone())
	}
	return ValueArray(res)
}
func (v ValueArray) castToBool() ValueBool {
	return len(v) > 0
}
func (v ValueArray) castToNumber() ValueNumber {
	return ValueNumber(len(v))
}
func (v ValueArray) castToString() ValueString {
	var builder strings.Builder
	for i, el := range v {
		if i != 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(string(el.castToString()))
	}
	return ValueString(builder.String())
}
func (v ValueArray) castToArray() ValueArray {
	return v
}
func (v ValueArray) castToMap() ValueMap {
	res := make(map[string]Value)
	for _, el := range v {
		res[string(el.castToString())] = ValueNull {}
	}
	return res
}
func (v ValueArray) add(w Value) Value {
	return append(v, w.castToArray()...)
}
func (v ValueArray) sub(w Value) Value {
	width := int(math.Round(float64(w.castToNumber())))
	if len(v) < width {
		return v
	} else {
		res := make([]Value, 2)
		res[0] = v[0:width]
		res[1] = v[width:]
		return ValueArray(res)
	}
}
func (v ValueArray) mul(w Value) Value {
	var res []Value
	target := int(math.Round(float64(w.castToNumber())))
	for i := 0; i < target; i += 1 {
		res = append(res, v...)
	}
	return ValueArray(res)
}
func(v ValueArray) div(w Value) Value {
	l := len(v)
	parts := int(math.Round(float64(w.castToNumber())))
	var res []Value
	part_width := l / parts
	remaining_els := l % parts
	progress := 0
	for i := 0; i < parts; i += 1 {
		if i < remaining_els {
			res = append(res, v[progress:progress+part_width+1])
			progress += part_width + 1
		} else {
			res = append(res, v[progress:progress+part_width])
			progress += part_width
		}
	}
	return ValueArray(res)
}
func (v ValueArray) index(w Value) Value {
	index := int(math.Round(float64(w.castToNumber())))
	return v[index]
}
func (v ValueArray) equals(w Value) ValueBool {
	rhs := w.castToArray()
	if len(v) != len(rhs) {
		return false
	}
	for i, el := range v {
		if !el.equals(rhs[i]) {
			return false
		}
	}
	return true
}

func (v ValueMap) withAssignment(path []Value, value Value) Value {
	if len(path) == 0 {
		return value
	}
	index := string(path[0].castToString())
	res := v.clone().(ValueMap)
	part, hasPart := res[index]
	if !hasPart {
		part = ValueNull {}
	}
	res[index] = part.withAssignment(path[1:], value)
	return res
}
func (v ValueMap) getPath(path []TreePathSegment) Value {
	if len(path) == 0 {
		return v
	}
	switch path[0].(type) {
		case string:
			return v[path[0].(string)].getPath(path[1:])
		default:
			panic("Tried to index map with int")
	}
}
func (v ValueMap) typ() ValueType {
	return TypeMap
}
func (v ValueMap) clone() Value {
	res := make(map[string]Value)
	for key, value := range v {
		res[key] = value.clone()
	}
	return ValueMap(res)
}
func (v ValueMap) castToBool() ValueBool {
	return len(v) > 0
}
func (v ValueMap) castToNumber() ValueNumber {
	return ValueNumber(len(v))
}
func (v ValueMap) castToString() ValueString {
	var builder strings.Builder
	first := true
	for key, val := range v {
		if !first {
			builder.WriteString(" ")
		}
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(string(val.castToString()))
	}
	return ValueString(builder.String())
}
func (v ValueMap) castToArray() ValueArray {
	var res []Value
	for key := range v {
		res = append(res, ValueString(key))
	}
	return res
}
func (v ValueMap) castToMap() ValueMap {
	return v
}
func (v ValueMap) add(w Value) Value {
	other := w.castToMap()
	res := v.clone().(ValueMap)
	for key, val := range other {
		res[key] = val
	}
	return res
}
func (v ValueMap) sub(w Value) Value {
	res := v.clone().(ValueMap)
	to_remove := w.castToArray()
	for _, key := range to_remove {
		delete(res, string(key.castToString()))
	}
	return res
}
func (v ValueMap) mul(w Value) Value {
	rhs := w.castToMap()
	res := make(map[string]Value)
	for key, val := range v {
		val2, hasRhsVal := rhs[key]
		if !hasRhsVal {
			val2 = ValueNull {}
		}
		res[key] = ValueArray {val, val2}
	}
	for key, val2 := range rhs {
		_, hasLhsVal := v[key]
		if !hasLhsVal {
			res[key] = ValueArray {ValueNull {}, val2}
		}
	}
	return ValueMap(res)
}
func (v ValueMap) div(w Value) Value {
	// TODO
	panic("Dividing a map not yet implemented")
}
func (v ValueMap) index(w Value) Value {
	index := string(w.castToString())
	res, hasValue := v[index]
	if !hasValue {
		return ValueNull {}
	}
	return res
}
func (v ValueMap) equals(w Value) ValueBool {
	rhs := w.castToMap()
	for key, lvalue := range v {
		rvalue, rhsHasValue := rhs[key]
		if !rhsHasValue || !bool(lvalue.equals(rvalue)) {
			return false
		}
	}
	for key := range rhs {
		_, lhsHasValue := v[key]
		if !lhsHasValue {
			return false
		}
	}
	return true
}

type VariableReference string
type IndexReference struct {
	parent StackValue
	index Value
}

type StackValue interface{
	toValue(*EvalState) Value
	toAddress() Address
}

func (v ValueNull) toValue(state *EvalState) Value {
	return v
}
func (v ValueNull) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v ValueBool) toValue(state *EvalState) Value {
	return v
}
func (v ValueBool) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v ValueNumber) toValue(state *EvalState) Value {
	return v
}
func (v ValueNumber) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v ValueString) toValue(state *EvalState) Value {
	return v
}
func (v ValueString) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v ValueArray) toValue(state *EvalState) Value {
	return v
}
func (v ValueArray) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v ValueMap) toValue(state *EvalState) Value {
	return v
}
func (v ValueMap) toAddress() Address {
	panic("Invalid assign to non variable")
}

func (v VariableReference) toValue(state *EvalState) Value {
	value, hasValue := state.variables[string(v)]
	if !hasValue {
		state.variables[string(v)] = ValueNull {}
		return ValueNull {}
	}
	return value.clone()
}
func (v VariableReference) toAddress() Address {
	return v
}

func (v IndexReference) toValue(state *EvalState) Value {
	return v.parent.toValue(state).index(v.index)
}
func (v IndexReference) toAddress() Address {
	return v
}

type Address interface {
	assign(*EvalState, Value)
	assignPath(*EvalState, []Value, Value)
}

func (v VariableReference) assign(state *EvalState, value Value) {
	state.variables[string(v)] = value
}
func (v VariableReference) assignPath(state *EvalState, path []Value, value Value) {
	state.variables[string(v)] = state.variables[string(v)].withAssignment(path, value)
}

func (v IndexReference) assign(state *EvalState, value Value) {
	v.assignPath(state, nil, value)
}
func (v IndexReference) assignPath(state *EvalState, path []Value, value Value) {
	v.parent.toAddress().assignPath(state, append([]Value{v.index}, path...), value)
}

type EvalState struct {
	stack []StackValue
	variables map[string]Value
	data Value
}

func (state *EvalState) push(value StackValue) {
	state.stack = append(state.stack, value)
}

func (state *EvalState) pop() StackValue {
	if len(state.stack) < 1 {
		panic("Error tried to pop empty stack")
	}
	index := len(state.stack) - 1
	value := state.stack[index]
	state.stack = state.stack[:index]
	return value
}

func (state *EvalState) popValue() Value {
	return state.pop().toValue(state)
}

func (state *EvalState) popAddress() Address {
	return state.pop().toAddress()
}

func (index PatternSegmentIndex) matches(_ *EvalState, path []TreePathSegment, pathSegment TreePathSegment) bool {
	switch pathSegment.(type) {
		case string:
			return string(index) == pathSegment.(string)
		case int:
			return string(index) == strconv.Itoa(pathSegment.(int))
		default:
			panic("Bug in treek, invalid TreePathSegment")
	}
}

func (filter PatternSegmentFilter) matches(state *EvalState, path []TreePathSegment, pathSegment TreePathSegment) bool {
	state.variables["path"] = pathToValueArray(path).clone()
	state.variables["$0"] = state.data.getPath(path).clone()
	result := evalExpr(state, Expression(filter))
	return bool(result.castToBool())
}

func (segment PatternSegmentBasic) matches(state *EvalState, path []TreePathSegment, pathSegment TreePathSegment) bool {
	switch segment {
		case PatternSegmentAll:
			return true
		default:
			panic("Invalid basic pattern segment")
	}
}

func matchPattern(state *EvalState, pattern Pattern, walkItem TreeWalkItem) bool {
	if len(pattern.segments) != len(walkItem.path)  || pattern.isFirst != walkItem.first{
		return false
	}
	for i, patternSegment := range pattern.segments {
		pathSegment := walkItem.path[i]
		if !patternSegment.matches(state, walkItem.path[0:i+1], pathSegment) {
			return false
		}
	}
	return true
}

func evalAction(state *EvalState, action Expression, node TreeWalkItem) {
	if len(action) == 0 {
		subroutinePrintln([]Value{state.data.getPath(node.path)})
		return
	}
	state.variables["path"] = pathToValueArray(node.path).clone()
	state.variables["$0"] = state.data.getPath(node.path).clone()
	evalExpr(state, action)
}

func (instruction InstructionBasic) eval(state *EvalState) {
	switch instruction {
		case InstructionAdd:
			rhs := state.popValue()
			lhs := state.popValue()
			state.push(lhs.add(rhs))
		case InstructionSub:
			rhs := state.popValue()
			lhs := state.popValue()
			state.push(lhs.sub(rhs))
		case InstructionDiv:
			rhs := state.popValue()
			lhs := state.popValue()
			state.push(lhs.div(rhs))
		case InstructionMul:
			rhs := state.popValue()
			lhs := state.popValue()
			state.push(lhs.mul(rhs))
		case InstructionIgnore:
			state.popValue()
		case InstructionPushNull:
			state.push(ValueNull {})
		case InstructionAssign:
			rhs := state.popValue()
			lhs := state.popAddress()
			lhs.assign(state, rhs)
			state.push(ValueNull {})
		case InstructionIndex:
			index := state.popValue()
			parent := state.pop()
			state.push(IndexReference {parent, index})
		case InstructionDup:
			val := state.pop()
			state.push(val)
			state.push(val.toValue(state).clone())
		case InstructionEqual:
			rhs := state.popValue()
			lhs := state.popValue()
			state.push(lhs.equals(rhs))
		case InstructionNot:
			val := state.popValue().castToBool()
			state.push(!val)
		default:
			panic("Error: Tried to execute invalid basic instruction")
	}
}

func (n InstructionPushNumber) eval(state *EvalState) {
	state.push(ValueNumber(n))
}

func (variable InstructionPushVariable) eval(state *EvalState) {
	state.push(VariableReference(variable))
}

func (s InstructionPushString) eval(state *EvalState) {
	state.push(ValueString(s))
}

type SubroutineFn func ([]Value) Value

func printSingle(arg Value) {
	switch arg.(type) {
		case ValueNull:
			fmt.Print("null")
		case ValueBool:
			fmt.Printf("%v", bool(arg.(ValueBool)))
		case ValueNumber:
			fmt.Printf("%v", float64(arg.(ValueNumber)))
		case ValueString:
			fmt.Printf("%q", string(arg.(ValueString)))
		case ValueArray:
			fmt.Print("[")
			for i, el := range arg.(ValueArray) {
				if i != 0 {
					fmt.Print(", ")
				}
				printSingle(el)
			}
			fmt.Print("]")
		case ValueMap:
			fmt.Print("{")
			isStart := true
			for key, value := range arg.(ValueMap) {
				if !isStart {
					fmt.Print(", ")
				}
				fmt.Printf("%q: ", key)
				printSingle(value)
				isStart = false
			}
			fmt.Print("}")
	}
}
func subroutinePrintln(args []Value) Value {
	for i, arg := range args {
		if i != 0 {
			fmt.Print(" ")
		}
		printSingle(arg)
	}
	fmt.Print("\n")
	return ValueNull{}
}

func (call InstructionCall) eval(state *EvalState) {
	args := make([]Value, call.nargs)
	for i := call.nargs - 1; i >= 0; i -= 1 {
		args[i] = state.popValue()
	}
	subroutines := map[Subroutine]SubroutineFn {
		SubroutinePrintln: subroutinePrintln,
	}
	subroutine, isSubroutine := subroutines[call.subroutine]
	if !isSubroutine {
		panic("Error: Invalid subroutine")
	}
	state.push(subroutine(args))
}

func evalExpr(state *EvalState, expr Expression) Value {
	for _, instruction := range expr {
		instruction.eval(state)
	}
	return state.popValue()
}

func pathToValueArray(path []TreePathSegment) ValueArray {
	value := make(ValueArray, len(path))
	for i, segment := range path {
		switch segment.(type) {
			case string:
				value[i] = ValueString(segment.(string))
			case int:
				value[i] = ValueNumber(float64(segment.(int)))
		}
	}
	return value
}

type TreeWalkItem struct {
	path []TreePathSegment
	first bool
}

func walkPaths(data Value, path []TreePathSegment, out chan TreeWalkItem) {
	out <- TreeWalkItem {path, true}
	switch data.(type) {
		case ValueNull, ValueBool, ValueNumber, ValueString:
		case ValueArray:
			for i, el := range data.(ValueArray) {
				walkPaths(el, append(path, i), out)
			}
		case ValueMap:
			for key, el := range data.(ValueMap) {
				walkPaths(el, append(path, key), out)
			}
	}
	out <- TreeWalkItem {path, false}
}
func pathsRoutine(data Value, path []TreePathSegment, out chan TreeWalkItem) {
	walkPaths(data, path, out)
	close(out)
}
func getPaths(data Value) chan TreeWalkItem {
	out := make(chan TreeWalkItem)
	go pathsRoutine(data, nil, out)
	return out
}

func Eval(program Program, data Value) {
	state := &EvalState {
		stack: nil,
		variables: make(map[string]Value),
		data: data,
	}
	paths := getPaths(data)
	for node := range paths {
		for _, block := range program.blocks {
			if matchPattern(state, block.pattern, node) {
				evalAction(state, block.action,  node)
			}
		}
	}
}
