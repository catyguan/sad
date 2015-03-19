package bma.servicecall.core;

import java.util.Map;
import java.util.TreeMap;

public class Answer {
	int status;
	String message;
	ValueMap result;
	ValueMap context;

	public Answer() {
		super();
	}

	public static Answer error2Answer(Answer a, Exception ex) {
		if (a == null) {
			a = new Answer();
		}
		a.setStatus(StatusConst.ERROR);
		a.setMessage(ex.getMessage());
		return a;
	}

	public int getStatus() {
		return status;
	}

	public void setStatus(int status) {
		this.status = status;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	public ValueMap getResult() {
		return result;
	}

	public void setResult(ValueMap result) {
		this.result = result;
	}

	public ValueMap sureResult() {
		if (this.result == null) {
			this.result = new ValueMap(null);
		}
		return this.result;
	}

	public ValueMap getContext() {
		return context;
	}

	public void setContext(ValueMap context) {
		this.context = context;
	}

	public ValueMap sureContext() {
		if (this.context == null) {
			this.context = new ValueMap(null);
		}
		return this.context;
	}

	public Map<String, Object> toMap() {
		Map<String, Object> m = new TreeMap<String, Object>();
		if (this.status != 0) {
			m.put("Status", this.status);
		}
		if (this.message != "") {
			m.put("Message", this.message);
		}
		if (this.result != null) {
			m.put("Result", this.result.toMap());
		}
		if (this.context != null) {
			m.put("Context", this.context.toMap());
		}
		return m;
	}

	public boolean isProcessing() {
		return this.isAsync();
	}

	public boolean isAsync() {
		int st = this.getStatus();
		switch (st) {
		case 202:
			return true;
		}
		return false;
	}

	public String getAsyncId() {
		String aid = "";
		ValueMap rs = this.getResult();
		if (rs != null) {
			aid = rs.getString(PropertyConst.ASYNC_ID);
		}
		return aid;
	}

	public boolean isContinue() {
		return this.getStatus() == 100;
	}

	public boolean isDone() {
		int st = this.getStatus();
		switch (st) {
		case 100:
		case 200:
		case 204:
			return true;
		}
		return false;
	}

	public void setSessionId(String v) {
		this.sureContext().put(PropertyConst.SESSION_ID, v);
	}

	public void checkError() {
		switch (this.status) {
		case 200:
		case 100:
		case 202:
		case 204:
		case 302:
			return;
		default:
			String msg = this.message;
			if (Util.empty(msg)) {
				msg = "unknow error";
			}
			throw new AppError(msg);
		}
	}

	@Override
	public String toString() {
		StringBuffer buf = new StringBuffer();
		if (this.status != 0) {
			buf.append("Status=").append(this.status).append(";");
		}
		if (this.message != "") {
			buf.append("Message=").append(this.message).append(";");
		}
		if (this.result != null) {
			buf.append("Result=").append(this.result).append(";");
		}
		if (this.context != null) {
			buf.append("Context=").append(this.context).append(";");
		}
		return buf.toString();
	}
}
