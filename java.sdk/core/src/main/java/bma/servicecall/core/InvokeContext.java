package bma.servicecall.core;

public interface InvokeContext {
	public void setProperty(String name, Object val);

	public Object getProperty(String name);

	public void removeProperty(String name);
}
