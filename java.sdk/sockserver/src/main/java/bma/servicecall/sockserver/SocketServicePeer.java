package bma.servicecall.sockserver;

import io.netty.channel.ChannelHandlerContext;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;

import bma.servicecall.core.Address;
import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;
import bma.servicecall.core.ServiceRequest;
import bma.servicecall.core.StatusConst;
import bma.servicecall.core.Util;
import bma.servicecall.core.ValueMap;
import bma.servicecall.sockcore.Message;

public class SocketServicePeer implements ServicePeer {

	protected ServiceCallSocketServer server;
	protected BlockingQueue<Message> wqueue;
	protected ChannelHandlerContext channelContext;
	protected Message msg;
	protected String asyncId;
	protected int mode; // 0-normal, 1-writed, 2-poll, 3-callback, 4-push
	protected Address callback;

	@Override
	public String getDriverType() {
		return "socket";
	}

	@Override
	public void beginTransaction() {

	}

	@Override
	public ServiceRequest readRequest(int waitTimeMS) {
		ServiceRequest sr = new ServiceRequest();
		Message qmsg;
		try {
			qmsg = this.wqueue.poll(waitTimeMS, TimeUnit.MILLISECONDS);
		} catch (InterruptedException e) {
			throw AppError.handle(e);
		}
		if (qmsg != null) {
			this.msg = qmsg;
			sr.setRequest(qmsg.getRequest());
			sr.setContext(qmsg.getContext());
			this.mode = 0;
			return sr;
		} else {
			throw new AppError("timeout");
		}
	}

	protected void reset() {
		this.msg = null;
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
			mode = 1;
			this.reset();
			return;
		case 4:
			if (Debuger.isEnable()) {
				Debuger.log("pushAnswer -> " + err + ", " + a);
			}
			ServiceCallSocketServer.doAnswer(this.channelContext,
					this.msg.getId(), a, err);
			if (a.getStatus() != StatusConst.ASYNC) {
				this.mode = 1;
				this.reset();
			}
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
			throw new AppError("SocketServicePeer already answer");
		default:
			if (this.channelContext == null) {
				throw new AppError("SocketServicePeer break");
			}
			if (Debuger.isEnable()) {
				Debuger.log("writeAnswer -> " + err + ", " + a);
			}
			try {
				ServiceCallSocketServer.doAnswer(this.channelContext,
						this.msg.getId(), a, err);
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
				ServiceCallSocketServer.doAnswer(this.channelContext,
						this.msg.getId(), a, null);
			} finally {
				this.reset();
			}
			return;
		}
		if (async.equals("push")) {
			this.mode = 4;
			Answer a = new Answer();
			a.setStatus(StatusConst.ASYNC);
			a.setResult(result);
			ServiceCallSocketServer.doAnswer(this.channelContext,
					this.msg.getId(), a, null);
			return;
		}
		if (async.equals("callback")) {
			ValueMap addrm = ctx.getMap(PropertyConst.CALLBACK);
			if (addrm == null) {
				throw new AppError(
						"SocketServicePeer Async callback miss address");
			}
			this.callback = Address.createAddressFromValue(addrm);

			this.mode = 3;
			Answer a = new Answer();
			a.setStatus(StatusConst.ASYNC);
			a.setResult(result);
			try {
				ServiceCallSocketServer.doAnswer(this.channelContext,
						this.msg.getId(), a, null);
			} finally {
				this.reset();
			}
			return;
		}
		throw new AppError("SocketServicePeer not support AsyncMode(" + async
				+ ")");
	}
}
