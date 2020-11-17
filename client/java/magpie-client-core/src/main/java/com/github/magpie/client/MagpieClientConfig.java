package com.github.magpie.client;

import lombok.Data;

import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

@Data
public class MagpieClientConfig {
    private boolean enabled = true;
    private List<InetSocketAddress> serverNodes = new ArrayList<>();
    private LoadBalancePolicy loadBalancePolicy = LoadBalancePolicy.round_robin;
    private long timeout = 5 * 60 * 1000; // 5 分钟

    public static Builder newBuilder() {
        return new Builder();
    }

    @Data
    public class HostAndPort {
        private String host;
        private int port;
    }

    public static class Builder {

        private MagpieClientConfig instance;

        public Builder() {
            this.instance = new MagpieClientConfig();
        }

        public Builder setTimeout(long timeout) {
            this.instance.setTimeout(timeout);
            return this;
        }

        public Builder setEnabled(boolean enabled) {
            this.instance.setEnabled(enabled);
            return this;
        }

        public Builder setServerNodes(List<InetSocketAddress> serverNodes) {
            this.instance.setServerNodes(serverNodes);
            return this;
        }

        public Builder parseServerNodesFromString(String hostAndPortsString) {
            String[] split = hostAndPortsString.split(",");
            Arrays.stream(split).forEach(s -> {
                String[] hostAndPort = s.trim().split(":");
                if (hostAndPort.length != 2) {
                    throw new RuntimeException("Wrong inet socket address");
                }
                InetSocketAddress serverNode = new InetSocketAddress(hostAndPort[0], Integer.parseInt(hostAndPort[1]));
                this.addServerNode(serverNode);
            });
            return this;
        }

        public Builder addServerNode(InetSocketAddress serverNode) {
            this.instance.getServerNodes().add(serverNode);
            return this;
        }

        public Builder setLoadBalancePolicy(LoadBalancePolicy loadBalancePolicy) {
            this.instance.setLoadBalancePolicy(loadBalancePolicy);
            return this;
        }

        public MagpieClientConfig build() {
            return instance;
        }
    }
}
