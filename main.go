package main

import (
	"os"

	"k8s.io/component-base/cli"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"scheduler-demo/internal"
)

func main() {
	// 此处用的 app.NewSchedulerCommand 是正儿八经 kubernetes scheduler 的 Command，说白了就是再起一个调度器。
	// 这边没有体现配置文件什么的，其实都在 NewSchedulerCommand 里面实现了。
	// 主要是 app.WithPlugin(pkg.Name, pkg.New) 这个是导入自定义插件

	command := app.NewSchedulerCommand(
		app.WithPlugin(internal.Name, internal.New))

	code := cli.Run(command)
	os.Exit(code)
}
