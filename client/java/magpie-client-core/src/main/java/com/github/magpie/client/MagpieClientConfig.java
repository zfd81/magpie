package com.github.magpie.client;

import lombok.Data;

import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.List;

@Data
public class MagpieClientConfig {
    private boolean enabled;
    private List<InetSocketAddress> serverNodes;
    private LoadBalancePolicy loadBalancePolicy;

    public Builder newBuilder() {
        return new Builder();
    }

    @Data
    public class HostAndPort {
        private String host;
        private int port;
    }

    public static class Builder {
        private boolean enabled;
        private List<InetSocketAddress> serverNodes = new ArrayList<>();
        private LoadBalancePolicy loadBalancePolicy;

        public Builder setEnabled(boolean enabled) {
            this.enabled = enabled;
            return this;
        }

        public Builder setServerNodes(List<InetSocketAddress> serverNodes) {
            this.serverNodes = serverNodes;
            return this;
        }

        public Builder addServerNode(InetSocketAddress serverNode) {
            this.serverNodes.add(serverNode);
            return this;
        }

        public Builder setLoadBalancePolicy(LoadBalancePolicy loadBalancePolicy) {
            this.loadBalancePolicy = loadBalancePolicy;
            return this;
        }

        public MagpieClientConfig build() {
            MagpieClientConfig instance = new MagpieClientConfig();
            instance.setEnabled(this.enabled);
            instance.setLoadBalancePolicy(this.loadBalancePolicy);
            instance.setServerNodes(this.serverNodes);
            return instance;
        }
    }
}
