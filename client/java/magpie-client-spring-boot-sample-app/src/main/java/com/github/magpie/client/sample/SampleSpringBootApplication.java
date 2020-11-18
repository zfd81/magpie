package com.github.magpie.client.sample;

import com.github.magpie.LoadResponse;
import com.github.magpie.QueryResponse;
import com.github.magpie.client.Callback;
import com.github.magpie.client.MagpieClient;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;

@SpringBootApplication
public class SampleSpringBootApplication implements CommandLineRunner {

    @Autowired
    private MagpieClient magpieClient;

    public static final String TABLE_NAME = "dummy.sql";
    private static final String QUERY_SQL = "select dummy_name, dummy_value from dummy where dummy_name = 'name01'";

    public static void main(String[] args) {
        SpringApplication.run(SampleSpringBootApplication.class, args);
    }

    @Override
    public void run(String... args) throws Exception {

        InputStream in = SampleSpringBootApplication.class.getResourceAsStream("/" + TABLE_NAME);
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
