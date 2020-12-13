package com.github.magpie.client;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.magpie.*;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.NameResolver;
import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.IOUtils;

import java.io.*;
import java.net.InetSocketAddress;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

@Slf4j
public class MagpieClient {

    private final MagpieClientConfig magpieClientConfig;
    private boolean isAvailable;

    private MagpieGrpc.MagpieBlockingStub magpieBlockingStub;
    private MagpieGrpc.MagpieStub magpieStub;
    private StorageGrpc.StorageBlockingStub storageBlockingStub;
    private MetaGrpc.MetaBlockingStub metaBlockingStub;
    private LogGrpc.LogBlockingStub logBlockingStub;
    private ClusterGrpc.ClusterBlockingStub clusterBlockingStub;

    private Storage storage;
    private Meta meta;
    private Log logNameSpace;
    private Cluster cluster;

    public MagpieClient(MagpieClientConfig magpieClientConfig) throws MagpieRpcException {
        this.magpieClientConfig = magpieClientConfig;
        this.isAvailable = this.magpieClientConfig.isEnabled();
        if (this.isAvailable) {
            this.initStub();
        }

        this.storage = new Storage();
        this.meta = new Meta();
        this.logNameSpace = new Log();
        this.cluster = new Cluster();
    }

    /**
     * 初始化 stub
     */
    private void initStub() throws MagpieRpcException {

        log.info("Initializing grpc stub...");
        log.info("magpieClientConfig: {}", magpieClientConfig);

        List<InetSocketAddress> availableServerNodes = getAvailableServerNodes();
        NameResolver.Factory nameResolverFactory = new MultiAddressNameResolverFactory(availableServerNodes);
        ManagedChannel channel = ManagedChannelBuilder.forTarget("service")
            .nameResolverFactory(nameResolverFactory)
            .defaultLoadBalancingPolicy(magpieClientConfig.getLoadBalancePolicy().name())
            .usePlaintext()
            .build();

        this.magpieBlockingStub = MagpieGrpc.newBlockingStub(channel);
        this.magpieStub = MagpieGrpc.newStub(channel);
        this.storageBlockingStub = StorageGrpc.newBlockingStub(channel);
        this.metaBlockingStub = MetaGrpc.newBlockingStub(channel);
        this.logBlockingStub = LogGrpc.newBlockingStub(channel);
        this.clusterBlockingStub = ClusterGrpc.newBlockingStub(channel);

        log.info("Rpc stub 初始化完成");
    }

    public Storage getStorage() {
        return this.storage;
    }

    public Meta getMeta() {
        return this.meta;
    }

    public Log getLog() {
        return this.logNameSpace;
    }

    public Cluster getCluster() {
        return this.cluster;
    }

