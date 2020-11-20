package com.github.magpie.client;

public class MagpieRpcException extends Exception {

    public MagpieRpcException() {
        super();
    }

    public MagpieRpcException(String message) {
        super(message);
    }

    public MagpieRpcException(String message, Throwable cause) {
        super(message, cause);
    }

    public MagpieRpcException(Throwable cause) {
        super(cause);
    }

    protected MagpieRpcException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}
