package bma.servicecall.core;

public interface StatusConst {
	public final int CONTINUE = 100;
	public final int DONE = 200;
	public final int ASYNC = 202;
	public final int OK = 204;
	public final int REDIRECT = 302;
	public final int INVALID = 400;
	public final int REJECT = 403;
	public final int TIMEOUT = 408;
	public final int ERROR = 500;
	public final int BADGATEWAY = 502;
}
