package bma.servicecall.core;

import java.util.Map;

public class Context extends ValueMap {

	public Context() {
		super(null);
	}

	public Context(Map<String, Value> d) {
		super(d);
	}

	@SuppressWarnings("rawtypes")
	public static Context create(Map d) {
		Context o = new Context();
		o.initValueMap(d);
		return o;
	}

	public String getSessionId() {
		return this.getString(PropertyConst.SESSION_ID);
	}
	
	public String getRequestId() {
		return this.getString(PropertyConst.REQ_ID);
	}
}
