package bma.servicecall.sockclient;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Queue;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.locks.ReentrantReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock.ReadLock;
import java.util.concurrent.locks.ReentrantReadWriteLock.WriteLock;

import bma.servicecall.core.Address;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.Util;
import bma.servicecall.core.ValueMap;
import bma.servicecall.sockcore.Coder;
import bma.servicecall.sockcore.SocketAPI;

public class SocketPool {

	public static class Config {
		public String net;
		public String host;
		public int port;
		public int timeoutMS;
		public int poolSize;
		public int idleMS = 60 * 1000;

		public void valid() {
			if (Util.empty(this.net)) {
				this.net = "tcp";
			}
			if (!this.net.equals("tcp")) {
				throw new AppError("only support tcp socket");
			}
			if (Util.empty(this.host)) {
				throw new AppError("host empty");
			}
			if (this.port <= 0) {
				throw new AppError("invalid port");
			}
			if (this.poolSize < 0) {
				this.poolSize = 0;
			}
			if (this.idleMS <= 0) {
				this.idleMS = 60 * 1000;
			}
		}
	}

	public static class Item {
		public SocketConn conn;
		public Date waitTime;
		public Date pingTime;
	}

	public static class PoolObject {
		public Config config;
		public Queue<Item> wait = new LinkedBlockingQueue<SocketPool.Item>();
		public AtomicBoolean closed = new AtomicBoolean();
	}

	private SocketConn poGetSocket(PoolObject po, int timeoutMS) {
		Item item = po.wait.poll();
		if (item != null) {
			return item.conn;
		}
		SocketConn conn = poDial(po, timeoutMS);
		return conn;
	}

	private boolean poPut(PoolObject po, Item item) {
		if (!po.closed.get() && po.wait.size() < po.config.poolSize) {
			po.wait.add(item);
			return true;
		} else {
			item.conn.close();
			return false;
		}
	}

	private void poReturnSocket(PoolObject po, SocketConn conn) {
		Date now = new Date();
		Item item = new Item();
		item.conn = conn;
		item.waitTime = now;
		item.pingTime = now;
		// fmt.Println("ReturnConnect", conn.LocalAddr())
		if (!this.poPut(po, item)) {
			// fmt.Println("ReturnFail")
		}
	}

	private SocketConn poDial(PoolObject po, int timeoutMS) {
		if (po.config.timeoutMS > 0) {
			timeoutMS = po.config.timeoutMS;
		}
		if (timeoutMS == 0) {
			timeoutMS = 5 * 1000;
		}
		Socket sock = new Socket();
		try {
			sock.connect(new InetSocketAddress(po.config.host, po.config.port),
					timeoutMS);
			SocketConn conn = new SocketConn();
			conn.setSocket(sock);
			return conn;
		} catch (IOException e) {
			throw AppError.handle(e);
		}
	}

	private void poClose(PoolObject po) {
		po.closed.set(true);
		while (true) {
			Item item = po.wait.poll();
			if (item == null) {
				break;
			}
			item.conn.close();
		}
	}

	private void poIdlePing(PoolObject po) {
		int idleDu = po.config.idleMS;
		int pingDu = 15 * 1000;

		int l = po.wait.size();
		for (int i = 0; i < l; i++) {
			if (po.closed.get()) {
				break;
			}
			Item item = po.wait.poll();
			if (item == null) {
				break;
			}
			Date now = new Date();
			if (now.getTime() - item.waitTime.getTime() > idleDu) {
				if (Debuger.isEnable()) {
					Debuger.log("'" + item.conn + "' idle break");
				}
				item.conn.close();
				continue;
			} else {
				if (now.getTime() - item.pingTime.getTime() > pingDu) {
					// ping
					if (poPing(po, item.conn)) {
						item.pingTime = now;
						poPut(po, item);
					}
				} else {
					poPut(po, item);
				}
			}
		}
	}

