package ionreader

import (
	"strings"

	"github.com/amazon-ion/ion-go/ion"
)

func IonToString(ion1 []byte) (string, error) {
	reader := ion.NewReaderBytes(ion1)
	str := strings.Builder{}
	writer := ion.NewTextWriter(&str)

	writeFromReaderToWriter(reader, writer)

	err := writer.Finish()
	if err != nil {
		return "", err
	}
	return str.String(), err
}
func writeFromReaderToWriter(reader ion.Reader, writer ion.Writer) {
	for reader.Next() {
		name, err := reader.FieldName()
		if err != nil {
			panic(err)
		}

		if name != nil {
			err := writer.FieldName(*name)
			if err != nil {
				panic(err)
			}
		}

		an, err := reader.Annotations()
		if err != nil {
			panic(err)
		}

		if len(an) > 0 {
			err := writer.Annotations(an...)
			if err != nil {
				panic(err)
			}
		}

		currentType := reader.Type()
		if reader.IsNull() {
			err := writer.WriteNullType(currentType)
			if err != nil {
				panic(err)
			}
			continue
		}

		switch currentType {
		case ion.SymbolType:
			a, err := reader.SymbolValue()
			if err != nil {
				panic(err)
			}
			writer.WriteSymbol(*a)
		case ion.BoolType:
			val, err := reader.BoolValue()
			if err != nil {
				panic("Something went wrong while reading a Boolean value: " + err.Error())
			}
			err = writer.WriteBool(*val)
			if err != nil {
				panic("Something went wrong while writing a Boolean value: " + err.Error())
			}
		case ion.IntType:
			val, err := reader.Int64Value()
			if err != nil {
				panic("Something went wrong while reading a Integer value: " + err.Error())
			}
			err = writer.WriteInt(*val)
			if err != nil {
				panic("Something went wrong while writing a Integer value: " + err.Error())
			}

		case ion.StringType:
			val, err := reader.StringValue()
			if err != nil {
				panic("Something went wrong while reading a String value: " + err.Error())
			}
			err = writer.WriteString(*val)
			if err != nil {
				panic("Something went wrong while writing a String value: " + err.Error())
			}
		case ion.FloatType:
			val, err := reader.FloatValue()
			if err != nil {
				panic(err)
			}
			err = writer.WriteFloat(*val)
			if err != nil {
				panic(err)
			}
		case ion.StructType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginStruct()
			if err != nil {
				panic(err)
			}
			writeFromReaderToWriter(reader, writer)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndStruct()
			if err != nil {
				panic(err)
			}
		case ion.SexpType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginList()
			if err != nil {
				panic(err)
			}
			writeFromReaderToWriter(reader, writer)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndList()
			if err != nil {
				panic(err)
			}
		case ion.ListType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginList()
			if err != nil {
				panic(err)
			}
			writeFromReaderToWriter(reader, writer)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndList()
			if err != nil {
				panic(err)
			}
		default:
			panic("This is an example, only taking in Bool, String and Struct")
		}
	}

	if reader.Err() != nil {
		panic(reader.Err().Error())
	}
}

var PageIndexVar int64 = 1

func replacerFromReaderToWriter(reader ion.Reader, writer ion.Writer, equivalents map[string]string) {
	for reader.Next() {
		name, err := reader.FieldName()
		if err != nil {
			panic(err)
		}

		if name != nil {
			err := writer.FieldName(*name)
			if err != nil {
				panic(err)
			}
		}

		an, err := reader.Annotations()
		if err != nil {
			panic(err)
		}

		if len(an) > 0 {
			err := writer.Annotations(an...)
			if err != nil {
				panic(err)
			}
		}

		currentType := reader.Type()
		if reader.IsNull() {
			err := writer.WriteNullType(currentType)
			if err != nil {
				panic(err)
			}
			continue
		}

		switch currentType {
		case ion.SymbolType:
			a, err := reader.SymbolValue()
			if err != nil {
				panic(err)
			}
			writer.WriteSymbol(*a)
		case ion.BoolType:
			val, err := reader.BoolValue()
			if err != nil {
				panic("Something went wrong while reading a Boolean value: " + err.Error())
			}
			err = writer.WriteBool(*val)
			if err != nil {
				panic("Something went wrong while writing a Boolean value: " + err.Error())
			}
		case ion.IntType:
			if name != nil && name.Text != nil && *name.Text == "page_index" {
				err = writer.WriteInt(PageIndexVar)
				if err != nil {
					panic("Something went wrong while writing a Integer value: " + err.Error())
				}
			} else {
				val, err := reader.Int64Value()
				if err != nil {
					panic("Something went wrong while reading a Integer value: " + err.Error())
				}
				err = writer.WriteInt(*val)
				if err != nil {
					panic("Something went wrong while writing a Integer value: " + err.Error())
				}
			}

		case ion.StringType:
			val, err := reader.StringValue()
			if err != nil {
				panic("Something went wrong while reading a String value: " + err.Error())
			}
			if k, exists := equivalents[*val]; exists {
				err = writer.WriteString(k)
			} else {
				err = writer.WriteString(*val)
			}

			if err != nil {
				panic("Something went wrong while writing a String value: " + err.Error())
			}
		case ion.FloatType:
			val, err := reader.FloatValue()
			if err != nil {
				panic(err)
			}
			err = writer.WriteFloat(*val)
			if err != nil {
				panic(err)
			}
		case ion.StructType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginStruct()
			if err != nil {
				panic(err)
			}
			replacerFromReaderToWriter(reader, writer, equivalents)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndStruct()
			if err != nil {
				panic(err)
			}
		case ion.SexpType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginSexp()
			if err != nil {
				panic(err)
			}
			replacerFromReaderToWriter(reader, writer, equivalents)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndSexp()
			if err != nil {
				panic(err)
			}
		case ion.ListType:
			err := reader.StepIn()
			if err != nil {
				panic(err)
			}
			err = writer.BeginList()
			if err != nil {
				panic(err)
			}
			replacerFromReaderToWriter(reader, writer, equivalents)
			err = reader.StepOut()
			if err != nil {
				panic(err)
			}
			err = writer.EndList()
			if err != nil {
				panic(err)
			}
		default:
			panic("This is an example, only taking in Bool, String and Struct")
		}
	}

	if reader.Err() != nil {
		panic(reader.Err().Error())
	}
}
