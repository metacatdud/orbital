package dom

import (
	"encoding/json"
	"reflect"
	"syscall/js"
)

func ConsoleLog(data ...interface{}) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func ConsoleError(data ...interface{}) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func ConsoleWarn(data ...interface{}) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func processConsoleData(data ...interface{}) []interface{} {
	var processedArgs []interface{}

	for _, d := range data {
		if d == nil {
			processedArgs = append(processedArgs, "nil")
			continue
		}

		switch v := d.(type) {
		case string, bool, int, int32, int64, uint, uint32, uint64, float32, float64:
			processedArgs = append(processedArgs, v)
		case js.Value:
			if v.IsUndefined() || v.IsNull() {
				processedArgs = append(processedArgs, "undefined or null")
			} else {
				processedArgs = append(processedArgs, v)
			}
		default:
			kind := reflect.TypeOf(d).Kind()
			if kind == reflect.Struct || kind == reflect.Ptr || kind == reflect.Map || kind == reflect.Slice {
				jsonData, err := json.Marshal(v)
				if err != nil {
					processedArgs = append(processedArgs, "Error marshalling to JSON:", err.Error())
				} else {
					processedArgs = append(processedArgs, string(jsonData))
				}
				continue
			}

			processedArgs = append(processedArgs, "Unsupported type:", reflect.TypeOf(d).String())
		}
	}

	return processedArgs
}
