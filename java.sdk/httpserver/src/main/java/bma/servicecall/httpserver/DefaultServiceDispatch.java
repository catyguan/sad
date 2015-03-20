package bma.servicecall.httpserver;

import io.netty.handler.codec.http.HttpRequest;

public class DefaultServiceDispatch implements ServiceDispatch {

	public static final DefaultServiceDispatch INSTANCE = new DefaultServiceDispatch();

	@Override
	public RequestTarget dispatch(HttpRequest req) {
		String uri = req.getUri();
		String[] splist = uri.split("/");
		int l = splist.length;
		RequestTarget r = new RequestTarget();
		switch (l) {
		case 0:
			r.setService("home");
			r.setMethod("index");
			break;
		case 1:
			r.setService(splist[0]);
			r.setMethod("index");
			break;
		default:
			r.setService(splist[l - 2]);
			r.setMethod(splist[l - 1]);
		}
		return r;
	}
}