	private static final byte[] pingData = { 9, 0, 0, 1, 0, 0, 0, 0, 0 };
	private static final byte[] pingRData = { 9, 0, 0, 1, 1, 0, 0, 0, 0 };

	private boolean poPing(PoolObject po, SocketConn conn) {
		// fmt.Println("do ping")
		try {
			conn.getSocket().setSoTimeout(5 * 1000);
			conn.getOut().write(pingData);
			byte[] bs = new byte[pingRData.length];
			Coder.readAll(conn.getIn(), bs);
			for (int i = 0; i < bs.length; i++) {
				if (bs[i] != pingRData[i]) {
					if (Debuger.isEnable()) {
						Debuger.log("'" + conn.getSocket()
								+ "' ping invalid response " + i + "/" + bs[i]);
					}
					return false;
				}
			}
			return true;
		} catch (Exception e) {
			if (Debuger.isEnable()) {
				Debuger.log("'" + conn.getSocket() + "' ping fail " + e);
			}
			return false;
		}
	}

	public static SocketPool gPool = new SocketPool();

	public static SocketPool pool() {
		return gPool;
	}

	private ReentrantReadWriteLock locker = new ReentrantReadWriteLock();
	private ReadLock rlock = locker.readLock();
	private WriteLock wlock = locker.writeLock();
	private Map<String, PoolObject> pools = new HashMap<String, PoolObject>();
	protected Config config = new Config();
	private Timer timer;

	public SocketConn getSocket(Address addr, SocketAPI api, long timeout) {
		if (api == null) {
			api = SocketAPI.parseSocketAPI(addr.getApi());
		}
		api.valid();

		PoolObject po;
		String key = api.key();
		this.rlock.lock();
		try {
			po = this.pools.get(key);
		} finally {
			this.rlock.unlock();
		}
		if (po == null) {
			String host = api.getHost();
			int port = api.getPort();
			Config cfg = new Config();
			cfg.net = api.getType();
			cfg.host = host;
			cfg.port = port;
			cfg.poolSize = this.config.poolSize;
			cfg.idleMS = this.config.idleMS;
			cfg.timeoutMS = this.config.timeoutMS;
			ValueMap opt = addr.getOption();
			if (opt != null) {
				if (opt.has("PoolSize")) {
					cfg.poolSize = opt.getInt("PoolSize");
				}
				if (opt.has("Timeout")) {
					cfg.timeoutMS = opt.getInt("Timeout");
				}
				if (opt.has("Idle")) {
					cfg.idleMS = opt.getInt("Idle");
				}
			}

			this.wlock.lock();
			try {
				po = this.pools.get(key);
				if (po == null) {
					po = new PoolObject();
					po.config = cfg;
					this.pools.put(key, po);
					if (this.timer == null) {
						this.timer = new Timer();
					}
					final PoolObject fpo = po;
					this.timer.schedule(new TimerTask() {

						@Override
						public void run() {
							poIdlePing(fpo);
						}
					}, 1000, 1000);
				}
			} finally {
				this.wlock.unlock();
			}
		}
		if (timeout <= 0) {
			timeout = this.config.timeoutMS;
		}
		SocketConn conn = poGetSocket(po, (int) timeout);
		conn.setKey(key);
		return conn;
	}

	public void returnSocket(SocketConn conn) {
		this.rlock.lock();
		try {
			PoolObject po = this.pools.get(conn.getKey());
			if (po != null) {
				poReturnSocket(po, conn);
			} else {
				conn.close();
			}
		} finally {
			this.rlock.unlock();
		}
	}

	public void closeSocket(SocketConn conn) {
		conn.close();
	}

	public void close() {
		this.wlock.lock();
		try {
			for (Entry<String, PoolObject> e : this.pools.entrySet()) {
				poClose(e.getValue());
			}
			this.pools.clear();
		} finally {
			this.wlock.unlock();
		}
	}
}
