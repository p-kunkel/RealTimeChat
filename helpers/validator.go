package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

const defaultTagName = "chat"

type field struct {
	Name      string
	IsNotNull bool
	Value     interface{}
	Tags      map[tag]bool
}

type tag struct {
	Name  string
	Value string
}

type structData []field

func LoadStructData(data interface{}, loadFieldValue bool) structData {
	var (
		result = structData{}
	)

	typeOf := reflect.TypeOf(data)
	value := reflect.ValueOf(data)

	for i := 0; i < typeOf.NumField(); i++ {
		var field field
		tagValue := typeOf.Field(i).Tag.Get(defaultTagName)
		field.Name = typeOf.Field(i).Name

		if !value.Field(i).IsZero() {
			if loadFieldValue {
				field.Value = value.Field(i).Interface()
			}
			field.IsNotNull = true
		}

		if len(tagValue) != 0 {
			field.Tags = map[tag]bool{}
			underTags := strings.Split(tagValue, ";")

			for _, underTag := range underTags {
				splitUnderTag := strings.Split(underTag, ":")
				underTagName := splitUnderTag[0]
				underTagValues := strings.Split(splitUnderTag[1], ",")

				for _, underTagValue := range underTagValues {
					field.Tags[tag{Name: strings.ToLower(underTagName), Value: strings.ToLower(underTagValue)}] = true
				}
			}
			result = append(result, field)
		}
	}
	return result
}

func (structData structData) CheckRequiredField(refersTo []string) error {
	var (
		emptyField = []string{}
	)

	for _, field := range structData {
		if ok := field.Tags[tag{Name: "required", Value: "true"}]; ok {
			for _, val := range refersTo {
				if ok := field.Tags[tag{Name: "refers_to", Value: strings.ToLower(val)}]; ok && !field.IsNotNull {
					emptyField = append(emptyField, field.Name)
				}
			}
		}
	}

	if len(emptyField) > 0 {
		return fmt.Errorf("%s is reqired", strings.Join(emptyField, ", "))
	}

	return nil
}
