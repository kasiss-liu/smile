package main

import (
	"fmt"

	"github.com/kasiss-liu/smile"
)

//实现一个监控器结构
type monitor struct{}

func (m *monitor) HandleStart(c *smile.MonitorInfo) {
	//	fmt.Println("func start")
}
func (m *monitor) HandleEnd(c *smile.MonitorInfo) {
	//	fmt.Println("func end")
}

//实现一个测试路由方法
func test(c *smile.Combination) error {
	fmt.Println("hello world")
	c.WriteString("helo")
	return nil
}

//主函数
func main() {
	//
	m := &monitor{}
	//获取一个服务器引擎
	engine := smile.Default()
	//注册路由
	routeGroup := smile.NewRouteGroup()
	routeGroup.SetGET("", test)
	engine.SetRouteGroup(routeGroup)
	//注册监视器
	smile.MonitorOn()
	engine.SetMonitor(m)
	engine.GzipOn()
	engine.Run(":8000")
}
