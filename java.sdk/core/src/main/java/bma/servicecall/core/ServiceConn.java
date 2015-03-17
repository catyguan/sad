package bma.servicecall.core;

public interface ServiceConn {
	public Answer invoke(InvokeContext ictx, Address addr, Request req,
			Context ctx);

	public Answer waitAnswer(int timeoutMS);

	public void end();

	public void close();
}
