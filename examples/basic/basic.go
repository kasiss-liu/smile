package main

import (
	"fmt"

	"github.com/kasiss-liu/smile"
)

//实现一个测试路由方法
func test(c *smile.Combination) error {
	fmt.Println("hello world")
	c.WriteString("helo")
	return nil
}

//主函数
func main() {
	//获取一个服务器引擎
	engine := smile.Default()
	//注册路由
	routeGroup := smile.NewRouteGroup()
	routeGroup.SetGET("", test)
	engine.SetRouteGroup(routeGroup)
	engine.GzipOn()
	engine.Run(":8000")
}
