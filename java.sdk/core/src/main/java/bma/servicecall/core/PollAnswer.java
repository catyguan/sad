package bma.servicecall.core;

import java.util.Date;

public class PollAnswer {
	private boolean done;
	private ServicePeer peer;
	private Answer answer;
	private Exception err;
	private Date waitTime;

	public boolean isDone() {
		return done;
	}

	public void setDone(boolean done) {
		this.done = done;
	}

	public ServicePeer getPeer() {
		return peer;
	}

	public void setPeer(ServicePeer peer) {
		this.peer = peer;
	}

	public Answer getAnswer() {
		return answer;
	}

	public void setAnswer(Answer answer) {
		this.answer = answer;
	}

	public Exception getErr() {
		return err;
	}

	public void setErr(Exception err) {
		this.err = err;
	}

	public Date getWaitTime() {
		return waitTime;
	}

	public void setWaitTime(Date waitTime) {
		this.waitTime = waitTime;
	}

}
