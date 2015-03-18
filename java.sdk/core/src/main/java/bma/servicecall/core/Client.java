package bma.servicecall.core;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.atomic.AtomicInteger;

public class Client implements InvokeContext {

	private Manager manager;
	private String mname;
	private int id;
	private AtomicInteger reqId = new AtomicInteger();
	private Map<String, ServiceConn> conns;
	private String sessionId;
	private Map<String, Object> props;

	protected Client(Manager m, String name, int id) {
		super();
		this.manager = m;
		this.mname = name;
		this.id = id;
		this.conns = new HashMap<String, ServiceConn>();
	}

	public String createReqId() {
		int sid = this.reqId.addAndGet(1);
		if (sid <= 0) {
			this.reqId.compareAndSet(sid, 0);
			sid = this.reqId.addAndGet(1);
		}
		return this.mname + "_" + this.id + "_" + sid;
	}

	public String getSessionId() {
		return sessionId;
	}

	public void setSessionId(String sessionId) {
		this.sessionId = sessionId;
	}

	public void setProperty(String n, Object val) {
		if (this.props == null) {
			this.props = new HashMap<String, Object>();
		}
		this.props.put(n, val);
	}

	public Object getProperty(String n) {
		if (this.props == null) {
			return null;
		}
		return this.props.get(n);
	}

	public void removeProperty(String n) {
		if (this.props == null) {
			return;
		}
		this.props.remove(n);
	}

	protected ServiceConn getConn(Address addr, boolean create) {
		String api = addr.getApi();
		ServiceConn conn = this.conns.get(api);
		if (conn != null) {
			return conn;
		}
		if (!create) {
			return null;
		}
		String type = addr.getType();
		conn = this.manager.createConn(type, api);
		this.conns.put(api, conn);
		return conn;
	}

	protected void closeConn(Address addr) {
		String api = addr.getApi();
		ServiceConn conn = this.conns.remove(api);
		if (conn != null) {
			conn.close();
		}
	}

	protected Answer doInvoke(Address addr, Request req, Context ctx) {
		ServiceConn conn = this.getConn(addr, true);
		Answer a;
		try {
			a = conn.invoke(this, addr, req, ctx);
		} catch (Throwable t) {
			this.closeConn(addr);
			throw AppError.handle(t);
		}
		ValueMap actx = a.getContext();
		if (actx != null) {
			String sid = actx.getString(PropertyConst.SESSION_ID);
			if (sid != null && sid != "") {
				this.sessionId = sid;
			}
		}
		int st = a.getStatus();
		switch (st) {
		case 100:
		case 200:
		case 202:
		case 204:
		case 302:
			break;
		default:
			this.closeConn(addr);
		}
		return a;
	}

	public Answer invoke(Address addr, Request req, Context ctx) {
		if (ctx == null) {
			ctx = new Context();
		}
		if (!ctx.has(PropertyConst.DEADLINE)) {
			int to = ctx.getInt(PropertyConst.TIMEOUT);
			if (to <= 0) {
				to = 30;
			} else {
				ctx.remove(PropertyConst.TIMEOUT);
			}
			long dl = Util.currentUnixTimestamp() + to;
			ctx.put(PropertyConst.DEADLINE, dl);
		}
		if (!ctx.has(PropertyConst.REQ_ID)) {
			ctx.put(PropertyConst.REQ_ID, this.createReqId());
		}
		for (;;) {
			if (Util.empty(this.sessionId)) {
				ctx.remove(PropertyConst.SESSION_ID);
			} else {
				ctx.put(PropertyConst.SESSION_ID, this.sessionId);
			}
			Answer a = this.doInvoke(addr, req, ctx);
			switch (a.getStatus()) {
			case 200:
			case 100:
			case 202:
			case 204:
				return a;
			case 302:
				ValueMap rs = a.getResult();
				if (rs == null) {
					throw new AppError("redirect address empty");
				}
				addr = Address.createAddressFromValue(rs);
				addr.Valid();
				Debuger.log("redirect -> " + addr.toString());
			default:
				return a;
			}
		}
	}

	public void close() {
		if (this.props != null) {
			for (Entry<String, Object> e : this.props.entrySet()) {
				Object v = e.getValue();
				if (v != null) {
					if (v instanceof Closable) {
						Closable cl = (Closable) v;
						cl.close();
					}
				}
			}
			this.props.clear();
		}
		for (Entry<String, ServiceConn> e : this.conns.entrySet()) {
			e.getValue().end();
		}
	}

	public Map<String, Object> Export() {
		Map<String, Object> r = new HashMap<String, Object>();
		if (this.sessionId != null && this.sessionId != "") {
			r.put("SessionId", this.sessionId);
		}
		if (this.props != null) {
			Map<String, Object> m = new HashMap<String, Object>();
			for (Entry<String, Object> e : this.props.entrySet()) {
				m.put(e.getKey(), Util.scopyv(e.getValue()));
			}
			r.put("Props", m);
		}
		return r;
	}

	@SuppressWarnings("rawtypes")
	public void Import(Map<String, Object> data) {
		if (data == null) {
			return;
		}
		Object sv = data.get("SessionId");
		if (sv != null) {
			if (sv instanceof String) {
				String s = (String) sv;
				this.sessionId = s;
			}
		}
		Object mv = data.get("Props");
		if (mv != null) {
			if (mv instanceof Map) {
				Map m = (Map) mv;
				for (Object o : m.entrySet()) {
					Entry e = (Entry) o;
					if (this.props == null) {
						this.props = new HashMap<String, Object>();
					}
					Object vv = Util.scopyv(e.getValue());
					this.props.put(e.getKey().toString(), vv);
				}
			}
		}
		return;
	}

	public Answer PollAnswer(Address addr, Answer an, Context ctx,
			Date endTime, int sleepDurMS) {
		String aid = an.getAsyncId();
		if (Util.empty(aid)) {
			throw new AppError("miss AsyncId");
		}
		Request req = new Request();
		ctx.put(PropertyConst.ASYNC_ID, aid);
		for (;;) {
			if (new Date().after(endTime)) {
				return null;
			}
			Answer an2 = this.invoke(addr, req, ctx);
			if (!an2.isAsync()) {
				return an2;
			}
			if (new Date().after(endTime)) {
				return null;
			}
			if (sleepDurMS <= 0) {
				return null;
			}
			try {
				Thread.sleep(sleepDurMS);
			} catch (InterruptedException e) {
				return null;
			}
		}
	}

	public Answer waitAnswer(Address addr, int duMS) {
		ServiceConn conn = this.getConn(addr, false);
		if (conn == null) {
			throw new AppError("invalid connection for '" + addr.toString()
					+ "'");
		}
		return conn.waitAnswer(duMS);
	}
}
