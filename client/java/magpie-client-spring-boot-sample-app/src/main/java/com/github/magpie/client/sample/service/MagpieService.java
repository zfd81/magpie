package com.github.magpie.client.sample.service;

import com.github.magpie.*;
import io.grpc.stub.StreamObserver;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Random;

public class MagpieService extends MagpieGrpc.MagpieImplBase {

    public StreamObserver<StreamRequest> load(StreamObserver<LoadResponse> responseObserver) {
        SimpleDateFormat simpleDateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss.SSS");
        return new StreamObserver<StreamRequest>() {

            private List<String> dataLines = new ArrayList<>();

            @Override
            public void onNext(StreamRequest value) {
                dataLines.add(value.getData());
            }

            @Override
            public void onError(Throwable t) {
                t.printStackTrace();
            }

            @Override
            public void onCompleted() {

                long start = System.currentTimeMillis();
                Date startDate = new Date();
                startDate.setTime(start);
                String startTime = simpleDateFormat.format(startDate);

                try {
                    Random random = new Random();
                    Thread.sleep(random.nextInt(10) * 1000);
                } catch (InterruptedException e) {
                    // Ignore
                }

                long end = System.currentTimeMillis();
                Date endDate = new Date();
                endDate.setTime(start);
                String endTime = simpleDateFormat.format(endDate);

                LoadResponse response;
                if (dataLines.isEmpty()) {
                    response = LoadResponse.newBuilder()
                        .setCode(500)
                        .setName(null)
                        .setStartTime(startTime)
                        .setEndTime(endTime)
                        .setElapsedTime(end - start)
                        .setRecordCount(dataLines.size() - 1)
                        .setMessage("错误：接收到的数据行数为0.")
                        .build();

                } else {
                    response = LoadResponse.newBuilder()
                        .setCode(200)
                        .setName(dataLines.get(0))
                        .setStartTime(startTime)
                        .setEndTime(endTime)
                        .setElapsedTime(end - start)
                        .setRecordCount(dataLines.size() - 1)
                        .setMessage("OK")
                        .build();
                }

                responseObserver.onNext(response);
                responseObserver.onCompleted();
            }
        };
    }

    @Override
    public void execute(QueryRequest request, StreamObserver<QueryResponse> responseObserver) {

        QueryResponse queryResponse = QueryResponse.newBuilder()
            .setCode(200)
            .setData("dummy data")
            .setMessage("dummy message")
            .setDataType(DataType.STRING)
            .build();
        responseObserver.onNext(queryResponse);
        responseObserver.onCompleted();
    }
}
