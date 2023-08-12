package govaluate

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
)

const (
	logicalErrorFormat    string = "Value '%v' cannot be used with the logical operator '%v', it is not a bool"
	modifierErrorFormat   string = "Value '%v' cannot be used with the modifier '%v', it is not a number"
	comparatorErrorFormat string = "Value '%v' cannot be used with the comparator '%v', it is not a number"
	ternaryErrorFormat    string = "Value '%v' cannot be used with the ternary operator '%v', it is not a bool"
	prefixErrorFormat     string = "Value '%v' cannot be used with the prefix '%v'"
)

type evaluationOperator func(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error)
type stageTypeCheck func(value interface{}) bool
type stageCombinedTypeCheck func(left interface{}, right interface{}) bool

type evaluationStage struct {
	symbol OperatorSymbol

	leftStage, rightStage *evaluationStage

	// the operation that will be used to evaluate this stage (such as adding [left] to [right] and return the result)
	operator evaluationOperator

	// ensures that both left and right values are appropriate for this stage. Returns an error if they aren't operable.
	leftTypeCheck  stageTypeCheck
	rightTypeCheck stageTypeCheck

	// if specified, will override whatever is used in "leftTypeCheck" and "rightTypeCheck".
	// primarily used for specific operators that don't care which side a given type is on, but still requires one side to be of a given type
	// (like string concat)
	typeCheck stageCombinedTypeCheck

	// regardless of which type check is used, this string format will be used as the error message for type errors
	typeErrorFormat string
}

var (
	_true  = interface{}(true)
	_false = interface{}(false)
)

func (this *evaluationStage) swapWith(other *evaluationStage) {

	temp := *other
	other.setToNonStage(*this)
	this.setToNonStage(temp)
}

func (this *evaluationStage) setToNonStage(other evaluationStage) {

	this.symbol = other.symbol
	this.operator = other.operator
	this.leftTypeCheck = other.leftTypeCheck
	this.rightTypeCheck = other.rightTypeCheck
	this.typeCheck = other.typeCheck
	this.typeErrorFormat = other.typeErrorFormat
}

func (this *evaluationStage) isShortCircuitable() bool {

	switch this.symbol {
	case AND:
		fallthrough
	case OR:
		fallthrough
	case TERNARY_TRUE:
		fallthrough
	case TERNARY_FALSE:
		fallthrough
	case COALESCE:
		return true
	}

	return false
}

func noopStageRight(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	return right, leftStage, rightStage, nil
}

func addStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	// string concat if either are strings
	if isString(left) || isString(right) {
		return fmt.Sprintf("%v%v", left, right), leftStage, rightStage, nil
	}

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return leftFloat64 + rightFloat64, leftStage, rightStage, nil
}
func subtractStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return leftFloat64 - rightFloat64, leftStage, rightStage, nil
}
func multiplyStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return leftFloat64 * rightFloat64, leftStage, rightStage, nil
}
func divideStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return leftFloat64 / rightFloat64, leftStage, rightStage, nil
}
func exponentStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return math.Pow(leftFloat64, rightFloat64), leftStage, rightStage, nil
}
func modulusStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return math.Mod(leftFloat64, rightFloat64), leftStage, rightStage, nil
}
func gteStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) >= right.(string)), leftStage, rightStage, nil
	}

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(leftFloat64 >= rightFloat64), leftStage, rightStage, nil
}
func gtStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) > right.(string)), leftStage, rightStage, nil
	}

	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(leftFloat64 > rightFloat64), leftStage, rightStage, nil
}
func lteStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) <= right.(string)), leftStage, rightStage, nil
	}
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(leftFloat64 <= rightFloat64), leftStage, rightStage, nil
}
func ltStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(left.(string) < right.(string)), leftStage, rightStage, nil
	}
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(leftFloat64 < rightFloat64), leftStage, rightStage, nil
}
func equalStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(reflect.DeepEqual(left.(string), right.(string))), leftStage, rightStage, nil
	}
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(reflect.DeepEqual(leftFloat64, rightFloat64)), leftStage, rightStage, nil
}
func notEqualStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if isString(left) && isString(right) {
		return boolIface(!reflect.DeepEqual(left.(string), right.(string))), leftStage, rightStage, nil
	}
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return boolIface(!reflect.DeepEqual(leftFloat64, rightFloat64)), leftStage, rightStage, nil
}
func andStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	return boolIface(left.(bool) && right.(bool)), leftStage, rightStage, nil
}
func orStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	return boolIface(left.(bool) || right.(bool)), leftStage, rightStage, nil
}
func negateStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return -rightFloat64, leftStage, rightStage, nil
}
func invertStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	return boolIface(!right.(bool)), leftStage, rightStage, nil
}
func bitwiseNotStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(^int64(rightFloat64)), leftStage, rightStage, nil
}
func ternaryIfStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if left.(bool) {
		return right, leftStage, rightStage, nil
	}
	return nil, leftStage, rightStage, nil
}
func ternaryElseStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	if left != nil {
		return left, leftStage, rightStage, nil
	}
	return right, leftStage, rightStage, nil
}

func regexStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	var pattern *regexp.Regexp
	var err error

	switch right.(type) {
	case string:
		pattern, err = regexp.Compile(right.(string))
		if err != nil {
			return nil, leftStage, rightStage, errors.New(fmt.Sprintf("Unable to compile regexp pattern '%v': %v", right, err))
		}
	case *regexp.Regexp:
		pattern = right.(*regexp.Regexp)
	}

	return pattern.Match([]byte(left.(string))), leftStage, rightStage, nil
}

func notRegexStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	ret, leftStage, rightStage, err := regexStage(left, right, leftStage, rightStage, parameters)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	return !(ret.(bool)), leftStage, rightStage, nil
}

func bitwiseOrStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(int64(leftFloat64) | int64(rightFloat64)), right, leftStage, nil
}
func bitwiseAndStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(int64(leftFloat64) & int64(rightFloat64)), right, leftStage, nil
}
func bitwiseXORStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(int64(leftFloat64) ^ int64(rightFloat64)), right, leftStage, nil
}
func leftShiftStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(uint64(leftFloat64) << uint64(rightFloat64)), right, leftStage, nil
}
func rightShiftStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
	leftFloat64, err := convert2Float64(left)
	if err != nil {
		return nil, leftStage, rightStage, err
	}

	rightFloat64, err := convert2Float64(right)
	if err != nil {
		return nil, leftStage, rightStage, err
	}
	return float64(uint64(leftFloat64) >> uint64(rightFloat64)), right, leftStage, nil
}

func makeParameterStage(parameterName string) evaluationOperator {

	return func(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
		value, err := parameters.Get(parameterName)
		if err != nil {
			return nil, leftStage, rightStage, err
		}

		return value, leftStage, rightStage, nil
	}
}

func makeLiteralStage(literal interface{}) evaluationOperator {
	return func(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {
		return literal, leftStage, rightStage, nil
	}
}

func makeFunctionStage(function ExpressionFunction) evaluationOperator {

	return func(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

		if right == nil {
			res, err := function()
			return res, leftStage, rightStage, err
		}

		switch right.(type) {
		case []interface{}:
			res, err := function(right.([]interface{})...)
			return res, leftStage, rightStage, err
		default:
			res, err := function(right)
			return res, leftStage, rightStage, err
		}
	}
}

func typeConvertParam(p reflect.Value, t reflect.Type) (ret reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Argument type conversion failed: failed to convert '%s' to '%s'", p.Kind().String(), t.Kind().String())
			err = errors.New(errorMsg)
			ret = p
		}
	}()

	return p.Convert(t), nil
}

func typeConvertParams(method reflect.Value, params []reflect.Value) ([]reflect.Value, error) {

	methodType := method.Type()
	numIn := methodType.NumIn()
	numParams := len(params)

	if numIn != numParams {
		if numIn > numParams {
			return nil, fmt.Errorf("Too few arguments to parameter call: got %d arguments, expected %d", len(params), numIn)
		}
		return nil, fmt.Errorf("Too many arguments to parameter call: got %d arguments, expected %d", len(params), numIn)
	}

	for i := 0; i < numIn; i++ {
		t := methodType.In(i)
		p := params[i]
		pt := p.Type()

		if t.Kind() != pt.Kind() {
			np, err := typeConvertParam(p, t)
			if err != nil {
				return nil, err
			}
			params[i] = np
		}
	}

	return params, nil
}

