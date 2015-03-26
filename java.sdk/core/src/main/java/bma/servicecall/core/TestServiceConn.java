package bma.servicecall.core;

import java.util.LinkedList;
import java.util.List;

public class TestServiceConn implements Driver, ServiceConn {

	private List<Answer> answers = new LinkedList<Answer>();

	public List<Answer> getAnswers() {
		return answers;
	}

	public void setAnswers(List<Answer> answer) {
		this.answers = answer;
	}

	public void addAnswer(Answer an) {
		this.answers.add(an);
	}

	@Override
	public Answer invoke(InvokeContext ictx, Address addr, Request req,
			Context ctx) {
		if (this.answers.isEmpty()) {
			throw new AppError("timeout");
		}
		return this.answers.remove(0);
	}

	@Override
	public Answer waitAnswer(int timeoutMS) {
		if (this.answers.isEmpty()) {
			try {
				Thread.sleep(timeoutMS);
			} catch (InterruptedException e) {
			}
			return null;
		}
		return this.answers.remove(0);
	}

	@Override
	public void end() {

	}

	@Override
	public void close() {

	}

	@Override
	public ServiceConn createConn(String type, String api) {
		return this;
	}

	public void initDriver(String name) {
		Manager.initDriver(name, this);
	}
}
