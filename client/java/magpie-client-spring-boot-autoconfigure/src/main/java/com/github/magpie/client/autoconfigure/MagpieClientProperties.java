package com.github.magpie.client.autoconfigure;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "magpie.client")
@Data
public class MagpieClientProperties {

    private Boolean enabled = true;
    private Long timeout;
    private String loadBalancePolicy;
    private String serverNodes;

}
