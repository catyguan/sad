package bma.servicecall.core;

import java.util.Map;

public class Request extends ValueMap {
	
	public Request() {
		super(null);
	}

	public Request(Map<String, Value> d) {
		super(d);
	}
	
	@SuppressWarnings("rawtypes")
	public static Request create(Map d) {
		Request o = new Request();
		o.initValueMap(d);
		return o;
	}

}
