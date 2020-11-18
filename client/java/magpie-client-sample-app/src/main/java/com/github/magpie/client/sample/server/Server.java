package com.github.magpie.client.sample.server;

import com.github.magpie.client.sample.service.MagpieService;
import io.grpc.ServerBuilder;

import java.io.IOException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Server {

    public static void main(String[] args) {
        final int nServers = 3;
        ExecutorService executorService = Executors.newFixedThreadPool(nServers);
        for (int i = 0; i < nServers; i++) {
            String name = "Server_" + i;
            int port = 50000 + i;
            executorService.submit(() -> {
                try {
                    startServer(name, port);
                } catch (IOException | InterruptedException e) {
                    e.printStackTrace();
                }
            });
        }
    }

    private static void startServer(String name, int port) throws IOException, InterruptedException {
        io.grpc.Server server = ServerBuilder
            .forPort(port)
            .addService(new MagpieService())
            .build();

        server.start();
        System.out.println(name + " server started, listening on port: " + server.getPort());
        server.awaitTermination();
    }
}
