package bma.servicecall.core;

import java.util.logging.Level;
import java.util.logging.Logger;

public class JDKLogger implements ServiceCallLogger {

	private Logger log;
	private boolean e;

	public JDKLogger(String name) {
		this(Logger.getLogger(name));
	}

	public JDKLogger(Logger l) {
		super();
		this.log = l;
		this.e = l.isLoggable(Level.INFO);
	}

	public void log(String msg) {
		this.log.log(Level.INFO, msg);
	}

	@Override
	public boolean isEnable() {
		return this.e;
	}
}
