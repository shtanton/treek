package main

import (
	"io"
	"encoding/json"
)

func tokenToValue(token json.Token) Value {
	switch token.(type) {
		case nil:
			return ValueNull {}
		case bool:
			return ValueBool(token.(bool))
		case float64:
			return ValueNumber(token.(float64))
		case string:
			return ValueString(token.(string))
		default:
			panic("Can't convert JSON token to value")
	}
}

func readValue(dec *json.Decoder) (value Value, empty bool) {
	if !dec.More() {
		return nil, true
	}
	t, err := dec.Token()
	if err == io.EOF {
		return nil, true
	} else if err != nil {
		panic("Invalid JSON")
	}
	switch t.(type) {
		case nil, string, float64, bool:
			v := tokenToValue(t)
			return v, false
		case json.Delim:
			switch rune(t.(json.Delim)) {
				case '[':
					var value []Value
					for dec.More() {
						v, empty := readValue(dec)
						if empty {
							break
						}
						value = append(value, v)
					}
					t, err := dec.Token()
					if err != nil {
						panic("Invalid JSON")
					}
					delim, isDelim := t.(json.Delim)
					if !isDelim || delim != ']' {
						panic("Expected ] in JSON")
					}
					v := ValueArray(value)
					return v, false
				case '{':
					value := make(map[string]Value)
					for dec.More() {
						t, _ := dec.Token()
						key, keyIsString := t.(string)
						if !keyIsString {
							panic("Invalid JSON")
						}
						v, empty := readValue(dec)
						if empty {
							panic("Invalid JSON")
						}
						value[key] = v
					}
					t, err := dec.Token()
					if err != nil {
						panic("Invalid JSON")
					}
					delim, isDelim := t.(json.Delim)
					if !isDelim || delim != '}' {
						panic("Expected } in JSON")
					}
					v := ValueMap(value)
					return v, false
				default:
					panic("Error parsing JSON")
			}
		default:
			panic("Invalid JSON token")
	}
}

func Json(r io.Reader) Value {
	dec := json.NewDecoder(r)
	value, isEmpty := readValue(dec)
	if isEmpty {
		panic("Missing JSON input")
	}
	return value
}
