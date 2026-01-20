## 调度器开发

### 功能描述

本调度器插件实现了以下功能：
1. **节点过滤**：当节点的 CPU 使用率超过 80% 时，任务将不会调度到该节点。
2. **节点打分**：根据 Pod 的 `scheduler.alpha.kubernetes.io/node-order` 注解，按照指定顺序为节点打分，优先调度到分数最高的节点。
   - 如果未指定顺序，节点将获得中等分数（50）。
   - 不在列表中的节点将获得最低优先级分数（1）。

### 目录架构
```
.
├── Dockerfile
├── REAME.md
├── deploy
│   ├── README.md      
│   ├── configmap.yaml 
│   ├── deployment.yaml
│   ├── pod.yaml       
│   └── rbac.yaml      
├── go.mod
├── go.sum
├── internal
│   └── sample.go
├── main.go
└── pkg
    └── utils.go
```

### 插件开发说明

调度器插件主要实现了 Kubernetes 调度框架中的以下接口：
- **PreFilter**：预过滤逻辑，当前未实现具体逻辑。
- **Filter**：过滤逻辑，检查节点的 CPU 使用率是否超过 80%。
- **Score**：打分逻辑，根据 Pod 的注解为节点分配优先级分数。

插件的实现代码位于 `/internal/sample.go`，并通过 `/main.go` 将插件注册到调度器中。

### 部署说明

1. **配置插件**
   在 `/deploy/configmap.yaml` 中配置调度器插件。

2. **部署调度器**
   使用 `/deploy/deployment.yaml` 部署调度器。

3. **权限配置**
   在 `/deploy/rbac.yaml` 中为调度器配置必要的权限。

### 测试方法

1. 创建测试 Pod：
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-cs001
     namespace: scheduler-plugins
     annotations:
       scheduler.alpha.kubernetes.io/node-order: node01,node02
     labels:
       app: test-cs
   spec:
     containers:
       - name: test-cs
         image: nginx
     schedulerName: node-order-scheduler
   ```

2. 观察调度结果：
   - 如果 `node01` 的 CPU 使用率低于 80%，Pod 将优先调度到 `node01`。
   - 如果 `node01` 的 CPU 使用率超过 80%，Pod 将调度到 `node02`。
   - 如果所有节点的 CPU 使用率均超过 80%，Pod 将无法调度。

3. 验证打分逻辑：
   - 修改 Pod 的注解 `scheduler.alpha.kubernetes.io/node-order`，观察节点的优先级变化。

### 注意事项

- 确保调度器的配置文件与插件名称一致。
- 部署前检查 Kubernetes 集群的版本是否兼容当前调度器插件。