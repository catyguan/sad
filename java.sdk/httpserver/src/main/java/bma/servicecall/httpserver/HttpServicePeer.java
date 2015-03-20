package bma.servicecall.httpserver;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.http.HttpRequest;

import java.util.HashMap;
import java.util.Map;

import bma.servicecall.core.Address;
import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.BaseServiceServ;
import bma.servicecall.core.Context;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;
import bma.servicecall.core.ServiceRequest;
import bma.servicecall.core.StatusConst;
import bma.servicecall.core.Util;
import bma.servicecall.core.ValueMap;

public class HttpServicePeer implements ServicePeer {

	protected ServiceCallWebServer server;
	protected ChannelHandlerContext channelContext;
	protected HttpRequest httpRequest;
	protected Request qo;
	protected Context co;
	protected String transId;
	protected String asyncId;
	protected int mode; // 0-normal, 1-writed, 2-poll, 3-callback
	protected Address callback;

	@Override
	public String getDriverType() {
		return "http";
	}

	@Override
	public void beginTransaction() {
		if (!Util.empty(this.transId)) {
			throw new AppError("already begin transaction");
		}

		BaseServiceServ serv = this.server.serv;
		this.transId = serv.createSeq();

		this.server.setTrans(this.transId, this);
	}

	public void endTransaction() {
		if (Util.empty(this.transId)) {
			return;
		}
		this.server.setTrans(this.transId, null);
		this.transId = null;
	}

	@Override
	public ServiceRequest readRequest(int waitTimeMS) {
		ServiceRequest sr = new ServiceRequest();
		synchronized (this) {
			if (this.qo != null) {
				sr.setContext(this.co);
				sr.setRequest(this.qo);
				return sr;
			}
			try {
				this.wait(waitTimeMS);
			} catch (InterruptedException e) {
				throw AppError.handle(e);
			}
			if (this.qo != null) {
				sr.setContext(this.co);
				sr.setRequest(this.qo);
				this.mode = 0;
				return sr;
			} else {
				throw new AppError("timeout");
			}
		}
	}

	protected void reset() {
		synchronized (this) {
			this.channelContext = null;
			this.httpRequest = null;
			this.qo = null;
			this.co = null;
		}
	}

	@Override
	public void writeAnswer(Answer a, Exception err) {
		switch (this.mode) {
		case 2:
			if (Util.empty(this.asyncId)) {
				throw new AppError("poll mode, asyncId empty");
			}
			this.server.serv.setPollAnswer(this.asyncId, a, err);
			this.asyncId = "";
			return;
		case 3:
			Request cbreq;
			if (err != null) {
				Map<String, Object> am = new HashMap<String, Object>();
				am.put("Status", StatusConst.ERROR);
				am.put("Message", err.getMessage());
				cbreq = Request.create(am);
			} else {
				Map<String, Object> am = a.toMap();
				cbreq = Request.create(am);
			}
			Context ctx = new Context();
			if (Debuger.isEnable()) {
				Debuger.log("callback invoke -> " + err + ", " + a);
			}
			try {
				Answer an = this.server.serv.doCallback(this.callback, cbreq,
						ctx);
				if (Debuger.isEnable()) {
					Debuger.log("callback answer -> " + an);
				}
			} catch (Exception err2) {
				if (Debuger.isEnable()) {
					Debuger.log("callback fail -> " + err2);
				}
			}
			break;
		case 1:
			throw new AppError("HttpServicePeer already answer");
		default:
			if (this.channelContext == null) {
				throw new AppError("HttpServicePeer break");
			}
			if (!Util.empty(this.transId) && a.getStatus() != 100) {
				this.endTransaction();
			}
			if (Debuger.isEnable()) {
				Debuger.log("writeAnswer -> " + err + ", " + a);
			}
			try {
				ServiceCallWebServer.doAnswer(this, this.channelContext,
						this.httpRequest, a, err);
			} finally {
				this.mode = 1;
				this.reset();
			}
			break;
		}

	}

	@Override
	public void sendAsync(Context ctx, ValueMap result, int timeoutMS) {
		String async = ctx.getString(PropertyConst.ASYNC_MODE);
		if (async == null || async.length() == 0 || async.equals("poll")) {
			String aid = this.server.serv.createPollAnswer(timeoutMS, this);

			this.mode = 2;
			this.asyncId = aid;

			Answer a = new Answer();
			a.setStatus(StatusConst.ASYNC);
			if (result == null) {
				result = new ValueMap(null);
			}
			result.put(PropertyConst.ASYNC_ID, aid);
			a.setResult(result);
			try {
				ServiceCallWebServer.doAnswer(this, this.channelContext,
						this.httpRequest, a, null);
			} finally {
				this.reset();
			}
			return;
		}
		if (async.equals("callback")) {
			ValueMap addrm = ctx.getMap(PropertyConst.CALLBACK);
			if (addrm == null) {
				throw new AppError(
						"HttpServicePeer Async callback miss address");
			}
			this.callback = Address.createAddressFromValue(addrm);

			this.mode = 3;
			Answer a = new Answer();
			a.setStatus(StatusConst.ASYNC);
			a.setResult(result);
			try {
				ServiceCallWebServer.doAnswer(this, this.channelContext,
						this.httpRequest, a, null);
			} finally {
				this.reset();
			}
			return;
		}
		throw new AppError("HttpServicePeer not support AsyncMode(" + async
				+ ")");
	}

	public void post(ChannelHandlerContext ctx, HttpRequest req, Request qo,
			Context co) {
		synchronized (this) {
			if (this.channelContext != null) {
				throw new AppError("peer executing");
			}
			this.channelContext = ctx;
			this.httpRequest = req;
			this.qo = qo;
			this.co = co;
			this.notifyAll();
		}

	}

}
