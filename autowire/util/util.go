/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"fmt"
	"reflect"
)

const (
	emptyString = ""
)

func GetStructName(v interface{}) string {
	if v == nil {
		return ""
	}
	typeOfInterface := GetTypeFromInterface(v)
	return typeOfInterface.Name()
}

func GetSDIDByStructPtr(v interface{}) string {
	if v == nil {
		return ""
	}
	typeOfInterface := GetTypeFromInterface(v)
	return fmt.Sprintf("%s.%s", typeOfInterface.PkgPath(), typeOfInterface.Name())
}

func GetTypeFromInterface(v interface{}) reflect.Type {
	valueOfInterface := reflect.ValueOf(v)
	valueOfElemInterface := valueOfInterface.Elem()
	return valueOfElemInterface.Type()
}
