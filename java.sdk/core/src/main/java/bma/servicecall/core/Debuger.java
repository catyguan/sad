package bma.servicecall.core;

public class Debuger {

	private static ServiceCallLogger log;

	public static void init(ServiceCallLogger l) {
		log = l;
	}

	public static boolean isEnable() {
		return log != null && log.isEnable();
	}

	public static void log(String msg) {
		if (log != null) {
			log.log(msg);
		}
	}
}
