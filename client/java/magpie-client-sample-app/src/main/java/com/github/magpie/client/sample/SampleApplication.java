package com.github.magpie.client.sample;

import com.github.magpie.LoadResponse;
import com.github.magpie.client.*;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

public class SampleApplication {

    public static final String TABLE_NAME = "dummy.sql";
    private static final String SERVER_NODES = "localhost:50000,localhost:50001,localhost:50002";

    public static void main(String[] args) throws IOException {

        MagpieClientConfig magpieClientConfig = MagpieClientConfig.newBuilder()
            .setEnabled(true)
            .setLoadBalancePolicy(LoadBalancePolicy.round_robin)
            .parseServerNodesFromString(SERVER_NODES)
            .setTimeout(5 * 60 * 1000)
            .build();

        MagpieClient magpieClient = new MagpieClient(magpieClientConfig);

        InputStream in = SampleApplication.class.getResourceAsStream("/" + TABLE_NAME);
        BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(in));
        Callback<LoadResponse> callback = new Callback<>();
        magpieClient.load(TABLE_NAME, bufferedReader, callback);

        System.out.println(callback.getResult());

    }
}
