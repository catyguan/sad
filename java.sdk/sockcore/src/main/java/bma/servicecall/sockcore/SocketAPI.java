package bma.servicecall.sockcore;

import bma.servicecall.core.AppError;
import bma.servicecall.core.Util;

public class SocketAPI {
	private String type = "tcp";
	private String host;
	private int port;
	private String service;
	private String method;

	public SocketAPI() {
		super();
	}

	public String getType() {
		return type;
	}

	public void setType(String type) {
		this.type = type;
	}

	public String getHost() {
		return host;
	}

	public void setHost(String host) {
		this.host = host;
	}

	public int getPort() {
		return port;
	}

	public void setPort(int port) {
		this.port = port;
	}

	public String getService() {
		return service;
	}

	public void setService(String service) {
		this.service = service;
	}

	public String getMethod() {
		return method;
	}

	public void setMethod(String method) {
		this.method = method;
	}

	public void valid() {
		if (Util.empty(this.type)) {
			throw new AppError("Type empty");
		}
		if (this.port < 0) {
			throw new AppError("Port(" + this.port + ") invalid");
		}
		if (Util.empty(this.service)) {
			throw new AppError("Service empty");
		}
		if (Util.empty(this.method)) {
			throw new AppError("Method empty");
		}
	}

	public String key() {
		StringBuffer buf = new StringBuffer();
		if (this.type != null) {
			buf.append(this.type);
		}
		buf.append(':');
		if (this.host != null) {
			buf.append(this.host);
		}
		buf.append(':');
		buf.append(this.port);
		return buf.toString();
	}

	@Override
	public String toString() {
		StringBuffer buf = new StringBuffer();
		if (this.type != null) {
			buf.append(this.type);
		}
		buf.append(':');
		if (this.host != null) {
			buf.append(this.host);
		}
		buf.append(':');
		buf.append(this.port);
		buf.append(':');
		if (this.service != null) {
			buf.append(this.service);
		}
		buf.append(':');
		if (this.method != null) {
			buf.append(this.method);
		}
		return buf.toString();
	}

	public static SocketAPI parseSocketAPI(String s) {
		String[] ps = s.split(":", 5);
		if (ps.length != 5) {
			throw new AppError("invalid SocketAPI - " + s);
		}
		SocketAPI o = new SocketAPI();
		o.type = ps[0];
		o.host = ps[1];
		if (!Util.empty(ps[2])) {
			o.port = Integer.parseInt(ps[2]);
		}
		o.service = ps[3];
		o.method = ps[4];
		o.valid();
		return o;
	}

}
