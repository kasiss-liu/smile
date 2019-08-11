package smile

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

//getFuncName return name of func
func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

}

//trimPath format path in one entrance
func trimPath(path string) string {
	return "/" + strings.Trim(path, "/")
}

func doPrintRoutes(routes []string) {
	for _, v := range routes {
		fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", v)
	}
}
