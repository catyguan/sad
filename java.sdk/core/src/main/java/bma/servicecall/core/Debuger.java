package bma.servicecall.core;

public class Debuger {

	private static ServiceCallLogger log;

	public static void init(ServiceCallLogger l) {
		log = l;
	}

	public static void log(String msg) {
		if (log != null) {
			log.log(msg);
		}
	}
}
