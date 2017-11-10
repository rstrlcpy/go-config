// Copyright 2017 Roman Strashkin.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func BuildConfig(config interface{}) {

	configValue := reflect.ValueOf(config)
	if configValue.Kind() != reflect.Ptr {
		panic("input argument is not a pointer")
	}

	if configValue.IsNil() {
		panic("input argument is a nil pointer")
	}

	configType := configValue.Elem().Type()
	if configType.Kind() != reflect.Struct {
		panic("input argument should be a poiner to a struct")
	}

	for i := 0; i < configType.NumField(); i++ {
		structField := configType.Field(i)
		tagValue, ok := structField.Tag.Lookup("config")
		if !ok {
			continue
		}

		tagFields := strings.Split(tagValue, ",")
		if len(tagFields) == 0 {
			panic(fmt.Sprintf("Field '%s': empty tag", structField.Name))
		}

		envVarName := tagFields[0]
		required := false
		if len(tagFields) > 1 {
			if tagFields[1] != "required" {
				panic(fmt.Sprintf("Field '%s': unsupported tag's field",
					tagFields[1]))
			}

			required = true
		}

		fValue := configValue.Elem().FieldByName(structField.Name)
		switch structField.Type.Kind() {
		case reflect.String:
			fValue.SetString(retrieveEnvVar(envVarName, required))
		case reflect.Int:
			fValue.SetInt(int64(retrieveIntEnvVar(envVarName, required)))
		case reflect.Bool:
			fValue.SetBool(retrieveBoolEnvVar(envVarName))
		default:
			panic(fmt.Sprintf("Field '%s': unsupported field's type (%s)",
				structField.Name, structField.Type.Kind()))
		}
	}

}

func retrieveBoolEnvVar(envVarName string) bool {
	envVarValue := os.Getenv(envVarName)
	if envVarValue == "" {
		return false
	}

	return true
}

func retrieveIntEnvVar(envVarName string, required bool) int {
	var intEnvVarValue int
	strEnvVarValue := retrieveEnvVar(envVarName, required)
	if strEnvVarValue == "" && !required {
		return intEnvVarValue
	}

	intEnvVarValue, err := strconv.Atoi(strEnvVarValue)
	if err != nil {
		panic(fmt.Sprintf("Environment variable '%s' is not an integer", envVarName))
	}

	return intEnvVarValue
}

func retrieveEnvVar(envVarName string, required bool) string {
	envVarValue := os.Getenv(envVarName)
	if envVarValue == "" && required {
		panic(fmt.Sprintf("Environment variable '%s' is not defined or empty", envVarName))
	}

	return envVarValue
}
