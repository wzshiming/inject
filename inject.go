package inject

import (
	"fmt"
	"reflect"
)

// Injector Save the current mapping and parent injectors.
type Injector struct {
	parent  *Injector
	mapping map[reflect.Type]reflect.Value
}

// NewInjector Create a new injector.
func NewInjector(parent *Injector) *Injector {
	return &Injector{
		parent:  parent,
		mapping: map[reflect.Type]reflect.Value{},
	}
}

// Child returns a sub-injector and inherit the current.
func (in *Injector) Child() *Injector {
	return NewInjector(in)
}

// Map map this value.
func (in *Injector) Map(val reflect.Value) error {
	typ := val.Type()
	in.mapping[typ] = val
	switch val.Kind() {
	case reflect.Interface, reflect.Ptr:
		if val.IsNil() {
			return fmt.Errorf(`Error: Mapped a null value: %v`, val)
		}
		return in.Map(val.Elem())
	}
	return nil
}

// Get Find the value of the type from the mapping and return zero value if not.
func (in *Injector) Get(typ reflect.Type) reflect.Value {
	if data, ok := in.Lookup(typ); ok {
		return data
	}
	return reflect.Zero(typ)
}

// Lookup Find the value of the type from the mapping.
func (in *Injector) Lookup(typ reflect.Type) (reflect.Value, bool) {
	if data, ok := in.mapping[typ]; ok {
		return data, ok
	}
	if in.parent != nil {
		return in.parent.Lookup(typ)
	}
	return reflect.Value{}, false
}

// InjectStruct inject this value, structure will find the not `inject:"-"` tags for injection.
func (in *Injector) InjectStruct(val reflect.Value) error {
	return in.inject(val, true)
}

// Inject inject this value.
func (in *Injector) Inject(val reflect.Value) error {
	return in.inject(val, false)
}

func (in *Injector) inject(val reflect.Value, structure bool) error {
	typ := val.Type()
	if data, ok := in.Lookup(typ); ok {
		if !val.CanSet() {
			return nil
		}
		val.Set(data)
		return nil
	}
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			if !val.CanSet() {
				return nil
			}
			val.Set(reflect.New(typ.Elem()))
		}
		return in.inject(val.Elem(), structure)
	case reflect.Interface:
		if !val.IsNil() {
			return in.inject(val.Elem(), structure)
		}
	case reflect.Struct:
		if structure {
			num := typ.NumField()
			for i := 0; i != num; i++ {
				field := typ.Field(i)
				if v, ok := field.Tag.Lookup("inject"); ok && v == "-" {
					continue
				}
				in.inject(val.Field(i), structure)
			}
			return nil
		}
	}
	return fmt.Errorf(`Error: No values that can be mapped: %v`, val.String())
}

// Call Find the parameter injection from the mapping and call the function.
func (in *Injector) Call(fun reflect.Value) ([]reflect.Value, error) {
	switch fun.Kind() {
	default:
		return nil, fmt.Errorf(`Error: Not a function type: %v`, fun.String())
	case reflect.Interface, reflect.Ptr:
		if fun.IsNil() {
			return nil, fmt.Errorf(`Error: Is a null value: %v`, fun.String())
		}
		return in.Call(fun.Elem())
	case reflect.Func:
		typ := fun.Type()
		num := typ.NumIn()
		if num == 0 {
			return fun.Call(nil), nil
		}

		args := make([]reflect.Value, 0, num)
		for i := 0; i != num; i++ {
			args = append(args, in.Get(typ.In(i)))
		}

		if typ.IsVariadic() {
			return fun.CallSlice(args), nil
		}
		return fun.Call(args), nil
	}
}
