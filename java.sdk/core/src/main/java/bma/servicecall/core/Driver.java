package bma.servicecall.core;

public interface Driver {
	public ServiceConn createConn(String type, String api);
}
