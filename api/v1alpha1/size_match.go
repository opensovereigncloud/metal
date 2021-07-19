package v1alpha1

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/inf.v0"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/util/jsonpath"
)

const CLabelPrefix = "machine.onmetal.de/size-"

func (s *Size) GetMatchLabel() string {
	return CLabelPrefix + s.Name
}

func (s *Size) Matches(inventory *Inventory) (bool, error) {
	for _, constraint := range s.Spec.Constraints {
		jp := jsonpath.New(constraint.Path)
		// Do not return errors if data is not found
		jp.AllowMissingKeys(true)
		err := jp.Parse(normalizeJSONPath(constraint.Path))
		if err != nil {
			return false, err
		}

		data, err := jp.FindResults(&inventory.Spec)
		if err != nil {
			return false, err
		}

		dataLen := len(data)
		// If validation data is empty, return "does not match"
		if dataLen == 0 {
			return false, nil
		}
		// If data has more than 2 arrays, multiple result sets were returned
		// we do not support that case
		if dataLen > 1 {
			return false, errors.New("multiple selection results are not supported")
		}

		validationData := data[0]
		validationDataLen := len(validationData)
		// If result array is empty for some reason, return "does not match"
		if validationDataLen == 0 {
			return false, nil
		}

		var valid bool
		// If result set has only one value, validate it as a single value
		// even if it is an aggregate, since result will be the same
		if validationDataLen == 1 {
			valid, err = constraint.MatchSingleValue(&validationData[0])
		} else {
			valid, err = constraint.MatchMultipleValues(constraint.Aggregate, validationData)
		}

		if err != nil {
			return false, err
		}

		if !valid {
			return false, nil
		}
	}

	return true, nil
}

func (r *ConstraintSpec) MatchSingleValue(value *reflect.Value) (bool, error) {
	// We do not have special constraints for nil and/or empty values
	// so returning false for now if the value is nil
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return false, nil
	}

	if r.Equal != nil {
		r, err := r.Equal.Compare(value)
		if err != nil {
			return false, err
		}
		return r == 0, nil
	}
	if r.NotEqual != nil {
		r, err := r.NotEqual.Compare(value)
		if err != nil {
			return false, err
		}
		return r != 0, nil
	}

	matches := true
	q, err := valueToQuantity(value)
	if err != nil {
		return false, err
	}
	if r.GreaterThanOrEqual != nil {
		r := q.Cmp(*r.GreaterThanOrEqual)
		matches = matches && (r >= 0)
	}
	if r.GreaterThan != nil {
		r := q.Cmp(*r.GreaterThan)
		matches = matches && (r > 0)
	}
	if r.LessThanOrEqual != nil {
		r := q.Cmp(*r.LessThanOrEqual)
		matches = matches && (r <= 0)
	}
	if r.LessThan != nil {
		r := q.Cmp(*r.LessThan)
		matches = matches && (r < 0)
	}

	return matches, nil
}

func (r *ConstraintSpec) MatchMultipleValues(aggregateType AggregateType, values []reflect.Value) (bool, error) {
	if aggregateType == "" {
		for _, value := range values {
			matches, err := r.MatchSingleValue(&value)
			if err != nil {
				return false, err
			}
			if !matches {
				return false, nil
			}
		}
		return true, nil
	}

	var agg resource.Quantity
	switch aggregateType {
	case CMinAggregateType:
		min, err := valueToQuantity(&values[0])
		if err != nil {
			return false, err
		}
		for i := 1; i < len(values); i++ {
			curr, err := valueToQuantity(&values[i])
			if err != nil {
				return false, err
			}
			if min.Cmp(*curr) > 0 {
				min = curr
			}
		}
		agg = *min
	case CMaxAggregateType:
		max, err := valueToQuantity(&values[0])
		if err != nil {
			return false, err
		}
		for i := 1; i < len(values); i++ {
			curr, err := valueToQuantity(&values[i])
			if err != nil {
				return false, err
			}
			if max.Cmp(*curr) < 0 {
				max = curr
			}
		}
		agg = *max
	case CCountAggregateType:
		agg = *resource.NewScaledQuantity(int64(len(values)), 0)
	case CAverageAggregateType:
		sum := resource.NewScaledQuantity(0, 0)
		for _, value := range values {
			q, err := valueToQuantity(&value)
			if err != nil {
				return false, err
			}
			sum.Add(*q)
		}
		decVal := sum.AsDec()
		divInt := len(values)
		div := inf.NewDec(int64(divInt), 0)
		res := decVal.QuoExact(decVal, div)
		if res == nil {
			return false, errors.Errorf("quotient of %s/%s is not finite decimal", decVal.String(), div.String())
		}
		agg = resource.MustParse(res.String())
	case CSumAggregateType:
		sum := resource.NewScaledQuantity(0, 0)
		for _, value := range values {
			q, err := valueToQuantity(&value)
			if err != nil {
				return false, err
			}
			sum.Add(*q)
		}
		agg = *sum
	default:
		return false, errors.Errorf("unknown aggregate type %s", aggregateType)
	}

	if r.Equal != nil {
		return agg.Cmp(*r.Equal.Numeric) == 0, nil
	}
	if r.NotEqual != nil {
		return agg.Cmp(*r.NotEqual.Numeric) != 0, nil
	}

	matches := true
	if r.GreaterThanOrEqual != nil {
		r := agg.Cmp(*r.GreaterThanOrEqual)
		matches = matches && (r >= 0)
	}
	if r.GreaterThan != nil {
		r := agg.Cmp(*r.GreaterThan)
		matches = matches && (r > 0)
	}
	if r.LessThanOrEqual != nil {
		r := agg.Cmp(*r.LessThanOrEqual)
		matches = matches && (r <= 0)
	}
	if r.LessThan != nil {
		r := agg.Cmp(*r.LessThan)
		matches = matches && (r < 0)
	}

	return matches, nil
}

