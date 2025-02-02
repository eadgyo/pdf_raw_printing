package ionreader

import (
	"fmt"

	"github.com/amazon-ion/ion-go/ion"
)

func ReadDouble(ion1 []byte, ion2 []byte) error {
	reader1 := ion.NewReaderBytes(ion1)
	reader2 := ion.NewReaderBytes(ion2)

	return readerDouble(reader1, reader2)
}

func readerDouble(reader1 ion.Reader, reader2 ion.Reader) error {
	for reader1.Next() {
		if !reader2.Next() {
			return fmt.Errorf("missing next for reader2")
		}

		name1, err := reader1.FieldName()
		if err != nil {
			return err
		}
		name2, err := reader2.FieldName()
		if err != nil {
			return err
		}

		if name1 != name2 {
			return fmt.Errorf("different name")
		}

		an1, err := reader1.Annotations()
		if err != nil {
			return err
		}
		an2, err := reader2.Annotations()
		if err != nil {
			return err
		}

		if len(an1) != len(an2) {
			return fmt.Errorf("missing annotations")
		}

		if len(an1) == 1 {
			if an1[0].LocalSID != an2[0].LocalSID {
				return fmt.Errorf("different annotation")
			}
		} else if len(an1) != 0 {
			panic("not handled here!")
		}

		currentType1 := reader1.Type()
		currentType2 := reader2.Type()

		if currentType1 != currentType2 {
			return fmt.Errorf("different type")
		}

		if reader1.IsNull() {
			if !reader2.IsNull() {
				return fmt.Errorf("expected nil")
			}
			continue
		}

		switch currentType1 {
		case ion.SymbolType:
			a1, err := reader1.SymbolValue()
			if err != nil {
				panic(err)
			}
			a2, err := reader2.SymbolValue()
			if err != nil {
				panic(err)
			}

			if a2 != a1 {
				return fmt.Errorf("different symbol")
			}
		case ion.BoolType:
			val1, err := reader1.BoolValue()
			if err != nil {
				panic("Something went wrong while reading a Boolean value: " + err.Error())
			}
			val2, err := reader2.BoolValue()
			if err != nil {
				panic("Something went wrong while reading a Boolean value: " + err.Error())
			}
			if val1 != val2 {
				return fmt.Errorf("different bool value")
			}
		case ion.IntType:
			if name1 != nil && name1.Text != nil && *name1.Text == "page_index" {
				if *name2.Text == "page_index" {
					return fmt.Errorf("different type of int page index")
				}
			} else {
				val1, err := reader1.Int64Value()
				if err != nil {
					panic("Something went wrong while reading a Integer value: " + err.Error())
				}
				val2, err := reader2.Int64Value()
				if err != nil {
					panic("Something went wrong while reading a Integer value: " + err.Error())
				}

				if val1 != val2 {
					return fmt.Errorf("different val1 val2")
				}
			}

		case ion.StringType:
			val1, err := reader1.StringValue()
			if err != nil {
				panic("Something went wrong while reading a String value: " + err.Error())
			}
			val2, err := reader2.StringValue()
			if err != nil {
				panic("Something went wrong while reading a String value: " + err.Error())
			}
			if val1 != val2 {
				return fmt.Errorf("different string value")
			}
		case ion.FloatType:
			val1, err := reader1.FloatValue()
			if err != nil {
				panic(err)
			}
			val2, err := reader2.FloatValue()
			if err != nil {
				panic(err)
			}

			if val1 != val2 {
				return fmt.Errorf("different float value")
			}
		case ion.StructType:
			err := reader1.StepIn()
			if err != nil {
				panic(err)
			}
			err = reader2.StepIn()
			if err != nil {
				panic(err)
			}
			readerDouble(reader1, reader2)
			err = reader1.StepOut()
			if err != nil {
				panic(err)
			}

			err = reader2.StepOut()
			if err != nil {
				panic(err)
			}
		case ion.SexpType:
			err := reader1.StepIn()
			if err != nil {
				panic(err)
			}
			err = reader2.StepIn()
			if err != nil {
				panic(err)
			}

			readerDouble(reader1, reader2)

			err = reader1.StepOut()
			if err != nil {
				panic(err)
			}
			err = reader2.StepOut()
			if err != nil {
				panic(err)
			}
		case ion.ListType:
			err := reader1.StepIn()
			if err != nil {
				panic(err)
			}
			err = reader2.StepIn()
			if err != nil {
				panic(err)
			}

			readerDouble(reader1, reader2)

			err = reader1.StepOut()
			if err != nil {
				panic(err)
			}
			err = reader2.StepOut()
			if err != nil {
				panic(err)
			}
		default:
			panic("unhandled type")
		}
	}
	return nil
}
