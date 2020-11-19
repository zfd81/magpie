# Magpie Client

Magpie 客户端,支持 load balance 。

# Maven 工程集成步骤

1. 添加 Maven 依赖
```xml
<dependency>
    <groupId>com.github.magpie</groupId>
    <artifactId>magpie-client-core</artifactId>
    <version>1.0-SNAPSHOT</version>
</dependency>
```
2. 调用方法
```java
public class SampleApplication {

    public static final String TABLE_NAME = "userInfo";
    public static final String DATA_FILE = "/userInfo.csv";
    private static final String SERVER_NODES = "127.0.0.1:8143";
    private static final String QUERY_SQL = "select id,name,pwd,age from userInfo where id = '1'";

    public static void main(String[] args) throws IOException {

        MagpieClientConfig magpieClientConfig = MagpieClientConfig.newBuilder()
            .setEnabled(true)
            .setLoadBalancePolicy(LoadBalancePolicy.round_robin)
            .parseServerNodesFromString(SERVER_NODES)
            .setTimeout(5 * 60 * 1000)
            .build();

        MagpieClient magpieClient = new MagpieClient(magpieClientConfig);

        // 加载数据
        InputStream in = SampleApplication.class.getResourceAsStream(DATA_FILE);
        Callback<LoadResponse> loadCallback = new Callback<>();
        magpieClient.load(TABLE_NAME, in, loadCallback);
        System.out.println("加载数据结果: ");
        System.out.println(loadCallback.getResult());

        // 同步查询
        QueryResponse executeResult = magpieClient.execute(QUERY_SQL);
        System.out.println("同步查询结果: ");
        System.out.println(executeResult);

        // 异步查询
        Callback<QueryResponse> queryCallback = new Callback<>();
        magpieClient.executeAsync(QUERY_SQL, queryCallback);
        System.out.println("异步查询结果: ");
        System.out.println(queryCallback.getResult());
    }
}
```

# Maven 工程集成示例
1. 运行 magpie-client-sample-app 应用里的 Server 。
2. 运行 magpie-client-sample-app 应用里的 SampleApplication 。
3. 在 SampleApplication 的控制台查看输出结果。

# Spring Boot 集成步骤

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
# 请求超时时间，单位毫秒, 默认值是 300000 （5分钟）
magpie.client.timeout=300000
# 服务器的地址，多个地址用逗号分隔
magpie.client.server-nodes=localhost:50000,localhost:50001,localhost:50002
```

# Spring Boot 集成示例

1. 运行 magpie-client-spring-boot-sample-app 应用里的 Server 。
2. 运行 magpie-client-spring-boot-sample-app 应用里的 SampleSpringBootApplication 。
3. 在 SampleSpringBootApplication 的控制台查看输出结果。
