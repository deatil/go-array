## go-array

<p align="center">
<a href="https://pkg.go.dev/github.com/deatil/go-array"><img src="https://pkg.go.dev/badge/deatil/go-array.svg" alt="Go Reference"></a>
<a href="https://codecov.io/gh/deatil/go-array" >
 <img src="https://codecov.io/gh/deatil/go-array/graph/badge.svg?token=SS2Z1IY0XL"/>
</a>
<img src="https://goreportcard.com/badge/github.com/deatil/go-array" />
</p>

### Desc

*  go-array can use `key` and get map or slice data.

[中文](README_CN.md) | English


### Download

~~~go
go get -u github.com/deatil/go-array
~~~


### Get Starting

~~~go
import "github.com/deatil/go-array/array"

arrData := map[string]any{
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
        "hh": map[int]any{
            1115: "hccccc",
            2225: "hddddd",
            3335: map[any]string{
                "qq1": "qq1ccccc",
                "qq2": "qq2ddddd",
                "qq3": "qq3fffff",
            },
        },
        "kJh21ay": map[string]any{
            "Hjk2": "fccDcc",
            "23rt": "^hgcF5c",
        },
    },
}

data := array.Get(arrData, "b.d.e")
// output: eee

data := array.Get(arrData, "b.dd.1")
// output: ddddd

data := array.Get(arrData, "b.hh.3335.qq2")
// output: qq2ddddd

data := array.Get(arrData, "b.kJh21ay.Hjk2", "defValString")
// output: fccDcc

data := array.Get(arrData, "b.kJh21ay.Hjk23333", "defValString")
// output: defValString
~~~


### LICENSE

*  The library LICENSE is `Apache2`, using the library need keep the LICENSE.


### Copyright

*  Copyright deatil(https://github.com/deatil).
