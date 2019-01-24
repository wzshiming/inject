package inject

import (
	"reflect"
	"testing"
	"time"
)

var inj = NewInjector(nil)

func TestInject(t *testing.T) {
	inj := inj.Child()

	type args struct {
		In  interface{}
		Out interface{}
	}

	type I int
	type T struct {
		D string
	}
	data := []args{
		{10, new(int)},
		{"string", new(string)},
		{I(10), new(I)},
		{T{"data"}, new(T)},
		{time.Now(), new(time.Time)},
	}

	for _, row := range data {
		in := reflect.ValueOf(row.In)
		out := reflect.ValueOf(row.Out)
		err := inj.Map(in)
		if err != nil {
			t.Error(err)
		}
		err = inj.Inject(out)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(reflect.Indirect(in).Interface(), reflect.Indirect(out).Interface()) {
			t.Errorf("Error: Injection failure: %s, %s ", in.String(), out.String())
		}
	}
}

func TestInjectStruct(t *testing.T) {
	inj := inj.Child()

	type args struct {
		Ins []interface{}
		Out interface{}
	}

	type T struct {
		S    string    `inject:""`
		I    *int      `inject:""`
		Time time.Time `inject:""`
	}
	data := []args{
		{[]interface{}{"A", 2, time.Now()}, new(T)},
		{[]interface{}{"C", 1, &time.Time{}}, new(T)},
	}

	for _, row := range data {
		for _, in := range row.Ins {
			err := inj.Map(reflect.ValueOf(in))
			if err != nil {
				t.Error(err)
			}
		}
		out := reflect.ValueOf(row.Out)

		err := inj.InjectStruct(out)
		if err != nil {
			t.Error(err)
		}
		out = reflect.Indirect(out)
		num := out.NumField()
		for i := 0; i != num; i++ {
			field := out.Field(i)
			in := reflect.ValueOf(row.Ins[i])
			if !reflect.DeepEqual(reflect.Indirect(in).Interface(), reflect.Indirect(field).Interface()) {
				t.Errorf("Error: Injection failure: %s, %s ", in.String(), field.String())
			}
		}
	}
}

func TestCall(t *testing.T) {
	inj := inj.Child()
	inj.Map(reflect.ValueOf(10))
	inj.Map(reflect.ValueOf([]string{"A", ""}))
	inj.Call(reflect.ValueOf(func(i int, s string) {
		if i != 10 {
			t.Fail()
			return
		}
		if s != "" {
			t.Fail()
			return
		}
	}))
	inj.Call(reflect.ValueOf(func(s ...string) {
		if len(s) != 2 {
			t.Fail()
			return
		}
	}))
}
