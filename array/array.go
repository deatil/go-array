package array

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/**
 * 获取数组数据 / array struct
 *
 * @create 2022-5-3
 * @author deatil
 */
type Array struct {
	// 分隔符 / key Delim
	keyDelim string

	// 原始数据 / source data
	source any
}

// 构造函数 / New
func New(source any) *Array {
	return &Array{
		keyDelim: ".",
		source:   source,
	}
}

// 解析 JSON 数据 / parse json data
func ParseJSON(source []byte) (*Array, error) {
	var dst any
	err := json.Unmarshal(source, &dst)

	return New(dst), err
}

// ParseJSONDecoder applies a json.Decoder to a *Container.
func ParseJSONDecoder(decoder *json.Decoder) (*Array, error) {
	var dst any
	if err := decoder.Decode(&dst); err != nil {
		return nil, err
	}

	return New(dst), nil
}

// 设置 keyDelim
// with keyDelim
func (this *Array) WithKeyDelim(data string) *Array {
	this.keyDelim = data

	return this
}

// String marshals an element to a JSON formatted string.
func (this *Array) String() string {
	return string(this.ToJSON())
}

// 返回数据
// return source data
func (this *Array) Value() any {
	return this.source
}

// 返回 JSON 数据
// return JSON data
func (this *Array) ToJSON() []byte {
	if data, err := json.Marshal(this.anyDataFormat(this.source)); err == nil {
		return data
	}

	return []byte("null")
}

// BytesIndent marshals an element to a JSON []byte blob formatted with a prefix
// and indent string.
func (this *Array) ToJSONIndent(prefix, indent string) []byte {
	if this.source != nil {
		if data, err := json.MarshalIndent(this.anyDataFormat(this.source), prefix, indent); err == nil {
			return data
		}
	}

	return []byte("null")
}

// 判断是否存在
// if key in source return true or false
func (this *Array) Exists(key string) bool {
	if this.Find(key) != nil {
		return true
	}

	return false
}

// 判断是否存在
// if key in source return true or false
func Exists(source any, key string) bool {
	return New(source).Exists(key)
}

// 获取
// get key data from source with default value
func (this *Array) Get(key string, defVal ...any) any {
	data := this.Find(key)
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
	return New(source).Get(key, defVal...)
}

// 查找
// find key data from source
func (this *Array) Find(key string) any {
	return this.Sub(key).Value()
}

// 查找
// find key data from source
func Find(source any, key string) any {
	return New(source).Find(key)
}

// 搜索
// Search data with key from source
func (this *Array) Search(path ...string) any {
	return this.search(this.source, path...)
}

// 搜索
// Search data with key from source
func Search(source any, path ...string) any {
	return New(source).Search(path...)
}

// 获取数据
// get data and return Array
func (this *Array) Sub(key string) *Array {
	path := strings.Split(key, this.keyDelim)
	source := this.search(this.source, path...)

	return &Array{
		keyDelim: this.keyDelim,
		source:   source,
	}
}

// 获取数据
// get data and return Array
func Sub(source any, key string) *Array {
	return New(source).Sub(key)
}

// Children returns a slice of all children of an array element. This also works
// for objects, however, the children returned for an source will be in a random
// order and you lose the names of the returned objects this way. If the
// underlying container value isn't an array or map nil is returned.
func (this *Array) Children() []*Array {
	source := this.anyDataFormat(this.source)

	if array, ok := source.([]any); ok {
		children := make([]*Array, len(array))
		for i := 0; i < len(array); i++ {
			children[i] = &Array{
				keyDelim: this.keyDelim,
				source:   array[i],
			}
		}

		return children
	}

	if mmap, ok := source.(map[string]any); ok {
		children := make([]*Array, 0, len(mmap))
		for _, obj := range mmap {
			children = append(children, &Array{
				keyDelim: this.keyDelim,
				source:   obj,
			})
		}

		return children
	}

	return nil
}

