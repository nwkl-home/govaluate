govaluate
====

基于[govaluate](https://github.com/Knetic/govaluate)规则引擎做一系列优化，以满足需要, 感谢作者提供这么优秀的项目

主要优化一下几个问题：
--

问题一: 通过函数输出数值类型，然后在做数学运算，就会报错

```go
    // 定义函数
    func returnInt(arg ...interface{}) (interface{}, error) {
        return 8, nil
    }

    func main() {
        expressionFunctions := make(map[string]govaluate.ExpressionFunction)
        // 返回int类型的函数
        expressionFunctions["returnInt"] = returnInt
        // returnInt() + 7 执行函数返回值 + 7
        expr, _ := govaluate.NewEvaluableExpressionWithFunctions("returnInt() + 7", expressionFunctions)
        parameters := make(map[string]interface{}, 0)
        fmt.Println(expr.Evaluate(parameters))
    }
	// 期望：15， 输出: <nil> Value '8' cannot be used with the modifier '+', it is not a number
```

优化：通过函数返回值进行计算的，需要兼容

问题二: 函数如果函数实参是[]interface{}类型，则实际参数是切片的内容，而非整个切片

```go
    func exportParamNum(args ...interface{}) (interface{}, error) {
		fmt.Println(args)
		return len(args), nil
	}

    func main() {
		expressionFunctions := make(map[string]govaluate.ExpressionFunction) 
		// 返回参数个数
		expressionFunctions["exportParamNum"] = exportParamNum
        // 增加一个长度为2的切片和一个字符串，期望返回参数个数是2
        expr, _ := govaluate.NewEvaluableExpressionWithFunctions("exportParamNum(a, b)", expressionFunctions)
        parameters := make(map[string]interface{}, 0)
        parameters["a"] = []interface{}{1, "a"}
        parameters["b"] = "b"
        fmt.Println(expr.Evaluate(parameters))
    }
	// 期望：[[1 a] b] 个数是2，输出: [1 a b]  3 <nil>
```
优化：传参发生不能改变，传什么就是什么

问题三: govaluate本身支持访问器访问参数中的某个结构体，但是前提要定义一个结构体，可是在实际使用中，基本没有定义结构体，大多是场景都是使用map，这样就很麻烦了，不能直接获取，需要定义获取的函数

```go
    func main() {
        parameters := make(map[string]interface{}, 0)
        a := make(map[string]interface{}, 0)
        a["b"] = 1
        parameters["a"] = a
        // 获取map里面的某个键值
        expr, err := govaluate.NewEvaluableExpression("a.b")
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println(expr.Evaluate(parameters))
	}
	//  期望：1 输出：Unable to access unexported field 'b' in token 'a.b'
```

优化：优化map取值方式，跟结构体取值方式一样，区别是像结构体一样定义方法

联系本人
--

欢迎探讨：+V nuanwei1314 
