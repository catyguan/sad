package bma.servicecall.core;

public interface ServiceProvider {
	public ServiceMethod getServiceMethod(String service, String method);
}
