package com.github.magpie.client;

public enum LoadBalancePolicy {
    // 轮询
    round_robin,
    // 用 Go 语言实现的外部均衡负载服务解决方案
    grpclb,
    // 选择第一个可用的节点
    pick_first
}
