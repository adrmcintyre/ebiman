//go:build wasm

package platform

import (
	"syscall/js"
)

var (
	objects = map[string]js.Value{}
	global  = js.Global()
)

func call(objName string, fn string, params ...any) js.Value {
	obj, ok := objects[objName]
	if !ok {
		obj = global.Get(objName)
		objects[objName] = obj
	}
	if obj.Type() != js.TypeObject {
		return js.Undefined()
	}
	if obj.Get(fn).Type() != js.TypeFunction {
		return js.Undefined()
	}
	return obj.Call(fn, params...)
}

func callString(objName string, fn string, params ...any) (string, bool) {
	jsValue := call(objName, fn, params...)
	if jsValue.Type() == js.TypeString {
		return jsValue.String(), true
	}
	return "", false
}

func getLocalStorage(key string) (string, bool) {
	return callString("localStorage", "getItem", key)
}

func setLocalStorage(key string, value string) {
	call("localStorage", "setItem", key, value)
}

func randomUUID() (string, bool) {
	return callString("crypto", "randomUUID")
}

func GetDeviceId() (string, bool) {
	if deviceId, ok := getLocalStorage("deviceId"); ok && deviceId != "" {
		return deviceId, true
	}

	if uuid, ok := randomUUID(); ok {
		setLocalStorage("deviceId", uuid)
		return uuid, true
	}

	return "", false
}
