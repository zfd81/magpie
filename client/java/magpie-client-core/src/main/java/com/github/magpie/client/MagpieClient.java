package com.github.magpie.client;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.magpie.*;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.NameResolver;
import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;

import java.io.*;
import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

@Slf4j
public class MagpieClient {

    private final MagpieClientConfig magpieClientConfig;
    private boolean isAvailable;

    private MagpieGrpc.MagpieBlockingStub blockingStub;
    private MagpieGrpc.MagpieStub asyncStub;

    public MagpieClient(MagpieClientConfig magpieClientConfig) throws MagpieRpcException {
        this.magpieClientConfig = magpieClientConfig;
        this.isAvailable = this.magpieClientConfig.isEnabled();
        if (this.isAvailable) {
            this.initStub();
        }
    }

    /**
     * 初始化 stub
     */
    private void initStub() throws MagpieRpcException {

        log.info("Initializing grpc stub...");
        log.info("magpieClientConfig: {}", magpieClientConfig);

        List<InetSocketAddress> realServerNodes = listMembers();
        NameResolver.Factory nameResolverFactory = new MultiAddressNameResolverFactory(realServerNodes);
        ManagedChannel channel = ManagedChannelBuilder.forTarget("service")
            .nameResolverFactory(nameResolverFactory)
            .defaultLoadBalancingPolicy(magpieClientConfig.getLoadBalancePolicy().name())
            .usePlaintext()
            .build();

        this.blockingStub = MagpieGrpc.newBlockingStub(channel);
        this.asyncStub = MagpieGrpc.newStub(channel);

        log.info("Rpc stub 初始化完成");
    }

    /**
     * 获得可用的服务端地址列表
     * @return 服务端地址列表
     */
    public List<InetSocketAddress> listMembers() throws MagpieRpcException {

        List<InetSocketAddress> result = new ArrayList<>();
        NameResolver.Factory nameResolverFactory = new MultiAddressNameResolverFactory(magpieClientConfig.getServerNodes());
        ManagedChannel channel = ManagedChannelBuilder.forTarget("service")
            .nameResolverFactory(nameResolverFactory)
            .defaultLoadBalancingPolicy(magpieClientConfig.getLoadBalancePolicy().name())
            .usePlaintext()
            .build();

        ClusterGrpc.ClusterBlockingStub clusterBlockingStub= ClusterGrpc.newBlockingStub(channel);
        RpcRequest request = RpcRequest.newBuilder().build();
        RpcResponse response = clusterBlockingStub.listMembers(request);
        channel.shutdown();

        ObjectMapper objectMapper = new ObjectMapper();
        try {
            List dataList = objectMapper.readValue(response.getData(), List.class);
            dataList.forEach(o -> {
                Map member = (Map<String, Object>) o;
                String addr = (String) member.get("addr");
                int port = (int) member.get("port");
                result.add(new InetSocketAddress(addr, port));
            });
        } catch (Exception e) {
            log.error("Error listMembers", e);
            throw new MagpieRpcException(e);
        }

        log.info("Available Magpie nodes: {}", result);
        return result;
    }

    /**
     * 异步加载表数据
     * @param tableName 表名
     * @param file 用于读取表数据信息的 File
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, File file, final Callback<LoadResponse> callback) throws IOException {
        load(tableName, new FileInputStream(file), callback);
    }

    /**
     * 异步加载表数据
     * @param tableName 表名
     * @param in 用于读取表数据信息的 InputStream
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, InputStream in, final Callback<LoadResponse> callback) throws IOException {
        load(tableName, new BufferedReader(new InputStreamReader(in)), callback);
    }

    /**
     * 异步加载表数据
     * @param tableName 表名
     * @param bufferedReader 用于读取表数据信息的 bufferedReader
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, BufferedReader bufferedReader, final Callback<LoadResponse> callback) throws IOException {

        CountDownLatch countDownLatch = new CountDownLatch(1);

        StreamObserver<LoadResponse> responseObserver = new StreamObserver<LoadResponse>() {
            @Override
            public void onNext(LoadResponse response) {
                log.info("响应: {}", response);
                callback.setResult(response);
            }

            @Override
            public void onError(Throwable t) {
                log.error("发生异常", t);
                countDownLatch.countDown();
            }

            @Override
            public void onCompleted() {
                log.info("成功完成");
                countDownLatch.countDown();
            }
        };

        StreamObserver<StreamRequest> requestObserver = asyncStub.load(responseObserver);

        requestObserver.onNext(StreamRequest.newBuilder().setData(tableName).build());
        String line;
        while ((line = bufferedReader.readLine()) != null) {
            requestObserver.onNext(StreamRequest.newBuilder().setData(line).build());
            // 判断调用结束状态。如果整个调用已经结束，不用再继续发送信息
            if (countDownLatch.getCount() == 0) {
                return;
            }
        }

        requestObserver.onCompleted();

        try {
            //如果在规定时间内没有请求完，则让程序停止
            if(!countDownLatch.await(magpieClientConfig.getTimeout(), TimeUnit.MILLISECONDS)){
                log.warn("没有在规定时间 {} ms 内完成", magpieClientConfig.getTimeout());
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    /**
     * 同步执行 sql
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    public QueryResponse execute(String sql) {
        return execute(QueryType.SELECT, sql);
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
     * @param sql 待执行的 sql
     * @param callback 接收相应的回调
     */
    public void executeAsync(String sql, final Callback<QueryResponse> callback) {
        executeAsync(QueryType.SELECT, sql, callback);
    }

    /**
     * 异步执行 sql
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @param callback 接收相应的回调
     */
    public void executeAsync(QueryType queryType, String sql, final Callback<QueryResponse> callback) {

        CountDownLatch countDownLatch = new CountDownLatch(1);

        StreamObserver<QueryResponse> responseObserver = new StreamObserver<QueryResponse>() {
            @Override
            public void onNext(QueryResponse response) {
                log.info("响应: {}", response);
                callback.setResult(response);
            }

            @Override
            public void onError(Throwable t) {
                log.error("发生异常", t);
                countDownLatch.countDown();
            }

            @Override
            public void onCompleted() {
                log.info("成功完成");
                countDownLatch.countDown();
            }
        };

        QueryRequest request = QueryRequest.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();

        asyncStub.execute(request, responseObserver);

        try {
            //如果在规定时间内没有请求完，则让程序停止
            if(!countDownLatch.await(magpieClientConfig.getTimeout(), TimeUnit.MILLISECONDS)){
                log.warn("没有在规定时间 {} ms 内完成", magpieClientConfig.getTimeout());
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

}
