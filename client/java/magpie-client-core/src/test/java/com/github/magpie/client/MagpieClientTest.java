package com.github.magpie.client;

import com.github.magpie.LoadResponse;
import com.github.magpie.Response;
import com.github.magpie.RpcResponse;
import com.github.magpie.StreamResponse;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class MagpieClientTest {

    public static final String TABLE_NAME = "userInfo";
    public static final String DATA_FILE = "/userInfo.csv";
    private static final String SERVER_NODES = "127.0.0.1:8143";
    private static final String QUERY_SQL = "select id,name,pwd,age from userInfo where id = '1'";

    private static MagpieClient magpieClient;

    @BeforeAll
    static void init() throws MagpieRpcException {
        MagpieClientConfig magpieClientConfig = MagpieClientConfig.newBuilder()
            .setEnabled(true)
            .setLoadBalancePolicy(LoadBalancePolicy.round_robin)
            .parseServerNodesFromString(SERVER_NODES)
            .setTimeout(5 * 60 * 1000)
            .build();

        magpieClient = new MagpieClient(magpieClientConfig);
    }

    /** 加载数据 */
    @org.junit.jupiter.api.Test
    void load() throws IOException {
        InputStream in = this.getClass().getResourceAsStream(DATA_FILE);
        Callback<LoadResponse> loadCallback = new Callback<>();
        magpieClient.load(TABLE_NAME, in, loadCallback);
        assertNotNull(loadCallback.getResult());
        System.out.println("加载数据结果: ");
    }

    /** 查询 */
    @org.junit.jupiter.api.Test
    void query() {
        Response executeResult = magpieClient.query(QUERY_SQL);
        assertNotNull(executeResult);
        System.out.println("查询结果: ");
        System.out.println(executeResult);
    }

    /** Storage - 查看键值对数量 */
    @Test
    void storageCount() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getStorage().count(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Storage - 根据Key或者Key前缀查看数据 */
    @Test
    void storageGet() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getStorage().get(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Meta - 创建表 */
    @Test
    void metaCreateTable() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getMeta().createTable(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Meta - 删除表 */
    @Test
    void metaDeleteTable() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getMeta().deleteTable(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Meta - 查看表结构信息 */
    @Test
    void metaDescribeTable() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getMeta().describeTable(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Meta - 列出所有表格 */
    @Test
    void metaListTables() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        RpcResponse response = magpieClient.getMeta().listTables(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Log - 请求日志 */
    @Test
    void logApply() {
        Map<String, String> params = new HashMap<>();
        long index = 0;
        String data = "";
        String team = "";
        String address = "";
        int port = 0;
        String timestamp = "";
        RpcResponse response = magpieClient.getLog().apply(index, data, team, address, port, timestamp);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Cluster - 服务端流式响应 */
    @Test
    void clusterDataSync() {
        Map<String, String> params = new HashMap<>();
        String data = "";
        Iterator<StreamResponse> response = magpieClient.getCluster().dataSync(params, data);
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Cluster - 列出集群所有成员信息 */
    @Test
    void clusterListMembers() {
        RpcResponse response = magpieClient.getCluster().listMembers();
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }

    /** Cluster - 查看集群成员状态信息 */
    @Test
    void clusterMemberStatus() {
        RpcResponse response = magpieClient.getCluster().memberStatus();
        assertNotNull(response);
        System.out.println("响应: ");
        System.out.println(response);
    }
}
