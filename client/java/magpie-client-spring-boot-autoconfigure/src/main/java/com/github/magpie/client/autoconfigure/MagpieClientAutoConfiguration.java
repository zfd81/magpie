package com.github.magpie.client.autoconfigure;

import com.github.magpie.client.LoadBalancePolicy;
import com.github.magpie.client.MagpieClient;
import com.github.magpie.client.MagpieClientConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConditionalOnClass(MagpieClient.class)
@EnableConfigurationProperties(MagpieClientProperties.class)
public class MagpieClientAutoConfiguration {

    @Autowired
    private MagpieClientProperties magpieClientProperties;

    @Bean
    @ConditionalOnMissingBean
    public MagpieClientConfig magpieClientConfig() {

        Boolean enabled = magpieClientProperties.getEnabled();
        enabled = enabled == null || enabled;

        String loadBalancePolicy  = magpieClientProperties.getLoadBalancePolicy();
        loadBalancePolicy = loadBalancePolicy == null ? "round_robin" : loadBalancePolicy;

        Long timeout  = magpieClientProperties.getTimeout();
        timeout = timeout == null ? 5 * 60 * 1000 : timeout;

        String serverNodes = magpieClientProperties.getServerNodes();

        return MagpieClientConfig.newBuilder()
            .setEnabled(enabled)
            .setTimeout(timeout)
            .parseServerNodesFromString(serverNodes)
            .setLoadBalancePolicy(LoadBalancePolicy.valueOf(loadBalancePolicy))
            .build();
    }

    @Bean
    @ConditionalOnMissingBean
    public MagpieClient magpieClient(MagpieClientConfig magpieClientConfig) {
        return new MagpieClient(magpieClientConfig);
    }

}
