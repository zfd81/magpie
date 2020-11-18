package com.github.magpie.client.sample.server;

public class MagpieServerException extends RuntimeException {
    public MagpieServerException() {
    }

    public MagpieServerException(String message) {
        super(message);
    }

    public MagpieServerException(String message, Throwable cause) {
        super(message, cause);
    }

    public MagpieServerException(Throwable cause) {
        super(cause);
    }

    public MagpieServerException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}
