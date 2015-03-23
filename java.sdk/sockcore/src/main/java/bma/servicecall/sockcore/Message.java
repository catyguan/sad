package bma.servicecall.sockcore;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;

public class Message {
	protected byte type;
	protected int id;
	protected boolean boolFlag;
	// Request
	protected String service;
	protected String method;
	protected Request request;
	protected Context context;
	// Answer
	protected Answer answer;

	public byte getType() {
		return type;
	}

	public void setType(byte type) {
		this.type = type;
	}

	public int getId() {
		return id;
	}

	public void setId(int id) {
		this.id = id;
	}

	public boolean isBoolFlag() {
		return boolFlag;
	}

	public void setBoolFlag(boolean boolFlag) {
		this.boolFlag = boolFlag;
	}

	public String getService() {
		return service;
	}

	public void setService(String service) {
		this.service = service;
	}

	public String getMethod() {
		return method;
	}

	public void setMethod(String method) {
		this.method = method;
	}

	public Request getRequest() {
		return request;
	}

	public void setRequest(Request request) {
		this.request = request;
	}

	public Context getContext() {
		return context;
	}

	public void setContext(Context context) {
		this.context = context;
	}

	public Answer getAnswer() {
		return answer;
	}

	public void setAnswer(Answer answer) {
		this.answer = answer;
	}

	public void reset() {
		this.type = 0;
		this.id = 0;
		this.boolFlag = false;
		this.service = null;
		this.method = null;
		this.request = null;
		this.context = null;
		this.answer = null;
	}
}
