package array

import (
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

func Example() {
	Get(arrData, "b.hhTy3.666.3")
}
