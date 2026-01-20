package internal

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"scheduler-demo/pkg"
)

// 未使用变量，但是因为类型是接口，所以 sample 需要实现接口才能被作为值传给变量
// 目的不是为了被引用而是检查接口的实现情况
var _ framework.PreFilterPlugin = &Sample{}
var _ framework.FilterPlugin = &Sample{}
var _ framework.ScorePlugin = &Sample{}

// Name 大写为了被引用
const (
	Name = "node-order"
)

// 固定格式 handle
type Sample struct {
	handle framework.Handle
}

func (s *Sample) Name() string {
	return Name
}

// 固定写法，内容可以填充
func (s *Sample) PreFilter(ctx context.Context, state *framework.CycleState, p *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	return nil, framework.NewStatus(framework.Success)
}

// PreFilterExtensions returns prefilter extensions, pod add and remove.
func (s *Sample) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// Filter : evaluate if node can respect maxNetworkCost requirements
// 固定写法
func (s *Sample) Filter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	allocatable := nodeInfo.Allocatable.MilliCPU
	if allocatable == 0 {
		return framework.NewStatus(framework.Unschedulable, "no cpu allocatable")
	}

	used := nodeInfo.Requested.MilliCPU
	usage := float64(used) / float64(allocatable)

	if usage >= 0.8 {
		return framework.NewStatus(
			framework.Unschedulable,
			fmt.Sprintf("cpu usage %.2f exceeds 80%%", usage),
		)
	}

	return framework.NewStatus(framework.Success)
}

// 固定写法，返回的 int64 是打分值，1-100，如果有多个score，最终分数是和， status 是返回状态
func (s *Sample) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {

	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.AsStatus(err)
	}

	alloc := nodeInfo.Allocatable.MilliCPU
	if alloc == 0 {
		return 0, framework.NewStatus(framework.Success)
	}

	used := nodeInfo.Requested.MilliCPU
	usage := float64(used) / float64(alloc)

	if usage >= 0.8 {
		return 0, framework.NewStatus(framework.Success)
	}

	nodeOrder := pkg.GetNodeOrderFromPod(p)

	// Pod 没指定顺序：给中等分，不抢
	if len(nodeOrder) == 0 {
		return 50, framework.NewStatus(framework.Success)
	}

	for i, n := range nodeOrder {
		if n == nodeName {
			score := int64(100 - i*20)
			if score < 1 {
				score = 1
			}
			return score, framework.NewStatus(framework.Success)
		}
	}

	// 不在列表里的 node，最低优先级
	return 1, framework.NewStatus(framework.Success)
}

func (s *Sample) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func (s *Sample) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	return nil
}

// New 函数，大写为了被引用
func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &Sample{
		handle: handle,
	}, nil
}
