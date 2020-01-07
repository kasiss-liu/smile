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

func doPrintRoutes(routesAssign []string, routesAutoCreate []string) {
	if len(routesAssign) > 0{
		fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", "assigned:")
		for _, v := range routesAssign {
			fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", v)
		}
	}
	if len(routesAutoCreate) > 0 {
		fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", "autocreated:")
		for _, v := range routesAutoCreate {
			fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", v)
		}
	}
	if len(routesAssign) == 0 && len(routesAutoCreate) == 0 {
		fmt.Fprintf(os.Stdout, "[SMILE Route]%s\r\n", "No Route Registered!")
	}
}