    /**
     * 获得可用的服务端地址列表
     * @return 服务端地址列表
     */
    private List<InetSocketAddress> getAvailableServerNodes() throws MagpieRpcException {

        List<InetSocketAddress> result = new ArrayList<>();
        NameResolver.Factory nameResolverFactory = new MultiAddressNameResolverFactory(magpieClientConfig.getServerNodes());
        ManagedChannel channel = ManagedChannelBuilder.forTarget("service")
            .nameResolverFactory(nameResolverFactory)
            .defaultLoadBalancingPolicy(magpieClientConfig.getLoadBalancePolicy().name())
            .usePlaintext()
            .build();

        ClusterGrpc.ClusterBlockingStub tempClusterBlockingStub= ClusterGrpc.newBlockingStub(channel);
        RpcRequest request = RpcRequest.newBuilder().build();
        RpcResponse response = tempClusterBlockingStub.listMembers(request);
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
     * Magpie - 获得集群成员列表
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    private MembersResponse members(String sql) {
        Request request = Request.newBuilder()
            .setQueryType(QueryType .SELECT)
            .setSql(sql)
            .build();
        return magpieBlockingStub.members(request);
    }

    /**
     * Magpie - 获得集群成员列表
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    private MembersResponse members(QueryType queryType, String sql) {
        Request request = Request.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();
        return magpieBlockingStub.members(request);
    }

    /**
     * Magpie - 异步加载表数据
     * @param tableName 表名
     * @param data 表数据信息
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, String data, final Callback<LoadResponse> callback) throws IOException {
        InputStream inputStream = IOUtils.toInputStream(data, "UTF-8");
        load(tableName, new BufferedReader(new InputStreamReader(inputStream)), callback);
    }

    /**
     * Magpie - 异步加载表数据
     * @param tableName 表名
     * @param file 用于读取表数据信息的 File
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, File file, final Callback<LoadResponse> callback) throws IOException {
        load(tableName, new FileInputStream(file), callback);
    }

    /**
     * Magpie - 异步加载表数据
     * @param tableName 表名
     * @param in 用于读取表数据信息的 InputStream
     * @param callback 接收响应的回调
     * @throws IOException
     */
    public void load(String tableName, InputStream in, final Callback<LoadResponse> callback) throws IOException {
        load(tableName, new BufferedReader(new InputStreamReader(in)), callback);
    }

    /**
     * Magpie - 异步加载表数据
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

        StreamObserver<StreamRequest> requestObserver = magpieStub.load(responseObserver);

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
     * Magpie - 同步查询
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    public Response query(String sql) {
        Response response = query(QueryType.SELECT, sql);
        log.info("响应: {}", response);
        return response;
    }

    /**
     * Magpie - 同步查询
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    private Response query(QueryType queryType, String sql) {
        Request request = Request.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();
        return magpieBlockingStub.query(request);
    }

    /**
     * Magpie - 同步更新
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    public Response update(String sql) {
        Response response = update(QueryType.UPDATE, sql);
        log.info("响应: {}", response);
        return response;
    }

    /**
     * Magpie - 同步更新
     * @param queryType 类型
     * @param sql 待执行的 sql
     * @return 响应信息
     */
    private Response update(QueryType queryType, String sql) {
        Request request = Request.newBuilder()
            .setQueryType(queryType)
            .setSql(sql)
            .build();
        return magpieBlockingStub.update(request);
    }

    public class Storage {
        /**
         * Storage - 查看键值对数量
         * @param params 参数
         * @param data 数据
         * @return 响应信息
         */
        public RpcResponse count(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = storageBlockingStub.count(request);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Storage - 根据Key或者Key前缀查看数据
         * @param params 参数
         * @param data 数据
         * @return 响应信息
         */
        public RpcResponse get(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = storageBlockingStub.get(request);
            log.info("响应: {}", response);
            return response;
        }
    }

    public class Meta {

        /**
         * Meta - 创建表
         *
         * @param params 参数
         * @param data   数据
         * @return 响应信息
         */
        public RpcResponse createTable(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = metaBlockingStub.createTable(request);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Meta - 删除表
         *
         * @param params 参数
         * @param data   数据
         * @return 响应信息
         */
        public RpcResponse deleteTable(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = metaBlockingStub.deleteTable(request);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Meta - 查看表结构信息
         *
         * @param params 参数
         * @param data   数据
         * @return 响应信息
         */
        public RpcResponse describeTable(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = metaBlockingStub.describeTable(request);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Meta - 列出所有表格
         *
         * @param params 参数
         * @param data   数据
         * @return 响应信息
         */
        public RpcResponse listTables(Map<String, String> params, String data) {
            RpcRequest request = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            RpcResponse response = metaBlockingStub.listTables(request);
            log.info("响应: {}", response);
            return response;
        }
    }

    public class Log {

        /**
         * Log - 请求日志
         *
         * @param index     索引
         * @param data      数据
         * @param team      组
         * @param address   地址
         * @param port      端口
         * @param timestamp 时间戳
         * @return 响应信息
         */
        public RpcResponse apply(long index, String data, String team, String address, int port, String timestamp) {
            Entry entry = Entry.newBuilder()
                .setIndex(index)
                .setData(data)
                .setTeam(team)
                .setAddress(address)
                .setPort(port)
                .setTimestamp(timestamp)
                .build();
            RpcResponse response = logBlockingStub.apply(entry);
            log.info("响应: {}", response);
            return response;
        }
    }

    public class Cluster {

        /**
         * Cluster - 服务端流式响应
         *
         * @param params 参数
         * @param data   数据
         * @return 响应信息迭代器
         * @throws IOException
         */
        public Iterator<StreamResponse> dataSync(Map<String, String> params, String data) {
            RpcRequest rpcRequest = RpcRequest.newBuilder()
                .putAllParams(params)
                .setData(data)
                .build();
            Iterator<StreamResponse> response = clusterBlockingStub.dataSync(rpcRequest);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Cluster - 列出集群所有成员信息
         *
         * @return 响应信息
         */
        public RpcResponse listMembers() {
            RpcRequest request = RpcRequest.newBuilder().build();
            RpcResponse response = clusterBlockingStub.listMembers(request);
            log.info("响应: {}", response);
            return response;
        }

        /**
         * Cluster - 查看集群成员状态信息
         *
         * @return 响应信息
         */
        public RpcResponse memberStatus() {
            RpcRequest request = RpcRequest.newBuilder().build();
            RpcResponse response = clusterBlockingStub.memberStatus(request);
            log.info("响应: {}", response);
            return response;
        }
    }
}
