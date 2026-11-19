package api

import (
	"fmt"
	"reflect"
	"testing"

	"yukumo-script/core/internal/domain"
)

// filler assigns a distinct, non-zero value to every field it visits, so a
// conversion that drops, duplicates, or swaps a field cannot round trip.
type filler struct {
	t    *testing.T
	next int
}

func (f *filler) fill(v reflect.Value, path string) {
	f.t.Helper()
	switch v.Kind() {
	case reflect.String:
		v.SetString(path)
	case reflect.Int:
		f.next++
		v.SetInt(int64(f.next * 10))
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			f.next++
			v.SetBytes([]byte{byte(f.next), byte(f.next + 1), byte(f.next + 2)})
			return
		}
		v.Set(reflect.MakeSlice(v.Type(), 2, 2))
		for i := 0; i < v.Len(); i++ {
			f.fill(v.Index(i), fmt.Sprintf("%s[%d]", path, i))
		}
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		f.fill(v.Elem(), path)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !v.Field(i).CanSet() {
				f.t.Fatalf("populate: cannot set %s.%s", v.Type().Name(), field.Name)
			}
			f.fill(v.Field(i), path+"."+field.Name)
		}
	default:
		f.t.Fatalf("populate: unsupported kind %s at %s; extend the filler before adding such a field", v.Kind(), path)
	}
}

func populate(t *testing.T, ptr any) {
	t.Helper()
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		t.Fatalf("populate: expected a non-nil pointer, got %T", ptr)
	}
	(&filler{t: t}).fill(v.Elem(), v.Elem().Type().Name())
}

// assertSameFieldsByName reports fields the api model is missing, and fields
// whose value differs under the same name, so a mapping cannot be silently
// dropped or pointed at the wrong field. Fields whose types differ between the
// layers (bool vs *bool, nested model types) are left to the round-trip checks.
func assertSameFieldsByName(t *testing.T, apiVal, domainVal reflect.Value) {
	t.Helper()
	apiType, domainType := apiVal.Type(), domainVal.Type()
	if apiType.NumField() != domainType.NumField() {
		t.Errorf("%s has %d fields, %s has %d", apiType.Name(), apiType.NumField(), domainType.Name(), domainType.NumField())
	}
	for i := 0; i < apiType.NumField(); i++ {
		name := apiType.Field(i).Name
		domainField := domainVal.FieldByName(name)
		if !domainField.IsValid() {
			t.Errorf("%s.%s has no %s counterpart", apiType.Name(), name, domainType.Name())
			continue
		}
		apiField := apiVal.Field(i)
		if apiField.Type() != domainField.Type() {
			continue
		}
		if !reflect.DeepEqual(apiField.Interface(), domainField.Interface()) {
			t.Errorf("%s.%s = %v, want %v", apiType.Name(), name, apiField.Interface(), domainField.Interface())
		}
	}
}

func TestSynthParamFieldCoverage(t *testing.T) {
	var d domain.SynthParam
	populate(t, &d)

	a := toParam(d)
	assertSameFieldsByName(t, reflect.ValueOf(a), reflect.ValueOf(d))
	if back := fromParam(a); !reflect.DeepEqual(d, back) {
		t.Fatalf("domain→api→domain round trip mismatch:\n in=%+v\nout=%+v", d, back)
	}

	var ap SynthParam
	populate(t, &ap)
	d2 := fromParam(ap)
	assertSameFieldsByName(t, reflect.ValueOf(ap), reflect.ValueOf(d2))
	if back := toParam(d2); !reflect.DeepEqual(ap, back) {
		t.Fatalf("api→domain→api round trip mismatch:\n in=%+v\nout=%+v", ap, back)
	}
}

func TestCharacterFieldCoverage(t *testing.T) {
	var d domain.Character
	populate(t, &d)

	a := toCharacter(d)
	assertSameFieldsByName(t, reflect.ValueOf(a), reflect.ValueOf(d))
	if len(a.Voices) != len(d.Voices) {
		t.Fatalf("voices = %d, want %d", len(a.Voices), len(d.Voices))
	}
	for i := range a.Voices {
		assertSameFieldsByName(t, reflect.ValueOf(a.Voices[i]), reflect.ValueOf(d.Voices[i]))
		if len(a.Voices[i].Params) != len(d.Voices[i].Params) {
			t.Fatalf("voices[%d].params = %d, want %d", i, len(a.Voices[i].Params), len(d.Voices[i].Params))
		}
	}
	if back := fromCharacter(a); !reflect.DeepEqual(d, back) {
		t.Fatalf("domain→api→domain round trip mismatch:\n in=%+v\nout=%+v", d, back)
	}

	var ac Character
	populate(t, &ac)
	d2 := fromCharacter(ac)
	if back := toCharacter(d2); !reflect.DeepEqual(ac, back) {
		t.Fatalf("api→domain→api round trip mismatch:\n in=%+v\nout=%+v", ac, back)
	}
}

func TestConfigFieldCoverage(t *testing.T) {
	var d domain.Config
	populate(t, &d)

	a := toConfig(d)
	assertSameFieldsByName(t, reflect.ValueOf(a), reflect.ValueOf(d))
	assertSameFieldsByName(t, reflect.ValueOf(a.Characters), reflect.ValueOf(d.Characters))
	assertSameFieldsByName(t, reflect.ValueOf(a.Paths), reflect.ValueOf(d.Paths))
	assertSameFieldsByName(t, reflect.ValueOf(a.Synth), reflect.ValueOf(d.Synth))
	if back := fromConfig(a); !reflect.DeepEqual(d, back) {
		t.Fatalf("domain→api→domain round trip mismatch:\n in=%+v\nout=%+v", d, back)
	}
}

// Nil domain pointers mean "unset", so the api model exposes the effective
// value and saving materializes it; both halves are part of the contract.
func TestConfigEffectiveBooleansFromUnset(t *testing.T) {
	in := domain.Config{Paths: domain.PathConfig{CharacterPath: "chars", OutputPath: "out"}}

	a := toConfig(in)
	if !a.Synth.AlwaysWrite || !a.Synth.ReuseCache {
		t.Fatalf("effective defaults not applied: %+v", a)
	}

	back := fromConfig(a)
	if back.Synth.AlwaysWrite == nil || !*back.Synth.AlwaysWrite ||
		back.Synth.ReuseCache == nil || !*back.Synth.ReuseCache {
		t.Fatalf("unset pointers were not written back explicitly: %+v", back)
	}
}
