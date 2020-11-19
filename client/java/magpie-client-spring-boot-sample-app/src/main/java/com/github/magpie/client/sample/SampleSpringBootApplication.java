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

    public static final String TABLE_NAME = "userInfo";
    public static final String DATA_FILE = "/userInfo.csv";
    private static final String QUERY_SQL = "select id,name,pwd,age from userInfo where id = '1'";

    public static void main(String[] args) {
        SpringApplication.run(SampleSpringBootApplication.class, args);
    }

    @Override
    public void run(String... args) throws Exception {

        // 加载数据
        InputStream in = SampleSpringBootApplication.class.getResourceAsStream(DATA_FILE);
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
