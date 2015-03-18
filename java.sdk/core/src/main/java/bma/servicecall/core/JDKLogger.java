package bma.servicecall.core;

import java.util.logging.Level;
import java.util.logging.Logger;

public class JDKLogger implements ServiceCallLogger {

	private Logger log;

	public JDKLogger(String name) {
		this(Logger.getLogger(name));
	}

	public JDKLogger(Logger l) {
		super();
		this.log = l;
	}

	public void log(String msg) {
		this.log.log(Level.INFO, msg);
	}

}
