package com.github.magpie.client.sample;

import com.github.magpie.LoadResponse;
import com.github.magpie.Response;
import com.github.magpie.client.*;

import java.io.IOException;
import java.io.InputStream;

public class SampleApplication {

    public static final String TABLE_NAME = "userInfo";
    public static final String DATA_FILE = "/userInfo.csv";
    private static final String SERVER_NODES = "127.0.0.1:8143";
    private static final String QUERY_SQL = "select id,name,pwd,age from userInfo where id = '1'";

    public static void main(String[] args) throws IOException, MagpieRpcException {

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

        // 查询
        Response executeResult = magpieClient.query(QUERY_SQL);
        System.out.println("查询结果: ");
        System.out.println(executeResult);
    }
}
