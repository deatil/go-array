package array

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var (
	arrData = map[string]any{
		"a": 123,
		"b": map[string]any{
			"c": "ccc",
			"d": map[string]any{
				"e": "eee",
				"f": map[string]any{
					"g": "ggg",
				},
			},
			"dd": []any{
				"ccccc",
				"ddddd",
				"fffff",
			},
			"ff": map[any]any{
				111: "fccccc",
				222: "fddddd",
				333: "dfffff",
			},
			"hhTy3": &map[int]any{
				111: "hccccc",
				222: "hddddd",
				333: map[any]string{
					"qq1": "qq1ccccc",
					"qq2": "qq2ddddd",
					"qq3": "qq3fffff",
				},
				666: []float64{
					12.3,
					32.5,
					22.56,
					789.156,
				},
			},
			"kJh21ay": map[string]any{
				"Hjk2": "fccDcc",
				"23rt": "^hgcF5c",
				"23rt5": []any{
					"adfa",
					1231,
				},
			},
		},
	}
)

func assertT(t *testing.T) func(any, string, string) {
	return func(actual any, expected string, msg string) {
		actualStr := toString(actual)
		if actualStr != expected {
			t.Errorf("Failed %s: actual: %v, expected: %v", msg, actualStr, expected)
		}
	}
}

func assertDeepEqualT(t *testing.T) func(any, any, string) {
	return func(actual any, expected any, msg string) {
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Failed %s: actual: %v, expected: %v", msg, actual, expected)
		}
	}
}

func Test_WithKeyDelim(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		index    string
		keyDelim string
		check    string
	}{
		{
			"index-1",
			"a",
			"a",
		},
		{
			"index-2",
			"-",
			"-",
		},
	}

	for _, v := range testData {
		arr := New("").WithKeyDelim(v.keyDelim)

		assert(arr.keyDelim, v.check, "WithKeyDelim fail, index "+v.index)
	}

}

func Test_Exists(t *testing.T) {
	testData := []struct {
		index string
		key   string
		check bool
	}{
		{
			"index-1",
			"a",
			true,
		},
		{
			"index-2",
			"b.dd.1",
			true,
		},
		{
			"index-3",
			"b.ff.222333",
			false,
		},
		{
			"index-4",
			"b.hhTy3.222.yu",
			false,
		},
		{
			"index-5",
			"b.hhTy3.333.qq2",
			true,
		},
	}

	for _, v := range testData {
		check := New(arrData).Exists(v.key)
		if check != v.check {
			t.Error("Exists fail, index " + v.index)
		}
	}

}

func Test_Exists_func(t *testing.T) {
	testData := []struct {
		index string
		key   string
		check bool
	}{
		{
			"index-1",
			"a",
			true,
		},
		{
			"index-2",
			"b.dd.1",
			true,
		},
		{
			"index-3",
			"b.ff.222333",
			false,
		},
		{
			"index-4",
			"b.hhTy3.222.yu",
			false,
		},
		{
			"index-5",
			"b.hhTy3.333.qq2",
			true,
		},
	}

	for _, v := range testData {
		check := Exists(arrData, v.key)
		if check != v.check {
			t.Error("Exists func fail, index " + v.index)
		}
	}

}

func Test_Get(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		def      string
		msg      string
	}{
		{
			"a",
			"123",
			"",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"",
			"Slice",
		},
		{
			"b.hhTy3.666.9999999",
			"222555",
			"222555",
			"default",
		},
	}

	for _, v := range testData {
		check := New(arrData).Get(v.key, v.def)

		assert(check, v.expected, v.msg)
	}

}

func Test_Get_func(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		def      string
		msg      string
	}{
		{
			"a",
			"123",
			"",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"",
			"Slice",
		},
		{
			"b.hhTy3.666.9999999",
			"222555",
			"222555",
			"default",
		},
	}

	for _, v := range testData {
		check := Get(arrData, v.key, v.def)

		assert(check, v.expected, v.msg)
	}

}

func Test_Find(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := New(arrData).Find(v.key)

		assert(check, v.expected, v.msg)
	}

}

func Test_Find_func(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := Find(arrData, v.key)

		assert(check, v.expected, v.msg)
	}

}

func Test_Search(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := New(arrData).Search(strings.Split(v.key, ".")...)

		assert(check, v.expected, v.msg)
	}

}

func Test_Search_func(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := Search(arrData, strings.Split(v.key, ".")...)

		assert(check, v.expected, v.msg)
	}

}

