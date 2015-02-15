package jsonquery

import (
	"fmt"
	"testing"
)

func TestFill(t *testing.T) {
	omap := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6}
	qomap := map[string]interface{}{"b": "int", "c": "int", "e": 5, "f": 6}
	out := map[string]interface{}{"b": 2, "c": 3, "e": 5, "f": 6}

	var ok bool

	if qomap, ok = fillmap(qomap, omap); !ok || len(qomap) != len(out) {
		t.Errorf("%v, want %v", qomap, out)
	}

	for k := range qomap {
		if qomap[k] != out[k] {
			t.Errorf("%v, want %v", qomap, out)
		}
	}

}

func TestSelect(t *testing.T) {
	s := `{"busstop":{"bus":{"lat":10,"len":20,"name":"A"}, "bus2":{"lat":10,"len":22,"name":"B"}}}`
	a := FromString(s)
	fmt.Println(a.Select(`{"lat":"float","len":"float"}`))

	//.Select(`{"lat":"float","len":"float"}`)
}
