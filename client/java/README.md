# Magpie Client

Magpie 客户端,支持 load balance 。


# 集成步骤

1. 添加 Maven 依赖
```xml
<dependency>
    <groupId>com.github.magpie</groupId>
    <artifactId>magpie-client-spring-boot-starter</artifactId>
    <version>1.0-SNAPSHOT</version>
</dependency>
```

2. 添加配置项
```properties
# 是否启用 magpie 客户端, 默认值是 true
magpie.client.enabled=true
# Load balance 策略， 可用的值有 round_robin, grpclb, pick_first, 默认值是 round_robin
magpie.client.load-balance-policy=round_robin
# 服务器的地址，多个地址用逗号分隔
magpie.client.server-nodes=localhost:50000,localhost:50001,localhost:50002
```

## 运行示例

1. 运行 magpie-client-spring-boot-sample-app 应用里的 Server 。
2. 运行 magpie-client-spring-boot-sample-app 应用里的 SampleApplication 。
3. 在 SampleApplication 的控制台查看输出结果。
