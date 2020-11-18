package com.github.magpie.client.sample;

import com.github.magpie.LoadResponse;
import com.github.magpie.QueryResponse;
import com.github.magpie.client.*;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

public class SampleApplication {

    public static final String TABLE_NAME = "dummy.sql";
    private static final String SERVER_NODES = "localhost:50000,localhost:50001,localhost:50002";
    private static final String QUERY_SQL = "select dummy_name, dummy_value from dummy where dummy_name = 'name01'";

    public static void main(String[] args) throws IOException {

        MagpieClientConfig magpieClientConfig = MagpieClientConfig.newBuilder()
            .setEnabled(true)
            .setLoadBalancePolicy(LoadBalancePolicy.round_robin)
            .parseServerNodesFromString(SERVER_NODES)
            .setTimeout(5 * 60 * 1000)
            .build();

        MagpieClient magpieClient = new MagpieClient(magpieClientConfig);

        // 加载数据
        InputStream in = SampleApplication.class.getResourceAsStream("/" + TABLE_NAME);
        BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(in));
        Callback<LoadResponse> loadCallback = new Callback<>();
        magpieClient.load(TABLE_NAME, bufferedReader, loadCallback);
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
