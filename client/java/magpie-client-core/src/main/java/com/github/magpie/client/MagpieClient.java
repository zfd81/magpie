package com.github.magpie.client;

import com.github.magpie.*;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.NameResolver;
import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

import static com.github.magpie.client.MagpieClientConfigParams.*;

@Slf4j
public class MagpieClient {

    private final MagpieClientConfig magpieClientConfig;
    private boolean isAvailable;

    MagpieGrpc.MagpieBlockingStub blockingStub;
    MagpieGrpc.MagpieStub asyncStub;

    public MagpieClient(MagpieClientConfig magpieClientConfig) throws MagpieClientInitializeException {
        this.magpieClientConfig = magpieClientConfig;
        this.isAvailable = this.magpieClientConfig.isEnabled();
        if (this.isAvailable) {
            this.initStub();
        }
    }

    private void initStub() throws MagpieClientInitializeException {

        log.info("Initializing grpc stub...");
        log.info("config: {}", magpieClientConfig);

        NameResolver.Factory nameResolverFactory = new MultiAddressNameResolverFactory(magpieClientConfig.getServerNodes());

        ManagedChannel channel = ManagedChannelBuilder.forTarget("service")
            .nameResolverFactory(nameResolverFactory)
            .defaultLoadBalancingPolicy(magpieClientConfig.getLoadBalancePolicy().name())
            .usePlaintext()
            .build();

        this.blockingStub = MagpieGrpc.newBlockingStub(channel);
        this.asyncStub = MagpieGrpc.newStub(channel);

        log.info("rpc stub initialized successfully.");
    }

    private InetSocketAddress parseFrom(String hostAndPort) {
        String[] hostAndPorts = hostAndPort.split(":");
        if (hostAndPorts.length != 2) {
            throw new RuntimeException("wrong inet socket address: " + hostAndPort);
        }
        return new InetSocketAddress(hostAndPorts[0], Integer.parseInt(hostAndPorts[1]));
    }

    /**
     * 异步加载表数据
     * @param tableName 表名
     * @param bufferedReader 用于读取表数据信息的 bufferedReader
     * @param responseObserver 响应
     * @throws IOException
     */
    public void load(String tableName, BufferedReader bufferedReader, StreamObserver<LoadResponse> responseObserver) throws IOException {

        StreamObserver<StreamRequest> requestObserver = asyncStub.load(responseObserver);

        requestObserver.onNext(StreamRequest.newBuilder().setData(tableName).build());
        String line;
        while ((line = bufferedReader.readLine()) != null) {
            requestObserver.onNext(StreamRequest.newBuilder().setData(line).build());
        }

        requestObserver.onCompleted();
    }

    /**
     * 同步执行 sql
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    public QueryResponse execute(QueryType queryType, String sql) {
        QueryRequest request = QueryRequest.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();
        return blockingStub.execute(request);
    }

    /**
     * 异步执行 sql
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @param responseObserver
     */
    public void executeAsync(QueryType queryType, String sql, StreamObserver<QueryResponse> responseObserver) {
        QueryRequest request = QueryRequest.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();
        asyncStub.execute(request, responseObserver);
    }
}
