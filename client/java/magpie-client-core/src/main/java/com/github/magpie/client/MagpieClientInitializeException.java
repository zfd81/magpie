package com.github.magpie.client;

public class MagpieClientInitializeException extends Exception {
    public MagpieClientInitializeException() {
        super();
    }

    public MagpieClientInitializeException(String message) {
        super(message);
    }

    public MagpieClientInitializeException(String message, Throwable cause) {
        super(message, cause);
    }

    public MagpieClientInitializeException(Throwable cause) {
        super(cause);
    }

    protected MagpieClientInitializeException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}
