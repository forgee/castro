package lua

import (
	"reflect"
	"strconv"

	"github.com/yuin/gopher-lua"
)

// GetStructVariables loads all the global variables
// from a lua file into a struct using reflect
func GetStructVariables(src interface{}, L *lua.LState) error {
	v := reflect.ValueOf(src).Elem()

	// Lopp all struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldTag := v.Type().Field(i)

		// If field contains the tag lua
		if t, ok := fieldTag.Tag.Lookup("lua"); ok {
			if t == "" {
				continue
			}

			// Get variable from the lua stack
			variable := L.GetGlobal(t)
			if variable.Type() == lua.LTNil {
				continue
			}

			// Determine what type of variable is and
			// set the field
			switch variable.Type() {

			// Variable is integer
			case lua.LTNumber:
				n, err := strconv.ParseInt(variable.String(), 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(n)

			// Variable is boolean
			case lua.LTBool:
				field.SetBool(lua.LVAsBool(variable))

			// Variable is string
			case lua.LTString:
				field.SetString(variable.String())
			}

		}
	}
	return nil
}
