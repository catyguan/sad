package bma.servicecall.core;

public interface ServiceHandler {
	public void execute(ServicePeer peer, String service, String method,
			Request req, Context ctx);
}
