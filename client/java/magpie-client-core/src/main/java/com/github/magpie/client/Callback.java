package com.github.magpie.client;

import lombok.Data;

@Data
public class Callback<T> {
    private T result;
}