func (r *ConstraintValSpec) Compare(value *reflect.Value) (int, error) {
	if r.Literal != nil {
		s, err := valueToString(value)
		if err != nil {
			return 0, err
		}
		return strings.Compare(s, *r.Literal), nil
	}

	if r.Numeric != nil {
		q, err := valueToQuantity(value)
		if err != nil {
			return 0, err
		}
		return q.Cmp(*r.Numeric), nil
	}

	return 0, errors.New("both numeric and literal constraints are nil")
}

func valueToString(value *reflect.Value) (string, error) {
	nonPointerValue := *value
	for {
		if nonPointerValue.Kind() == reflect.Ptr {
			nonPointerValue = nonPointerValue.Elem()
		} else {
			break
		}
	}

	if nonPointerValue.Kind() != reflect.String {
		return "", errors.Errorf("unsupported kind %s for literal comparison", nonPointerValue.Kind().String())
	}
	s := nonPointerValue.Interface().(string)

	return s, nil
}

func valueToQuantity(value *reflect.Value) (*resource.Quantity, error) {
	nonPointerValue := *value
	for {
		if nonPointerValue.Kind() == reflect.Ptr {
			nonPointerValue = nonPointerValue.Elem()
		} else {
			break
		}
	}

	switch nonPointerValue.Kind() {
	case reflect.String:
		v := nonPointerValue.Interface().(string)
		q, err := resource.ParseQuantity(v)
		if err != nil {
			return nil, err
		}
		return &q, nil
	case reflect.Int:
		v := nonPointerValue.Interface().(int)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Int8:
		v := nonPointerValue.Interface().(int8)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Int16:
		v := nonPointerValue.Interface().(int16)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Int32:
		v := nonPointerValue.Interface().(int32)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Int64:
		v := nonPointerValue.Interface().(int64)
		q := resource.NewScaledQuantity(v, 0)
		return q, nil
	case reflect.Uint:
		v := nonPointerValue.Interface().(uint)
		q := resource.MustParse(fmt.Sprintf("%d", v))
		return &q, nil
	case reflect.Uint8:
		v := nonPointerValue.Interface().(uint8)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Uint16:
		v := nonPointerValue.Interface().(uint16)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Uint32:
		v := nonPointerValue.Interface().(uint32)
		q := resource.NewScaledQuantity(int64(v), 0)
		return q, nil
	case reflect.Uint64:
		v := nonPointerValue.Interface().(uint64)
		q := resource.MustParse(fmt.Sprintf("%d", v))
		return &q, nil
	case reflect.Float32:
		v := nonPointerValue.Interface().(float32)
		q := resource.MustParse(fmt.Sprintf("%f", v))
		return &q, nil
	case reflect.Float64:
		v := nonPointerValue.Interface().(float64)
		q := resource.MustParse(fmt.Sprintf("%f", v))
		return &q, nil
	case reflect.Struct:
		v, ok := nonPointerValue.Interface().(resource.Quantity)
		if !ok {
			return nil, errors.Errorf("unsupported struct type %s for numeric comparison", nonPointerValue.Type().String())
		}
		return &v, nil
	}

	return nil, errors.Errorf("unsupported kind %s for numeric comparison", nonPointerValue.Kind().String())
}
