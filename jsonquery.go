package jsonquery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Results []interface{}

func From(v interface{}) Results {
	return Results{v}
}

func FromString(s string) Results {
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	//TODO: store err in Results Object
	if err != nil {
		panic(err)
	}
	return From(v)
}

func FromURL(url string) Results {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return FromString(string(body))
}

func (objs Results) Select(query string) Results {
	// finds q in objs recursively
	var qobj interface{}
	err := json.Unmarshal([]byte(query), &qobj)
	if err != nil {
		panic(err)
	}
	rs := Results{}
	search([]interface{}(objs), qobj, &rs)
	return rs
}

func (objs Results) String() string {
	s, err := json.Marshal(objs)
	if err != nil {
		return "Marshal failed"
	}
	return string(s)
}

func (objs Results) Flatten() Result {
	r := make(map[string][]interface{})
	for _, o := range objs {
		if m, ok := o.(map[string]interface{}); ok {
			for k := range m {
				r[k] = append(r[k], m[k])
			}
		}
	}
	return Result(r)
}

type Result map[string][]interface{}

func (r Result) Rename(strs ...string) Result {
	for i := 0; i+1 < len(strs); i += 2 {
		r[strs[i+1]] = r[strs[i]]
		delete(r, strs[i])
	}
	return r
}

// recursively search in o for qo. Append matches to *rs in the form of filled
// in query interfaces qo.
func search(o, qo interface{}, rs *Results) {
	// if direct match with qo
	if v, ok := fill(qo, o); ok {
		*rs = append(*rs, v)
		// return if we only want top result
	}

	// array, map, value
	switch xs := o.(type) {
	case []interface{}:
		for _, x := range xs {
			search(x, qo, rs)
		}
	case map[string]interface{}:
		for _, x := range xs {
			search(x, qo, rs)
		}
	}
}

func fill(qo, o interface{}) (interface{}, bool) {

	// recursively check if qo and o are both []interface{}
	// if qo and o are both map[string]interface{}
	// check if all keys of qo are in o and recursively check the values
	// else check if types match
	// *** how do we repr type in qo? ***
	if s, ok := qo.(string); ok && s == "" {
		return o, true
	}

	switch x := o.(type) {
	case []interface{}:
		if _, ok := qo.([]interface{}); ok {
			return o, true
		}

	case map[string]interface{}:
		if qomap, ok := qo.(map[string]interface{}); ok {
			return fillmap(qomap, x)
		}

	case bool:
		if s, ok := qo.(string); ok && s == "bool" {
			return x, true
		}
		if v, ok := qo.(bool); ok && v == x {
			return x, true
		}
	case float64:
		if s, ok := qo.(string); ok && s == "float" {
			return x, true
		}

		if v, ok := qo.(float64); ok && v == x {
			return x, true
		}
	case int:
		if s, ok := qo.(string); ok {
			if s == "int" {
				return x, true
			}
			if s == "float" {
				return float64(x), true
			}
		}

		if v, ok := qo.(int); ok && v == x {
			return x, true
		}
	case string:
		if s, ok := qo.(string); ok && s == "string" || s == x {
			return x, true
		}
	default:
		if x == nil {
			return nil, true
		}
	}

	return nil, false
}

func fillmap(qomap, omap map[string]interface{}) (map[string]interface{}, bool) {
	// check is everything in qo map is in q and if the types match
	r := map[string]interface{}{}
	for k := range qomap {
		if _, haskey := omap[k]; !haskey {
			return nil, false
		}

		var typesmatch bool
		if r[k], typesmatch = fill(qomap[k], omap[k]); !typesmatch {
			return nil, false
		}

	}
	return r, true
}