// ChildrenMap returns a map of all the children of an source element. IF the
// underlying value isn't a source then an empty map is returned.
func (this *Array) ChildrenMap() map[string]*Array {
	source := this.anyDataFormat(this.source)

	if mmap, ok := source.(map[string]any); ok {
		children := make(map[string]*Array, len(mmap))
		for name, obj := range mmap {
			children[name] = &Array{
				keyDelim: this.keyDelim,
				source:   obj,
			}
		}

		return children
	}

	return map[string]*Array{}
}

// 设置数据
// set data with key
func (this *Array) SetKey(value any, key string) (*Array, error) {
	path := strings.Split(key, this.keyDelim)

	return this.Set(value, path...)
}

// 设置数据
// set data with path
func (this *Array) Set(value any, path ...string) (*Array, error) {
	if len(path) == 0 {
		this.source = value
		return this, nil
	}

	if this.source == nil {
		this.source = map[string]any{}
	}

	source := this.source

	for target := 0; target < len(path); target++ {
		pathSeg := path[target]

		switch typedObj := source.(type) {
		case map[string]any:
			if target == len(path)-1 {
				source = value
				typedObj[pathSeg] = source
			} else if source = typedObj[pathSeg]; source == nil {
				typedObj[pathSeg] = map[string]any{}
				source = typedObj[pathSeg]
			}
		case []any:
			if pathSeg == "-" {
				if target < 1 {
					return nil, errors.New("unable to append new array index at root of path")
				}

				if target == len(path)-1 {
					source = value
				} else {
					source = map[string]any{}
				}

				typedObj = append(typedObj, source)
				if _, err := this.Set(typedObj, path[:target]...); err != nil {
					return nil, err
				}
			} else {
				index, err := strconv.Atoi(pathSeg)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", target, pathSeg, err)
				}

				if index < 0 {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' is invalid", target, pathSeg)
				}

				if len(typedObj) <= index {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", target, pathSeg, len(typedObj))
				}

				if target == len(path)-1 {
					source = value
					typedObj[index] = source
				} else if source = typedObj[index]; source == nil {
					return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", target, pathSeg)
				}
			}
		default:
			sourceType := reflect.TypeOf(source)

			sourceValue := reflect.ValueOf(source)
			pathSegValue := reflect.ValueOf(pathSeg)

			pathSegValue, ok := this.convertTo(sourceType.Key(), pathSegValue)
			if !ok {
				return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was error", target, pathSeg)
			}

			switch {
			case sourceType.Kind() == reflect.Map:
				if target == len(path)-1 {
					valueValue := reflect.ValueOf(value)
					valueValue, ok := this.convertTo(sourceType.Elem(), valueValue)
					if !ok {
						return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was error", target, pathSeg)
					}

					sourceValue.SetMapIndex(pathSegValue, valueValue)

					source = valueValue.Interface()
				} else if source = sourceValue.MapIndex(pathSegValue).Interface(); source == nil {
					valueValue := reflect.ValueOf(map[string]any{})
					valueValue, ok := this.convertTo(sourceType.Elem(), valueValue)
					if !ok {
						return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was error", target, pathSeg)
					}

					sourceValue.SetMapIndex(pathSegValue, valueValue)

					source = valueValue.Interface()
				}
			case sourceType.Kind() == reflect.Slice:
				if pathSeg == "-" {
					if target < 1 {
						return nil, errors.New("unable to append new array index at root of path")
					}

					if target == len(path)-1 {
						source = value
					} else {
						source = map[string]any{}
					}

					valueValue := reflect.ValueOf(source)
					valueValue, ok := this.convertTo(sourceType.Elem(), valueValue)
					if !ok {
						return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was error", target, pathSeg)
					}

					sourceValue = reflect.AppendSlice(sourceValue, valueValue)

					if _, err := this.Set(sourceValue, path[:target]...); err != nil {
						return nil, err
					}
				} else {
					index, err := strconv.Atoi(pathSeg)
					if err != nil {
						return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", target, pathSeg, err)
					}

					if index < 0 {
						return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' is invalid", target, pathSeg)
					}

					if sourceValue.Len() <= index {
						return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", target, pathSeg, sourceValue.Len())
					}

					if target == len(path)-1 {
						source = value

						valueValue := reflect.ValueOf(source)
						valueValue, ok := this.convertTo(sourceValue.Index(index).Type(), valueValue)
						if !ok {
							return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was error", target, pathSeg)
						}

						sourceValue.Index(index).Set(valueValue)
					} else if source = sourceValue.Index(index).Interface(); source == nil {
						return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", target, pathSeg)
					}
				}
			default:
				return nil, errors.New("encountered value collision whilst building path")
			}
		}
	}

	return &Array{
		keyDelim: this.keyDelim,
		source:   source,
	}, nil
}

