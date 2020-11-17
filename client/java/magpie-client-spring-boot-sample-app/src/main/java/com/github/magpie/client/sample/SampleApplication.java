package com.github.magpie.client.sample;

import com.github.magpie.LoadResponse;
import com.github.magpie.StreamRequest;
import com.github.magpie.client.MagpieClient;
import com.github.magpie.client.sample.service.MagpieService;
import io.grpc.stub.StreamObserver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;

@SpringBootApplication
public class SampleApplication implements CommandLineRunner {

    @Autowired
    private MagpieClient magpieClient;

    public static final String TABLE_NAME = "dummy.sql";

    public static void main(String[] args) {
        SpringApplication.run(SampleApplication.class, args);
    }

    @Override
    public void run(String... args) throws Exception {
        StreamObserver<LoadResponse> responseObserver = new StreamObserver<LoadResponse>() {
            @Override
            public void onNext(LoadResponse response) {
                System.out.println("receive: " + response);
            }

            @Override
            public void onError(Throwable t) {
                System.out.println("error: ");
            }

            @Override
            public void onCompleted() {
                System.out.println("completed: ");
            }
        };

        InputStream in = SampleApplication.class.getResourceAsStream("/" + TABLE_NAME);
        BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(in));
        magpieClient.load(TABLE_NAME, bufferedReader, responseObserver);
    }
}
