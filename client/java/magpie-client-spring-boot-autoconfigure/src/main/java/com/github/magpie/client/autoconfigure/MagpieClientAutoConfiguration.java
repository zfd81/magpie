package com.github.magpie.client.autoconfigure;

import com.github.magpie.client.LoadBalancePolicy;
import com.github.magpie.client.MagpieClient;
import com.github.magpie.client.MagpieClientConfig;
import com.github.magpie.client.MagpieClientInitializeException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.net.InetSocketAddress;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static com.github.magpie.client.MagpieClientConfigParams.*;

@Configuration
@ConditionalOnClass(MagpieClient.class)
@EnableConfigurationProperties(MagpieClientProperties.class)
public class MagpieClientAutoConfiguration {

    @Autowired
    private MagpieClientProperties magpieClientProperties;

    @Bean
    @ConditionalOnMissingBean
    public MagpieClientConfig magpieClientConfig() throws MagpieClientInitializeException {

        Boolean enabled = magpieClientProperties.getEnabled();
        enabled = enabled == null ? true : enabled;

        String loadBalancePolicy  = magpieClientProperties.getLoadBalancePolicy();
        loadBalancePolicy = loadBalancePolicy == null ? "round_robin" : loadBalancePolicy;

        String serverNodes = magpieClientProperties.getServerNodes();
        if (serverNodes == null || serverNodes.length() == 0) {
            throw new MagpieClientInitializeException("Server nodes not defined");
        }

        MagpieClientConfig magpieClientConfig = new MagpieClientConfig();
        magpieClientConfig.setEnabled(enabled);
        magpieClientConfig.setServerNodes(parseServerNodesFromString(serverNodes));
        magpieClientConfig.setLoadBalancePolicy(LoadBalancePolicy.valueOf(loadBalancePolicy));
        return magpieClientConfig;
    }

    private List<InetSocketAddress> parseServerNodesFromString(String hostAndPortsString) {
        List<InetSocketAddress> result = new ArrayList<>();
        String[] split = hostAndPortsString.split(",");
        Arrays.stream(split).forEach(s -> {
            String[] hostAndPort = s.trim().split(":");
            if (hostAndPort.length != 2) {
                throw new RuntimeException("Wrong inet socket address: " + hostAndPort);
            }
            InetSocketAddress inetSocketAddress = new InetSocketAddress(hostAndPort[0], Integer.parseInt(hostAndPort[1]));
            result.add(inetSocketAddress);
        });
        return result;
    }

    @Bean
    @ConditionalOnMissingBean
    public MagpieClient magpieClient(MagpieClientConfig magpieClientConfig) throws MagpieClientInitializeException {
        return new MagpieClient(magpieClientConfig);
    }

}
