package bma.servicecall.core;

import java.util.LinkedList;
import java.util.List;

public class TestServicePeer implements ServicePeer {

	private String driverType;
	private List<ServiceRequest> nextRequests = new LinkedList<ServiceRequest>();

	public void addNextRequest(ServiceRequest req) {
		this.nextRequests.add(req);
	}

	public List<ServiceRequest> getNextRequests() {
		return nextRequests;
	}

	public void setNextRequests(List<ServiceRequest> nextRequests) {
		this.nextRequests = nextRequests;
	}

	public void setDriverType(String driverType) {
		this.driverType = driverType;
	}

	@Override
	public String getDriverType() {
		return this.driverType;
	}

	@Override
	public void beginTransaction() {
		// do nothing
	}

	@Override
	public ServiceRequest readRequest(int waitTimeMS) {
		if (this.nextRequests.isEmpty()) {
			throw new AppError("timeout");
		}
		return this.nextRequests.remove(0);
	}

	@Override
	public void writeAnswer(Answer a, Exception err) {
		// do nothing
	}

	@Override
	public void sendAsync(Context ctx, ValueMap result, int timeoutMS) {
		// do nothing
	}

}
