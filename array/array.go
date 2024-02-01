package array

import (
	"reflect"
	"strconv"
	"strings"
)

var defaultArray = New()

/**
 * 获取数组数据 / array struct
 *
 * @create 2022-5-3
 * @author deatil
 */
type Array struct {
	// 分隔符 / key Delim
	keyDelim string
}

// 构造函数 / NewArray
func NewArray() Array {
	return Array{
		keyDelim: ".",
	}
}

// 构造函数 / New array
func New() Array {
	return NewArray()
}

// 设置 keyDelim
// with keyDelim
func (this Array) WithKeyDelim(data string) Array {
	this.keyDelim = data

	return this
}

// 判断是否存在
// if key in source return true or false
func (this Array) Exists(source any, key string) bool {
	if this.Find(source, key) != nil {
		return true
	}

	return false
}

// 判断是否存在
// if key in source return true or false
func Exists(source any, key string) bool {
	return defaultArray.Exists(source, key)
}

// 获取
// get key data from source with default value
func (this Array) Get(source any, key string, defVal ...any) any {
	data := this.Find(source, key)
	if data != nil {
		return data
	}

	if len(defVal) > 0 {
		return defVal[0]
	}

	return nil
}

// 获取
// get key data from source with default value
func Get(source any, key string, defVal ...any) any {
	return defaultArray.Get(source, key, defVal...)
}

// 查找
// find key data from source
func (this Array) Find(source any, key string) any {
	var (
		val  any
		path = strings.Split(key, this.keyDelim)
	)

	newSource, isMap := this.anyDataMapFormat(source)
	if isMap {
		// map
		val = this.searchMap(newSource, path)
		if val != nil {
			return val
		}
	}

	// 格式化
	source = this.anyDataFormat(source)

	// 索引
	val = this.searchIndexWithPathPrefixes(source, path)
	if val != nil {
		return val
	}

	return nil
}

// 查找
// find key data from source
func Find(source any, key string) any {
	return defaultArray.Find(source, key)
}

// 数组
// searchMap
func (this Array) searchMap(source map[string]any, path []string) any {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if !ok {
		return nil
	}

	if len(path) == 1 {
		return next
	}

	switch n := next.(type) {
	case map[any]any:
		return this.searchMap(toStringMap(n), path[1:])
	case map[string]any:
		return this.searchMap(n, path[1:])
	default:
		if nextMap, isMap := this.anyMapFormat(next); isMap {
			return this.searchMap(toStringMap(nextMap), path[1:])
		}
	}

	return nil
}

// 索引查询
// searchIndexWithPathPrefixes
func (this Array) searchIndexWithPathPrefixes(source any, path []string) any {
	if len(path) == 0 {
		return source
	}

	for i := len(path); i > 0; i-- {
		prefixKey := strings.Join(path[0:i], this.keyDelim)

		var val any
		switch sourceIndexable := source.(type) {
		case []any:
			val = this.searchSliceWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		case map[string]any:
			val = this.searchMapWithPathPrefixes(sourceIndexable, prefixKey, i, path)
		}

		if val != nil {
			return val
		}
	}

	return nil
}

// 切片
// searchSliceWithPathPrefixes
func (this Array) searchSliceWithPathPrefixes(
	sourceSlice []any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
	index, err := strconv.Atoi(prefixKey)
	if err != nil || len(sourceSlice) <= index {
		return nil
	}

	next := sourceSlice[index]

	if pathIndex == len(path) {
		return next
	}

	n := this.anyDataFormat(next)
	if n != nil {
		return this.searchIndexWithPathPrefixes(n, path[pathIndex:])
	}

	return nil
}

// map 数据
// searchMapWithPathPrefixes
func (this Array) searchMapWithPathPrefixes(
	sourceMap map[string]any,
	prefixKey string,
	pathIndex int,
	path []string,
) any {
	next, ok := sourceMap[prefixKey]
	if !ok {
		return nil
	}

	if pathIndex == len(path) {
		return next
	}

	n := this.anyDataFormat(next)
	if n != nil {
		return this.searchIndexWithPathPrefixes(n, path[pathIndex:])
	}

	return nil
}

func (this Array) IsPathShadowedInDeepMap(path []string, m map[string]any) string {
	var parentVal any

	for i := 1; i < len(path); i++ {
		parentVal = this.searchMap(m, path[0:i])
		if parentVal == nil {
			return ""
		}

		switch parentVal.(type) {
		case map[any]any:
			continue
		case map[string]any:
			continue
		default:
			parentValKind := reflect.TypeOf(parentVal).Kind()
			if parentValKind == reflect.Map {
				continue
			}

			return strings.Join(path[0:i], this.keyDelim)
		}
	}

	return ""
}

// any data 数据格式化
// any data format
func (this Array) anyDataFormat(data any) any {
	switch n := data.(type) {
	case map[any]any:
		return toStringMap(n)
	case map[string]any, []any:
		return n
	default:
		dataMap, isMap := this.anyMapFormat(data)
		if isMap {
			return toStringMap(dataMap)
		}

		if dataSlice, isSlice := this.anySliceFormat(data); isSlice {
			return dataSlice
		}
	}

	return nil
}

// any data map 数据格式化
// any data map format
func (this Array) anyDataMapFormat(data any) (map[string]any, bool) {
	switch n := data.(type) {
	case map[any]any:
		return toStringMap(n), true
	case map[string]any:
		return n, true
	default:
		dataMap, isMap := this.anyMapFormat(data)
		if isMap {
			return toStringMap(dataMap), true
		}
	}

	return nil, false
}

// any map 数据格式化
// any map format
func (this Array) anyMapFormat(data any) (map[any]any, bool) {
	m := make(map[any]any)
	isMap := false

	dataValue := reflect.ValueOf(data)
	for dataValue.Kind() == reflect.Pointer {
		dataValue = dataValue.Elem()
	}

	// 获取最后的数据
	newData := dataValue.Interface()

	newDataKind := reflect.TypeOf(newData).Kind()
	if newDataKind == reflect.Map {
		iter := reflect.ValueOf(newData).MapRange()
		for iter.Next() {
			k := iter.Key().Interface()
			v := iter.Value().Interface()

			m[k] = v
		}

		isMap = true
	}

	return m, isMap
}

// any slice 数据格式化
// any slice format
func (this Array) anySliceFormat(data any) ([]any, bool) {
	m := make([]any, 0)
	isSlice := false

	dataValue := reflect.ValueOf(data)
	for dataValue.Kind() == reflect.Pointer {
		dataValue = dataValue.Elem()
	}

	// 获取最后的数据
	newData := dataValue.Interface()

	newDataKind := reflect.TypeOf(newData).Kind()
	if newDataKind == reflect.Slice {
		newDataValue := reflect.ValueOf(newData)
		newDataLen := newDataValue.Len()

		for i := 0; i < newDataLen; i++ {
			v := newDataValue.Index(i).Interface()

			m = append(m, v)
		}

		isSlice = true
	}

	return m, isSlice
}