// SetIndex attempts to set a value of an array element based on an index.
func (this *Array) SetIndex(value any, index int) (*Array, error) {
	if array, ok := this.Value().([]any); ok {
		if index >= len(array) {
			return nil, errors.New("out of bounds")
		}

		array[index] = value

		return &Array{
			keyDelim: this.keyDelim,
			source:   array[index],
		}, nil
	}

	// 反射设置
	sourceType := reflect.TypeOf(this.Value())
	sourceValue := reflect.ValueOf(this.Value())

	if sourceType.Kind() == reflect.Slice {
		if index >= sourceValue.Len() {
			return nil, errors.New("out of bounds")
		}

		valueValue := reflect.ValueOf(value)
		valueValue, ok := this.convertTo(sourceValue.Index(index).Type(), valueValue)
		if !ok {
			return nil, fmt.Errorf("failed: field '%v' was error", value)
		}

		sourceValue.Index(index).Set(valueValue)

		return &Array{
			keyDelim: this.keyDelim,
			source:   sourceValue.Interface(),
		}, nil
	}

	return nil, errors.New("not an array")
}

// 删除根据 key
// delete data with key
func (this *Array) DeleteKey(key string) error {
	path := strings.Split(key, this.keyDelim)

	return this.Delete(path...)
}

// 删除更加路径
// delete data with path
func (this *Array) Delete(path ...string) error {
	if this == nil || this.source == nil {
		return errors.New("source is nil")
	}

	if len(path) == 0 {
		return errors.New("invalid search path")
	}

	source := this.source
	target := path[len(path)-1]
	if len(path) > 1 {
		source = this.Search(path[:len(path)-1]...)
	}

	if obj, ok := source.(map[string]any); ok {
		if _, ok = obj[target]; !ok {
			return errors.New("field not found")
		}

		delete(obj, target)
		return nil
	}

	if array, ok := source.([]any); ok {
		if len(path) < 2 {
			return errors.New("unable to delete array index at root of path")
		}

		index, err := strconv.Atoi(target)
		if err != nil {
			return fmt.Errorf("failed to parse array index '%v': %v", target, err)
		}

		if index >= len(array) {
			return errors.New("out of bounds")
		}
		if index < 0 {
			return errors.New("out of bounds")
		}

		array = append(array[:index], array[index+1:]...)
		this.Set(array, path[:len(path)-1]...)
		return nil
	}

	// 通用删除
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	var dst any
	dstValue := reflect.ValueOf(&dst)

	if sourceType.Kind() == reflect.Map {
		iter := sourceValue.MapRange()
		for iter.Next() {
			k := iter.Key().Interface()

			if toString(k) != target {
				dstValue.SetMapIndex(iter.Key(), iter.Value())
			}
		}

		source = dstValue.Interface()

		return nil
	}

	if sourceType.Kind() == reflect.Slice {
		if len(path) < 2 {
			return errors.New("unable to delete array index at root of path")
		}

		index, err := strconv.Atoi(target)
		if err != nil {
			return fmt.Errorf("failed to parse array index '%v': %v", target, err)
		}

		if index >= sourceValue.Len() {
			return errors.New("out of bounds")
		}
		if index < 0 {
			return errors.New("out of bounds")
		}

		dstValue = reflect.AppendSlice(sourceValue.Slice(0, index), sourceValue.Slice(index+1, sourceValue.Len()))
		this.Set(dstValue.Interface(), path[:len(path)-1]...)
		return nil
	}

	return errors.New("source is error")
}

