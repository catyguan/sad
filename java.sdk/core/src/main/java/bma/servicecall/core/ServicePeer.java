package bma.servicecall.core;

public interface ServicePeer {
	public String getDriverType();

	public void beginTransaction();

	public ServiceRequest readRequest(int waitTimeMS);

	public void writeAnswer(Answer a, Exception err);

	public void sendAsync(Context ctx, ValueMap result, int timeoutMS);
}
