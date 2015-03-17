package bma.servicecall.core;

public interface ServiceMethod {
	public void execute(ServicePeer peer, Request req, Context ctx);
}
