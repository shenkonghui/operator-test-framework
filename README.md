# operator-test-framework
operator测试框架是一种通用的测试框架，主要针对K8s资源变更情况进行测试。


## 快速开始
执行测试命令
```yaml
git clone http://shenkonghui/shenkonghui/operator-test-framework.git
 
cd operator-test-framework && go build 
 
./operator-test-framework run -c <job-file-dir> -v=2
```
输出
```yaml
{
    "test_run_results":[
        {
            "test_run_name":"pod-test",
            "status":"completed",
            "time":"1.13101142s",
            "case_results":[
                {
                    "case_name":"Init",
                    "status":"completed",
                    "error_log":"",
                    "time":"1.13098585s"
                }
            ]
        }
    ]
}
```
## 流程
![image.png](https://cdn.nlark.com/yuque/0/2021/png/487266/1619581580347-475bfb18-d471-498b-a9ef-299b7023ce39.png#align=left&display=inline&height=849&margin=%5Bobject%20Object%5D&name=image.png&originHeight=849&originWidth=498&size=168256&status=done&style=none&width=498)
## 配置文件
支持目录 + 文件 ，使用参数-c 指定
### 字段说明

- name  测试名称
- parameter 参数，会获取key中的属性#{xxx}进行替换, 可以直接命令行-p指定参数
- job 一个配置文件就是一个job，下面有多个case
   - initTime 初始化时间
   - timeout 超时时间
   - interval 检查间隔
   - variable 变量
      - key 键
      - valueFromCmd 值从命令行获取
      - value 静态值
   - verificate 验证结果
      - name 验证名称
      - cmd 执行命令
      - value 验证的属性值，可以从variable.value 获取变量进行动态验证，比如revison，保证cr更新了
      - operator 验证操作，目前包含等于(equal)/不等于(noEqual), 默认等于



### demo
#### 常规pod
如下的配置，只是检查一个名为"redis-cl-1" 的 pod是否是running的状态。
验证的过程主要是执行kubectl get pod #name -o=jsonpath='{.status.phase}' 检查是否是Running状态
```yaml
name: "pod-test"
parameter:
  - key: name
    value: redis-cl-1
jobs:
  # 初始化状态检查
  - name: "Init"
    timeout: 10s
    interval: 5s
    verificate:
      - cmd: "kubectl get pod #{name} -o=jsonpath='{.status.phase}'"
        value: "Running"
        name: "Running"
```


