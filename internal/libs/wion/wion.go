package wion

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	"github.com/eadgyo-forked/ion-go/ion"
)

var regValue = regexp.MustCompile(`^(.*?)(,|$)`)
var regOption = regexp.MustCompile(`^(\w+)=(.*?)$`)

var ItemSharedSymbols ion.SharedSymbolTable = ion.NewSharedSymbolTable("$ion", 1, Items_symbols_string)

func init() {
	ion.V1SystemSymbolTable = ion.NewSharedSymbolTable("$ion", 1, AllItemsSymbols())
}

type Wion struct {
	name       string
	typeWion   string
	annotation string
}

func extractWions(content string) (*Wion, error) {
	values := regValue.FindStringSubmatch(content)
	if len(values) == 0 {
		return nil, fmt.Errorf("empty")
	}

	w := Wion{
		name: values[1],
	}

	content = content[len(values[0]):]

	for len(content) != 0 {
		values = regValue.FindStringSubmatch(content)
		if len(values) == 0 {
			return nil, fmt.Errorf("weird")
		}

		options := regOption.FindStringSubmatch(values[1])
		if len(options) == 0 {
			return nil, fmt.Errorf("unrecognized option")
		}

		switch options[1] {
		case "type":
			w.typeWion = options[2]
		case "annotation":
			w.annotation = options[2]
		default:
			return nil, fmt.Errorf("unrecognized option type")
		}

		content = content[len(values[0]):]
	}

	return &w, nil
}

type WionsPerStruct struct {
	Fields  map[reflect.Value]Wion
	Wions   map[string]Wion
	Ordered []reflect.Value
}

func getFieldWionsStruct(t reflect.Value) (*WionsPerStruct, error) {
	wionsStruct := WionsPerStruct{
		Fields:  map[reflect.Value]Wion{},
		Wions:   map[string]Wion{},
		Ordered: []reflect.Value{},
	}

	for i := 0; i < t.NumField(); i++ {
		fieldwion, err := extractWions(t.Type().Field(i).Tag.Get("wion"))
		if err != nil {
			return nil, err
		}
		structField := t.Field(i)
		wionsStruct.Fields[structField] = *fieldwion
		if fieldwion.name != "" {
			wionsStruct.Wions[fieldwion.name] = *fieldwion
		}
	}

	return &wionsStruct, nil
}

func getAnnotation(name string) ion.SymbolToken {
	symbol := ItemSharedSymbols.Find(name)
	if symbol == nil {
		panic("not find symbol")
	}

	return *symbol
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func writeStruct(writer ion.Writer, vt reflect.Value) error {
	wionsStruct, err := getFieldWionsStruct(vt)
	if err != nil {
		return err
	}

	// write annotation struct
	this, exists := wionsStruct.Wions["this"]
	if exists {
		// name are designed for field not struct
		if this.name != "" && this.name != "this" {
			panic("unhandled name here!")
			//writer.FieldName(getAnnotation(this.name))
		}

		if this.annotation != "" {
			writer.Annotation(getAnnotation(this.annotation))
		}

	}
	if this.typeWion == "" {
		// begin struct
		must(writer.BeginStruct())
	} else if this.typeWion == "list" {
		must(writer.BeginList())
	}

	// for fields
	for i := 0; i < vt.NumField(); i++ {
		k := vt.Field(i)
		v := wionsStruct.Fields[k]
		writeElement(writer, k, &v)
	}

	if this.typeWion == "" {
		// end struct
		must(writer.EndStruct())
	} else if this.typeWion == "list" {
		must(writer.EndList())
	}

	return nil

}

func mustnot(s string) {
	if s != "" {
		panic("not handled here")
	}
}

func writeElement(writer ion.Writer, vt reflect.Value, wion *Wion) error {
	if wion.name == "this" {
		return nil
	}

	if wion.name != "" {
		must(writer.FieldName(getAnnotation(wion.name)))
	}

	if wion.annotation != "" {
		writer.Annotation(getAnnotation(wion.annotation))
	}

	switch vt.Kind() {
	case reflect.Bool:
		mustnot(wion.typeWion)
		return writer.WriteBool(vt.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		mustnot(wion.typeWion)
		return writer.WriteInt(vt.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		mustnot(wion.typeWion)
		return writer.WriteInt(int64(vt.Uint()))
	case reflect.Float32, reflect.Float64:
		mustnot(wion.typeWion)
		return writer.WriteFloat(vt.Float())
	case reflect.String:
		if wion.typeWion == "symbol" {
			return writer.WriteSymbolFromString(vt.String())
		}
		return writer.WriteString(vt.String())
	case reflect.Interface:
		return writeStruct(writer, vt.Elem())
	case reflect.Struct:
		return writeStruct(writer, vt)
	case reflect.Map:
		panic("unhandeld here for now!")
	case reflect.Slice, reflect.Array:
		if wion.typeWion == "sexp" {
			writer.BeginSexp()
		} else {
			writer.BeginList()
		}

		for i := 0; i < vt.Len(); i++ {
			writeElement(writer, vt.Index(i), &Wion{})
		}

		if wion.typeWion == "sexp" {
			writer.EndSexp()
		} else {
			writer.EndList()
		}

	case reflect.Pointer:
		panic("unhandeld here for now!")
	default:
		panic("not recognized type!")
	}

	return nil
}

func Marshal(v any) ([]byte, error) {

	str := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&str)
	if err := writeStruct(writer, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	if err := writer.Finish(); err != nil {
		return nil, err
	}
	return str.Bytes(), nil
}

func MarshalString(v any) (string, error) {
	str := bytes.Buffer{}
	writer := ion.NewTextWriter(&str)
	if err := writeStruct(writer, reflect.ValueOf(v)); err != nil {
		return "", err
	}
	if err := writer.Finish(); err != nil {
		return "", err
	}
	return str.String(), nil
}

// func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
// 	// If we have a non-pointer value whose type implements
// 	// Marshaler with a value receiver, then we're better off taking
// 	// the address of the value - otherwise we end up with an
// 	// allocation as we cast the value to an interface.
// 	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(marshalerType) {
// 		return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
// 	}
// 	if t.Implements(marshalerType) {
// 		return marshalerEncoder
// 	}
// 	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(textMarshalerType) {
// 		return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
// 	}
// 	if t.Implements(textMarshalerType) {
// 		return textMarshalerEncoder
// 	}

// 	switch t.Kind() {
// 	case reflect.Bool:
// 		return boolEncoder
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return intEncoder
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		return uintEncoder
// 	case reflect.Float32:
// 		return float32Encoder
// 	case reflect.Float64:
// 		return float64Encoder
// 	case reflect.String:
// 		return stringEncoder
// 	case reflect.Interface:
// 		return interfaceEncoder
// 	case reflect.Struct:
// 		return newStructEncoder(t)
// 	case reflect.Map:
// 		return newMapEncoder(t)
// 	case reflect.Slice:
// 		return newSliceEncoder(t)
// 	case reflect.Array:
// 		return newArrayEncoder(t)
// 	case reflect.Pointer:
// 		return newPtrEncoder(t)
// 	default:
// 		return unsupportedTypeEncoder
// 	}
// }
