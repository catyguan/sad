package bma.servicecall.core;

import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

public class Client {

	private Manager manager;
	private String mname;
	private int id;
	private AtomicInteger reqId = new AtomicInteger();
	private Map<String, ServiceConn> conns;
	private String sessionId;
	private Map<String, Object> props;

	protected Client(Manager m, String name, int id) {
		super();
		this.manager = m;
		this.mname = name;
		this.id = id;
	}

	public String createReqId() {
		int sid = this.reqId.addAndGet(1);
		if (sid <= 0) {
			sid = this.reqId.addAndGet(1);
		}
		return this.mname + "_" + this.id + "_" + sid;
	}

	public String getSessionId() {
		return sessionId;
	}

	public void setSessionId(String sessionId) {
		this.sessionId = sessionId;
	}

}