func Test_ParseJSON(t *testing.T) {
	assert := assertT(t)

	jsonParsed, err := ParseJSON([]byte(`{
		"outer":{
			"inner":{
				"value1":21,
				"value2":35
			},
			"alsoInner":{
				"value1":99,
				"array1":[
					11, 23
				]
			}
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}

	value := jsonParsed.Find("outer.inner.value1")
	expected := "21"

	assert(value, expected, "ParseJSON fail")

	value2 := jsonParsed.Find("outer.alsoInner.array1.1")
	expected2 := "23"

	assert(value2, expected2, "ParseJSON 2 fail")
}

func Test_Sub_And_Value(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := New(arrData).Sub(v.key).Value()

		assert(check, v.expected, v.msg)
	}

}

func Test_Sub_And_Value_func(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"a",
			"123",
			"map[string]any",
		},
		{
			"b.dd.1",
			"ddddd",
			"[]any",
		},
		{
			"b.ff.222",
			"fddddd",
			"map[any]any",
		},
		{
			"b.hhTy3.222",
			"hddddd",
			"&map[int]any",
		},
		{
			"b.hhTy3.333.qq2",
			"qq2ddddd",
			"map[any]string",
		},
		{
			"b.hhTy3.666.3",
			"789.156",
			"Slice",
		},
	}

	for _, v := range testData {
		check := Sub(arrData, v.key).Value()

		assert(check, v.expected, v.msg)
	}

}

func Test_Sub_And_ToJSON(t *testing.T) {
	assert := assertT(t)

	testData := []struct {
		key      string
		expected string
		msg      string
	}{
		{
			"b.dd",
			`["ccccc","ddddd","fffff"]`,
			"[]any",
		},
		{
			"b.d",
			`{"e":"eee","f":{"g":"ggg"}}`,
			"map[any]any",
		},
		{
			"b.hhTy3.333",
			`{"qq1":"qq1ccccc","qq2":"qq2ddddd","qq3":"qq3fffff"}`,
			"&map[int]any",
		},
	}

	for _, v := range testData {
		check := New(arrData).Sub(v.key).ToJSON()

		assert(check, v.expected, v.msg)
	}

}

func Test_Children(t *testing.T) {
	jsonParsed, _ := ParseJSON([]byte(`{"map":{"objectOne":{"num":1}}, "array":[ "first", "second", "third" ]}`))

	expected := []string{"first", "second", "third"}

	children := jsonParsed.Sub("array").Children()
	for i, child := range children {
		if expected[i] != child.Value().(string) {
			t.Errorf("Child unexpected: %v != %v", expected[i], child.Value().(string))
		}
	}

	mapChildren := jsonParsed.Sub("map").Children()
	for key, val := range mapChildren {
		switch key {
		case 0:
			if val := val.Sub("num").Value().(float64); val != 1 {
				t.Errorf("%v != %v", val, 1)
			}
		default:
			t.Errorf("Unexpected key: %v", key)
		}
	}
}

func Test_ChildrenMap(t *testing.T) {
	json1, _ := ParseJSON([]byte(`{
		"objectOne":{"num":1},
		"objectTwo":{"num":2},
		"objectThree":{"num":3}
	}`))

	objectMap := json1.ChildrenMap()
	if len(objectMap) != 3 {
		t.Errorf("Wrong num of elements in objectMap: %v != %v", len(objectMap), 3)
		return
	}

	for key, val := range objectMap {
		switch key {
		case "objectOne":
			if val := val.Sub("num").Value().(float64); val != 1 {
				t.Errorf("%v != %v", val, 1)
			}
		case "objectTwo":
			if val := val.Sub("num").Value().(float64); val != 2 {
				t.Errorf("%v != %v", val, 2)
			}
		case "objectThree":
			if val := val.Sub("num").Value().(float64); val != 3 {
				t.Errorf("%v != %v", val, 3)
			}
		default:
			t.Errorf("Unexpected key: %v", key)
		}
	}
}

func Test_Flatten(t *testing.T) {
	assert := assertDeepEqualT(t)

	json1, _ := ParseJSON([]byte(`{"foo":[{"bar":"1"},{"bar":"2"}]}`))

	flattenData, err := json1.Flatten()
	if err != nil {
		t.Fatal(err)
	}

	check := map[string]any{
		"foo.0.bar": "1",
		"foo.1.bar": "2",
	}

	assert(flattenData, check, "Flatten fail")
}

func Test_FlattenIncludeEmpty(t *testing.T) {
	assert := assertDeepEqualT(t)

	json1, _ := ParseJSON([]byte(`{"foo":[{"bar":"1"},{"bar":"2"},{"bar222":{}}]}`))

	flattenData, err := json1.FlattenIncludeEmpty()
	if err != nil {
		t.Fatal(err)
	}

	check := map[string]any{
		"foo.0.bar":    "1",
		"foo.1.bar":    "2",
		"foo.2.bar222": struct{}{},
	}

	assert(flattenData, check, "FlattenIncludeEmpty fail")
}

func Test_Set(t *testing.T) {
	gObj := New(nil)

	if _, err := gObj.Set([]interface{}{}, "foo"); err != nil {
		t.Fatal(err)
	}
	if _, err := gObj.Set(1, "foo", "-"); err != nil {
		t.Fatal(err)
	}
	if _, err := gObj.Set([]interface{}{}, "foo", "-", "baz"); err != nil {
		t.Fatal(err)
	}
	if _, err := gObj.Set(2, "foo", "1", "baz", "-"); err != nil {
		t.Fatal(err)
	}
	if _, err := gObj.Set(3, "foo", "1", "baz", "-"); err != nil {
		t.Fatal(err)
	}
	if _, err := gObj.Set(5, "foo", "-"); err != nil {
		t.Fatal(err)
	}

	exp := `{"foo":[1,{"baz":[2,3]},5]}`
	if act := gObj.String(); act != exp {
		t.Errorf("Unexpected value: %v != %v", act, exp)
	}
}

func Test_SetMap(t *testing.T) {
	obj := New(arrData)

	_, err := obj.Set("yyyyyyyyy", "b", "ff", "555")
	if err != nil {
		t.Fatal(err)
	}

	res := fmt.Sprintf("%v", obj.Find("b.ff"))

	check := `map[111:fccccc 222:fddddd 333:dfffff 555:yyyyyyyyy]`

	if res != check {
		t.Errorf("SetMap fail.got %v, want %v", res, check)
	}
}

func Test_Deletes(t *testing.T) {
	jsonParsed, _ := ParseJSON([]byte(`{
		"outter":{
			"inner":{
				"value1":10,
				"value2":22,
				"value3":32
			},
			"alsoInner":{
				"value1":20,
				"value2":42,
				"value3":92
			},
			"another":{
				"value1":null,
				"value2":null,
				"value3":null
			}
		}
	}`))

	if err := jsonParsed.Delete("outter", "inner", "value2"); err != nil {
		t.Error(err)
	}
	if err := jsonParsed.Delete("outter", "inner", "value4"); err == nil {
		t.Error("value4 should not have been found in outter.inner")
	}
	if err := jsonParsed.Delete("outter", "another", "value1"); err != nil {
		t.Error(err)
	}
	if err := jsonParsed.Delete("outter", "another", "value4"); err == nil {
		t.Error("value4 should not have been found in outter.another")
	}
	if err := jsonParsed.DeleteKey("outter.alsoInner.value1"); err != nil {
		t.Error(err)
	}
	if err := jsonParsed.DeleteKey("outter.alsoInner.value4"); err == nil {
		t.Error("value4 should not have been found in outter.alsoInner")
	}
	if err := jsonParsed.DeleteKey("outter.another.value2"); err != nil {
		t.Error(err)
	}
	if err := jsonParsed.Delete("outter.another.value4"); err == nil {
		t.Error("value4 should not have been found in outter.another")
	}

	expected := `{"outter":{"alsoInner":{"value2":42,"value3":92},"another":{"value3":null},"inner":{"value1":10,"value3":32}}}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from deletes: %v != %v", actual, expected)
	}
}

func Test_DeletesWithSlices(t *testing.T) {
	rawJSON := `{
		"outter":[
			{
				"foo":{
					"value1":10,
					"value2":22,
					"value3":32
				},
				"bar": [
					20,
					42,
					92
				]
			},
			{
				"baz":{
					"value1":null,
					"value2":null,
					"value3":null
				}
			}
		]
	}`

	jsonParsed, err := ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "1", "baz", "value1"); err != nil {
		t.Error(err)
	}

	expected := `{"outter":[{"bar":[20,42,92],"foo":{"value1":10,"value2":22,"value3":32}},{"baz":{"value2":null,"value3":null}}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}

	jsonParsed, err = ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "1", "baz"); err != nil {
		t.Error(err)
	}

	expected = `{"outter":[{"bar":[20,42,92],"foo":{"value1":10,"value2":22,"value3":32}},{}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}

	jsonParsed, err = ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "1"); err != nil {
		t.Error(err)
	}

	expected = `{"outter":[{"bar":[20,42,92],"foo":{"value1":10,"value2":22,"value3":32}}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}

	jsonParsed, err = ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "0", "bar", "0"); err != nil {
		t.Error(err)
	}

	expected = `{"outter":[{"bar":[42,92],"foo":{"value1":10,"value2":22,"value3":32}},{"baz":{"value1":null,"value2":null,"value3":null}}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}

	jsonParsed, err = ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "0", "bar", "1"); err != nil {
		t.Error(err)
	}

	expected = `{"outter":[{"bar":[20,92],"foo":{"value1":10,"value2":22,"value3":32}},{"baz":{"value1":null,"value2":null,"value3":null}}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}

	jsonParsed, err = ParseJSON([]byte(rawJSON))
	if err != nil {
		t.Fatal(err)
	}
	if err = jsonParsed.Delete("outter", "0", "bar", "2"); err != nil {
		t.Error(err)
	}

	expected = `{"outter":[{"bar":[20,42],"foo":{"value1":10,"value2":22,"value3":32}},{"baz":{"value1":null,"value2":null,"value3":null}}]}`
	if actual := jsonParsed.String(); actual != expected {
		t.Errorf("Unexpected result from array deletes: %v != %v", actual, expected)
	}
}

func Example() {
	Get(arrData, "b.hhTy3.666.3")
}
