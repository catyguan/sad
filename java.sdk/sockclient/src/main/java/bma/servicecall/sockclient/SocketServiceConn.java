package bma.servicecall.sockclient;

import java.util.Date;
import java.util.concurrent.atomic.AtomicInteger;

import bma.servicecall.core.Address;
import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.InvokeContext;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServiceConn;
import bma.servicecall.sockcore.Message;
import bma.servicecall.sockcore.MessageReader;
import bma.servicecall.sockcore.MessageWriter;
import bma.servicecall.sockcore.SocketAPI;
import bma.servicecall.sockcore.SocketCoreConst;

public class SocketServiceConn implements ServiceConn {

	private static final AtomicInteger gMessageId = new AtomicInteger();

	private SocketConn conn;

	public Answer invoke(InvokeContext ictx, Address addr, Request req,
			Context ctx) {
		try {
			return doInvoke(ictx, addr, req, ctx);
		} catch (Exception ex) {
			throw AppError.handle(ex);
		}
	}

	public Answer doInvoke(InvokeContext ictx, Address addr, Request req,
			Context ctx) throws Exception {
		SocketAPI sapi = SocketAPI.parseSocketAPI(addr.getApi());
		SocketConn conn = null;
		if (this.conn != null) {
			conn = this.conn;
		}

		long dltm = ctx.getLong(PropertyConst.DEADLINE);
		long du = dltm * 1000 - new Date().getTime();
		if (du <= 0) {
			throw new AppError("timeout");
		}
		if (conn == null) {
			if (Debuger.isEnable()) {
				Debuger.log("'" + sapi + "' connect...");
			}
			conn = SocketPool.pool().getSocket(addr, sapi, du);
			if (conn != null) {
				this.conn = conn;
			}
		} else {
			if (Debuger.isEnable()) {
				Debuger.log("'" + sapi + "' use trans socket");
			}
		}
		try {
			conn.getSocket().setSoTimeout((int) du);

			// opt := addr.GetOption()
			MessageWriter mw = new MessageWriter(conn.getOut());

			if (Debuger.isEnable()) {
				Debuger.log("'" + sapi + "' write request to '"
						+ conn.getSocket().getRemoteSocketAddress() + "'");
			}
			int mid = gMessageId.addAndGet(1);
			if (mid <= 0) {
				gMessageId.compareAndSet(1, 0);
				mid = gMessageId.addAndGet(1);
			}
			mw.sendRequest(mid, sapi.getService(), sapi.getMethod(), req, ctx);
			MessageReader mr = new MessageReader(this.conn.getIn());
			Message msg = new Message();
			while (true) {
				byte mt = mr.nextMessage(msg);
				switch (mt) {
				case SocketCoreConst.MT_ANSWER: {
					Answer a = msg.getAnswer();
					switch (a.getStatus()) {
					case 100:
						if (Debuger.isEnable()) {
							Debuger.log("keep connection for transaction");
						}
						break;
					case 202:
						String amode = ctx.getString(PropertyConst.ASYNC_MODE);
						if (amode != null && amode.equals("callback")) {
							this.ret();
						} else {
							if (Debuger.isEnable()) {
								Debuger.log("keep connection for async");
							}
						}
						break;
					case 200:
					case 204:
					case 302:
						this.ret();
						break;
					default:
						this.close();
						break;
					}
					return a;
				}
				default:
					if (Debuger.isEnable()) {
						Debuger.log("unknow message(" + mt + ") - (" + msg
								+ ")");
					}
					break;
				}
			}
		} catch (Exception e) {
			if (Debuger.isEnable()) {
				Debuger.log("'" + sapi + "' fail '" + e + "'");
			}
			this.conn.close();
			this.conn = null;
			throw e;
		}
	}

	@Override
	public Answer waitAnswer(int timeoutMS) {
		SocketConn conn = this.conn;
		if (conn == null) {
			throw new AppError("invalid connection for ServerPush");
		}

		try {
			conn.getSocket().setSoTimeout(timeoutMS);

			MessageReader mr = new MessageReader(conn.getIn());
			Message msg = new Message();
			while (true) {
				byte mt = mr.nextMessage(msg);

				Answer an = msg.getAnswer();
				switch (mt) {
				case SocketCoreConst.MT_ANSWER:
					switch (an.getStatus()) {
					case 202:
						break;
					case 200:
					case 204:
						this.ret();
					default:
						this.close();
					}
					return an;
				default:
					if (Debuger.isEnable()) {
						Debuger.log("unknow message(" + mt + ") - (" + msg
								+ ")");
					}
				}
			}
		} catch (Exception e) {
			this.close();
			throw AppError.handle(e);
		}
	}

	protected void ret() {
		if (this.conn != null) {
			if (Debuger.isEnable()) {
				Debuger.log("return conn "
						+ this.conn.getSocket().getLocalAddress());
			}
			SocketPool.pool().returnSocket(this.conn);
			this.conn = null;
		}
	}

	@Override
	public void end() {
		this.close();
	}

	@Override
	public void close() {
		if (this.conn != null) {
			SocketPool.pool().closeSocket(this.conn);
			this.conn = null;
		}
	}

}
