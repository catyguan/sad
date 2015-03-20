package bma.servicecall.httpserver;

import io.netty.handler.codec.http.HttpRequest;

public interface ServiceDispatch {

	public RequestTarget dispatch(HttpRequest req);
}
