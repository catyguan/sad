package bma.servicecall.core;

import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Manager {

	private final static Map<String, Driver> gDS = new TreeMap<String, Driver>();

	public static Driver getDriver(String type) {
		return gDS.get(type);
	}

	public static void initDriver(String type, Driver df) {
		gDS.put(type, df);
	}

	private String name;
	private AtomicInteger clientSeq = new AtomicInteger();

	public Manager(String n) {
		super();
		if (Util.empty(n)) {
			n = "jvscm" + System.currentTimeMillis();
		}
		this.name = n;
	}

	public Client createClient() {
		int id = this.clientSeq.addAndGet(1);
		if (id <= 0) {
			id = this.clientSeq.addAndGet(1);
		}
		return new Client(this, this.name, id);
	}

	protected ServiceConn createConn(String type, String api) {
		Driver df = getDriver(type);
		if (df == null) {
			throw new AppError("unknow driver(" + type + ")");
		}
		return df.createConn(type, api);
	}
}
