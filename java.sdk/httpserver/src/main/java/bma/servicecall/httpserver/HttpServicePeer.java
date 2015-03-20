package bma.servicecall.httpserver;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.http.HttpRequest;

import java.util.HashMap;
import java.util.Map;

import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.BaseServiceServ;
import bma.servicecall.core.Context;
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

	@Override
	public String getDriverType() {
		return "http";
	}

	@Override
	public void beginTransaction() {
		if (this.transId != "") {
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
		switch(this.mode) {
		case 2:
			if(Util.empty(this.asyncId)) {
				throw new AppError("poll mode, asyncId empty");
			}
			this.server.serv.setPollAnswer(this.asyncId, a, err);
			this.asyncId = "";
			return;
		case 3:
			Request cbreq;			
			if(err != null) {
				Map<String, Object> am = new HashMap<String, Object>();
				am.put("Status", StatusConst.ERROR);
				am.put("Message", err.getMessage());
				cbreq = Request.create(am);
			} else {
				Map<String, Object> am = a.toMap();
				cbreq = Request.create(am);
			}
			ctx := sccore.NewContext()
			sccore.DoLog("callback invoke -> %v, %v", err, a)
			an, err2 := this.mux.serv.DoCallback(this.callback, req, ctx)
			sccore.DoLog("callback answer -> %v, %v", err2, an)
			return err
		case 1:
			return fmt.Errorf("HttpServicePeer already answer")
		default:
			if(this.channelContext==null) {
				throw new AppError("HttpServicePeer break");
			}
			if this.transId != "" && a.GetStatus() != 100 {
				this.endTransaction()
			}
			sccore.DoLog("writeAnswer -> %v, %v", err, a)
			err2 := doAnswer(this, this.w, a, err)
			this.mode = 1
			close(this.end)
			return err2
		}
		
	}

	@Override
	public void sendAsync(Context ctx, ValueMap result, int timeoutMS) {
		// TODO Auto-generated method stub

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