// Flatten a array or slice into an source of key/value pairs for each
// field, where the key is the full path of the structured field in dot path
// notation matching the spec for the method Path.
func (this *Array) Flatten() (map[string]any, error) {
	return this.flatten(false)
}

// FlattenIncludeEmpty a array or slice into an source of key/value pairs
// for each field, just as Flatten, but includes empty arrays and objects, where
// the key is the full path of the structured field in dot path notation matching
// the spec for the method Path.
func (this *Array) FlattenIncludeEmpty() (map[string]any, error) {
	return this.flatten(true)
}

func (this *Array) flatten(includeEmpty bool) (map[string]any, error) {
	flattened := map[string]any{}

	source := this.anyDataFormat(this.source)

	switch t := source.(type) {
	case map[string]any:
		this.walkObject("", t, flattened, includeEmpty)
	case []any:
		this.walkArray("", t, flattened, includeEmpty)
	default:
		return nil, errors.New("not a map or slice")
	}

	return flattened, nil
}

// 搜索
// Search data with key from source
func (this *Array) search(source any, path ...string) any {
	var val any

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

// 数组
// searchMap
func (this *Array) searchMap(source map[string]any, path []string) any {
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
func (this *Array) searchIndexWithPathPrefixes(source any, path []string) any {
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
func (this *Array) searchSliceWithPathPrefixes(
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
func (this *Array) searchMapWithPathPrefixes(
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

func (this *Array) isPathShadowedInDeepMap(path []string, m map[string]any) string {
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
func (this *Array) anyDataFormat(data any) any {
	if data == nil {
		return nil
	}

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
func (this *Array) anyDataMapFormat(data any) (map[string]any, bool) {
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
func (this *Array) anyMapFormat(data any) (map[any]any, bool) {
	m := make(map[any]any)
	isMap := false

	if data == nil {
		return m, isMap
	}

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
func (this *Array) anySliceFormat(data any) ([]any, bool) {
	m := make([]any, 0)
	isSlice := false

	if data == nil {
		return m, isSlice
	}

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

func (this *Array) walkObject(path string, obj, flat map[string]any, includeEmpty bool) {
	if includeEmpty && len(obj) == 0 {
		flat[path] = struct{}{}
	}

	for elePath, value := range obj {
		if len(path) > 0 {
			elePath = path + "." + elePath
		}

		v := this.anyDataFormat(value)

		switch t := v.(type) {
		case map[string]any:
			this.walkObject(elePath, t, flat, includeEmpty)
		case []any:
			this.walkArray(elePath, t, flat, includeEmpty)
		default:
			flat[elePath] = value
		}
	}
}

func (this *Array) walkArray(path string, arr []any, flat map[string]any, includeEmpty bool) {
	if includeEmpty && len(arr) == 0 {
		flat[path] = []struct{}{}
	}

	for i, value := range arr {
		elePath := strconv.Itoa(i)
		if len(path) > 0 {
			elePath = path + "." + elePath
		}

		ele := this.anyDataFormat(value)

		switch t := ele.(type) {
		case map[string]any:
			this.walkObject(elePath, t, flat, includeEmpty)
		case []any:
			this.walkArray(elePath, t, flat, includeEmpty)
		default:
			flat[elePath] = value
		}
	}
}

func (this *Array) convertTo(typ reflect.Type, src any) (reflect.Value, bool) {
	if !reflect.ValueOf(src).CanConvert(typ) {
		return reflect.Value{}, false
	}

	return reflect.ValueOf(src).Convert(typ), true
}
