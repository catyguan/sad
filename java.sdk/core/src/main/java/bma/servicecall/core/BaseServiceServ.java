package bma.servicecall.core;

import java.util.ArrayList;
import java.util.Date;
import java.util.LinkedHashMap;
import java.util.Map.Entry;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class BaseServiceServ {
	private ClientFactory clientFactory;
	private ReentrantReadWriteLock lock = new ReentrantReadWriteLock();
	public Lock rlock = lock.readLock();
	public Lock wlock = lock.writeLock();
	private LinkedHashMap<String, PollAnswer> polls;
	private long seed;
	private AtomicInteger seq = new AtomicInteger();

	public BaseServiceServ() {
		super();
	}

	public BaseServiceServ(ClientFactory cl) {
		super();
		this.clientFactory = cl;
	}

	public void setClientFactory(ClientFactory cl) {
		this.clientFactory = cl;
	}

	public String createSeq() {
		int s = this.seq.addAndGet(1);
		if (s <= 0) {
			this.seq.compareAndSet(s, 0);
			s = this.seq.addAndGet(1);
		}
		String k = this.seed + "_" + s;
		return Util.md5(k.getBytes());
	}

	public void checkPollAnswerTimeout() {
		Date now = new Date();
		ArrayList<String> outs = new ArrayList<String>();
		this.rlock.lock();

		try {
			if (this.polls != null) {
				for (Entry<String, PollAnswer> e : this.polls.entrySet()) {
					PollAnswer pa = e.getValue();
					if (now.after(pa.getWaitTime())) {
						outs.add(e.getKey());
					}
				}
			}
		} finally {
			this.rlock.unlock();
		}

		if (!outs.isEmpty()) {
			for (String key : outs) {
				this.wlock.lock();
				try {
					if (this.polls.remove(key) != null) {
						if (Debuger.isEnable()) {
							Debuger.log("pollAnswer(" + key + ")timeout");
						}
					}
				} finally {
					this.wlock.unlock();
				}
			}
		}
	}

	public String createPollAnswer(int duMS, ServicePeer peer) {
		String aid = this.createSeq();
		PollAnswer pa = new PollAnswer();
		pa.setPeer(peer);
		pa.setWaitTime(new Date(new Date().getTime() + duMS));
		this.wlock.lock();
		try {
			if (this.polls == null) {
				this.polls = new LinkedHashMap<String, PollAnswer>();
			}
			this.polls.put(aid, pa);
			return aid;
		} finally {
			this.wlock.unlock();
		}
	}

	public void setPollAnswer(String aid, Answer an, Exception err) {
		this.wlock.lock();
		try {
			if (this.polls != null) {
				PollAnswer pa = this.polls.get(aid);
				if (pa == null) {
					if (Debuger.isEnable()) {
						Debuger.log("pollAnswer(" + aid + ") invalid");
					}
					return;
				}
				pa.setDone(true);
				pa.setAnswer(an);
				pa.setErr(err);
				Debuger.log("poll async answer '" + aid + "'");
				return;
			}
			Debuger.log("poll async miss '" + aid + "'");
		} finally {
			this.wlock.unlock();
		}
	}

	public PollAnswer pollAsync(String aid) {
		if (Util.empty(aid)) {
			return null;
		}
		PollAnswer pa = null;
		this.rlock.lock();
		try {
			if (this.polls != null) {
				PollAnswer pa2 = this.polls.get(aid);
				if (pa2.isDone()) {
					pa = pa2;
				}
			}
		} finally {
			this.rlock.unlock();
		}

		if (pa != null) {
			this.wlock.lock();
			try {
				this.polls.remove(aid);
			} finally {
				this.wlock.unlock();
			}
			Debuger.log("'" + aid + "' poll success");
			// pa.Timer.Stop()
			return pa;
		} else {
			Debuger.log("'" + aid + "' polling");
			return null;
		}
	}

	public Answer doCallback(Address addr, Request req, Context ctx) {
		if (this.clientFactory == null) {
			throw new AppError("clientFactory is null");
		}
		Client cl = this.clientFactory.createClient();
		try {
			Answer answer = cl.invoke(addr, req, ctx);
			return answer;
		} finally {
			cl.close();
		}
	}
}