func makeAccessorStage(pair []string, isFunction bool) evaluationOperator {

	reconstructed := strings.Join(pair, ".")

	return func(left, right, leftStage, rightStage interface{}, parameters Parameters) (ret, leftStageRet, rightStageRet interface{}, err error) {

		var params []reflect.Value

		value, err := parameters.Get(pair[0])
		if err != nil {
			return nil, leftStage, rightStage, err
		}

		// while this library generally tries to handle panic-inducing cases on its own,
		// accessors are a sticky case which have a lot of possible ways to fail.
		// therefore every call to an accessor sets up a defer that tries to recover from panics, converting them to errors.
		defer func() {
			if r := recover(); r != nil {
				errorMsg := fmt.Sprintf("Failed to access '%s': %v", reconstructed, r.(string))
				err = errors.New(errorMsg)
				leftStageRet = leftStage
				rightStageRet = rightStage
				ret = nil
			}
		}()

		for i := 1; i < len(pair); i++ {

			coreValue := reflect.ValueOf(value)

			var corePtrVal reflect.Value

			// if this is a pointer, resolve it.
			if coreValue.Kind() == reflect.Ptr {
				corePtrVal = coreValue
				coreValue = coreValue.Elem()
			}

			if coreValue.Kind() == reflect.Struct {
				if isFunction {
					field := coreValue.FieldByName(pair[i])
					if field != (reflect.Value{}) {
						value = field.Interface()
						continue
					}
				} else {
					method := coreValue.MethodByName(pair[i])
					if method == (reflect.Value{}) {
						if corePtrVal.IsValid() {
							method = corePtrVal.MethodByName(pair[i])
						}
						if method == (reflect.Value{}) {
							return nil, leftStage, rightStage, errors.New("No method or field '" + pair[i] + "' present on parameter '" + pair[i-1] + "'")
						}
					}

					switch right.(type) {
					case []interface{}:

						givenParams := right.([]interface{})
						params = make([]reflect.Value, len(givenParams))
						for idx, _ := range givenParams {
							params[idx] = reflect.ValueOf(givenParams[idx])
						}

					default:

						if right == nil {
							params = []reflect.Value{}
							break
						}

						params = []reflect.Value{reflect.ValueOf(right.(interface{}))}
					}

					params, err = typeConvertParams(method, params)

					if err != nil {
						return nil, leftStage, rightStage, errors.New("Method call failed - '" + pair[0] + "." + pair[1] + "': " + err.Error())
					}

					returned := method.Call(params)
					retLength := len(returned)

					if retLength == 0 {
						return nil, leftStage, rightStage, errors.New("Method call '" + pair[i-1] + "." + pair[i] + "' did not return any values.")
					}

					if retLength == 1 {

						value = returned[0].Interface()
						continue
					}

					if retLength == 2 {

						errIface := returned[1].Interface()
						err, validType := errIface.(error)

						if validType && errIface != nil {
							return returned[0].Interface(), leftStage, rightStage, err
						}

						value = returned[0].Interface()
						continue
					}
				}
			} else if coreValue.Kind() == reflect.Map {
				var key = reflect.ValueOf(pair[i])
				valueValue := coreValue.MapIndex(key)
				if !valueValue.IsValid() {
					return nil, leftStage, rightStage, errors.New("No field '" + pair[i] + "' present on parameter '" + pair[i-1] + "'")
				}
				value = valueValue.Interface()
				continue
			}

			// return nil, leftStage, rightStage, errors.New("Method call '" + pair[0] + "." + pair[1] + "' did not return either one value, or a value and an error. Cannot interpret meaning.")

			return nil, leftStage, rightStage, errors.New("Unable to access '" + pair[i] + "', '" + pair[i-1] + "' is not a struct or map")

		}

		value = castToFloat64(value)
		return value, leftStage, rightStage, nil
	}
}

func separatorStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	var ret []interface{}
	isOk := false
	if leftStage != nil {
		retT, ok := leftStage.([]interface{})
		if ok {
			isOk = true
			ret = retT
		}
	}

	if !isOk {
		ret = append(ret, left)
	}
	ret = append(ret, right)
	//switch left.(type) {
	//case []interface{}:
	//	ret = append(left.([]interface{}), right)
	//default:
	//	ret = []interface{}{left, right}
	//}

	return ret, ret, nil, nil
}

func inStage(left, right, leftStage, rightStage interface{}, parameters Parameters) (interface{}, interface{}, interface{}, error) {

	for _, value := range right.([]interface{}) {
		if left == value {
			return true, leftStage, rightStage, nil
		}
	}
	return false, leftStage, rightStage, nil
}

func isString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	}
	return false
}

func isRegexOrString(value interface{}) bool {

	switch value.(type) {
	case string:
		return true
	case *regexp.Regexp:
		return true
	}
	return false
}

func isBool(value interface{}) bool {

	switch value.(type) {
	case bool:
		return true
	}
	return false
}

func isNumber(value interface{}) bool {
	_, err := convert2Float64(value)
	return err == nil
}

func isFloat64(value interface{}) bool {

	switch value.(type) {
	case float64:
		return true
	case float32:
		return true
	}

	return false
}

/*
	Addition usually means between numbers, but can also mean string concat.
	String concat needs one (or both) of the sides to be a string.
*/
func additionTypeCheck(left interface{}, right interface{}) bool {

	if isNumber(left) && isNumber(right) {
		return true
	}
	if !isString(left) && !isString(right) {
		return false
	}
	return true
}

/*
	Comparison can either be between numbers, or lexicographic between two strings,
	but never between the two.
*/
func comparatorTypeCheck(left interface{}, right interface{}) bool {

	if isNumber(left) && isNumber(right) {
		return true
	}
	if isString(left) && isString(right) {
		return true
	}
	return false
}

func isArray(value interface{}) bool {
	switch value.(type) {
	case []interface{}:
		return true
	}
	return false
}

/*
	Converting a boolean to an interface{} requires an allocation.
	We can use interned bools to avoid this cost.
*/
func boolIface(b bool) interface{} {
	if b {
		return _true
	}
	return _false
}
