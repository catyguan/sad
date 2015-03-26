package bma.servicecall.httpserver;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpResponseEncoder;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.multipart.Attribute;
import io.netty.handler.codec.http.multipart.DefaultHttpDataFactory;
import io.netty.handler.codec.http.multipart.HttpPostRequestDecoder;
import io.netty.handler.codec.http.multipart.InterfaceHttpData;
import io.netty.handler.codec.http.multipart.InterfaceHttpData.HttpDataType;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.util.SelfSignedCertificate;
import io.netty.util.concurrent.DefaultEventExecutorGroup;

import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.TreeMap;

import org.codehaus.jackson.JsonParser;
import org.codehaus.jackson.map.ObjectMapper;

import bma.servicecall.boot.ServerBooter;
import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.BaseServiceServ;
import bma.servicecall.core.ClientFactory;
import bma.servicecall.core.Context;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.PollAnswer;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServiceMethod;
import bma.servicecall.core.ServiceMux;
import bma.servicecall.core.StatusConst;
import bma.servicecall.core.Util;

public class ServiceCallWebServer implements ServerBooter {

	private static final ObjectMapper mapper;

	static {
		mapper = new ObjectMapper();
		mapper.configure(JsonParser.Feature.ALLOW_SINGLE_QUOTES, true);
		mapper.configure(JsonParser.Feature.ALLOW_COMMENTS, true);
		mapper.configure(JsonParser.Feature.ALLOW_UNQUOTED_FIELD_NAMES, true);
	}

	public static ObjectMapper getDefaultMapper() {
		return mapper;
	}

	private EventLoopGroup bossGroup;
	private EventLoopGroup workerGroup;
	private DefaultEventExecutorGroup executorGroup;
	private Channel listener;
	private Map<String, HttpServicePeer> trans = new TreeMap<String, HttpServicePeer>();
	private Timer timer;

	private int executors = 10;
	private boolean ssl;
	private int port;
	private int maxContentLength = 10 * 1024 * 1024;
	private ServiceDispatch dispatch = DefaultServiceDispatch.INSTANCE;
	private boolean debugLog;
	private ServiceMux serviceMux;
	protected BaseServiceServ serv = new BaseServiceServ();

	public HttpServicePeer getTrans(String tid) {
		this.serv.rlock.lock();
		try {
			return this.trans.get(tid);
		} finally {
			this.serv.rlock.unlock();
		}
	}

	public void setTrans(String tid, HttpServicePeer peer) {
		this.serv.wlock.lock();
		try {
			if (peer == null) {
				this.trans.remove(tid);
			} else {
				this.trans.put(tid, peer);
			}
		} finally {
			this.serv.wlock.unlock();
		}
	}

	public int getExecutors() {
		return executors;
	}

	public void setExecutors(int executors) {
		this.executors = executors;
	}

	public void setClientFactory(ClientFactory cl) {
		serv.setClientFactory(cl);
	}

	public ServiceMux getServiceMux() {
		return serviceMux;
	}

	public void setServiceMux(ServiceMux serviceMux) {
		this.serviceMux = serviceMux;
	}

	public boolean isLog() {
		return debugLog;
	}

	public void setLog(boolean log) {
		this.debugLog = log;
	}

	public ServiceDispatch getDispatch() {
		return dispatch;
	}

	public void setDispatch(ServiceDispatch dispatch) {
		this.dispatch = dispatch;
	}

	public boolean isSsl() {
		return ssl;
	}

	public void setSsl(boolean ssl) {
		this.ssl = ssl;
	}

	public int getPort() {
		return port;
	}

	public void setPort(int port) {
		this.port = port;
	}

	public int getMaxContentLength() {
		return maxContentLength;
	}

	public void setMaxContentLength(int maxContentLength) {
		this.maxContentLength = maxContentLength;
	}

	public void start() throws Exception {
		final SslContext sslCtx;
		if (ssl) {
			SelfSignedCertificate ssc = new SelfSignedCertificate();
			sslCtx = SslContext.newServerContext(ssc.certificate(),
					ssc.privateKey());
		} else {
			sslCtx = null;
		}
		// Configure the server.
		bossGroup = new NioEventLoopGroup(1);
		workerGroup = new NioEventLoopGroup();
		executorGroup = new DefaultEventExecutorGroup(this.executors);

		ServerBootstrap b = new ServerBootstrap();
		b.group(bossGroup, workerGroup).channel(NioServerSocketChannel.class);
		if (this.debugLog) {
			b.handler(new LoggingHandler(LogLevel.INFO));
		}
		b.childHandler(new ChannelInitializer<SocketChannel>() {
			@Override
			protected void initChannel(SocketChannel ch) throws Exception {
				ChannelPipeline p = ch.pipeline();
				if (sslCtx != null) {
					p.addLast(sslCtx.newHandler(ch.alloc()));
				}
				p.addLast(new HttpRequestDecoder());
				p.addLast(new HttpObjectAggregator(maxContentLength));
				p.addLast(new HttpResponseEncoder());
				p.addLast(new ServiceCallServerHandler());
			}
		});

		Debuger.log("listen at " + this.port);
		this.listener = b.bind(port).sync().channel();

		this.timer = new Timer(true);
		this.timer.schedule(new TimerTask() {

			@Override
			public void run() {
				serv.checkPollAnswerTimeout();
			}
		}, 1000, 1000);
	}

	public void waitClose() throws InterruptedException {
		if (this.listener != null) {
			this.listener.closeFuture().sync();
		}
	}

	public void close() {
		if (this.timer != null) {
			this.timer.cancel();
			this.timer = null;
		}
		if(this.listener!=null) {
			this.listener.close();
			this.listener = null;
		}
		if (this.bossGroup != null) {
			this.bossGroup.shutdownGracefully();
			this.bossGroup = null;
		}
		if (this.workerGroup != null) {
			this.workerGroup.shutdownGracefully();
			this.workerGroup = null;
		}
		if (this.executorGroup != null) {
			this.executorGroup.shutdownGracefully();
			this.executorGroup = null;
		}
	}

	@Override
	public void startServer() {
		try {
			this.start();
		} catch (Exception e) {
			throw AppError.handle(e);
		}
	}

	@Override
	public void stopServer() {
		this.close();
	}

	public static void doAnswer(HttpServicePeer peer,
			ChannelHandlerContext ctx, HttpRequest req, Answer aa,
			Exception aerr) {
		Answer a = Answer.error2Answer(aa, aerr);
		if (a.getStatus() == 0) {
			a.setStatus(200);
		}
		if (peer != null && !Util.empty(peer.transId)) {
			a.sureContext().put(PropertyConst.TRANSACTION_ID, peer.transId);
		}
		Map<String, Object> m = a.toMap();

		ObjectMapper om = getDefaultMapper();
		String content;
		try {
			content = om.writeValueAsString(m);
		} catch (Exception e) {
			throw AppError.handle(e);
		}

		HttpUtil.sendResponse(ctx, req, HttpResponseStatus.OK, content);
	}

	public class ServiceCallServerHandler extends
			SimpleChannelInboundHandler<HttpRequest> {
		@Override
		public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
				throws Exception {
			if (!ctx.channel().isOpen()) {
				return;
			}
			if (Debuger.isEnable()) {
				Debuger.log("exceptionCaught - " + cause);
//				cause.printStackTrace();
			}
			String msg = cause.getMessage();
			if (debugLog) {
				StringWriter sw = new StringWriter();
				PrintWriter pw = new PrintWriter(sw);
				cause.printStackTrace(pw);
				msg = sw.toString();
			}
			HttpUtil.sendError(ctx, null, null, msg, cause);
			// super.exceptionCaught(ctx, cause);
		}

		@SuppressWarnings("rawtypes")
		@Override
		protected void channelRead0(ChannelHandlerContext ctx, HttpRequest req)
				throws Exception {
			if (!req.getDecoderResult().isSuccess()) {
				HttpUtil.sendResponse(ctx, req, HttpResponseStatus.BAD_REQUEST,
						"Bad Request");
				return;
			}
			if ("/favicon.ico".equals(req.getUri())) {
				HttpUtil.sendResponse(ctx, req, HttpResponseStatus.NOT_FOUND,
						"not found");
				return;
			}
			RequestTarget rt = dispatch.dispatch(req);
			if (Debuger.isEnable()) {
				Debuger.log("http '" + req.getUri() + "' -> " + rt.getService()
						+ ":" + rt.getMethod());
			}
			if (serviceMux == null) {
				throw new AppError("serviceMux invalid");
			}
			ServiceMethod sm = serviceMux.find(rt.getService(), rt.getMethod());
			if (sm == null) {
				HttpUtil.sendResponse(ctx, req, HttpResponseStatus.NOT_FOUND,
						"service(" + rt.getService() + ":" + rt.getMethod()
								+ ") not found");
				return;
			}
			HttpPostRequestDecoder decoder = new HttpPostRequestDecoder(
					new DefaultHttpDataFactory(false), req);
			String qv = "";
			InterfaceHttpData qdata = decoder.getBodyHttpData("q");
			if (qdata != null
					&& qdata.getHttpDataType() == HttpDataType.Attribute) {
				Attribute attribute = (Attribute) qdata;
				qv = attribute.getValue();
			}
			String cv = "";
			InterfaceHttpData cdata = decoder.getBodyHttpData("c");
			if (qdata != null
					&& qdata.getHttpDataType() == HttpDataType.Attribute) {
				Attribute attribute = (Attribute) cdata;
				cv = attribute.getValue();
			}

			ObjectMapper om = getDefaultMapper();

			if (Debuger.isEnable()) {
				Debuger.log("Q : " + qv);
				Debuger.log("C : " + cv);
			}

			Map qm = null;
			if (!Util.empty(qv)) {
				qm = om.readValue(qv, Map.class);
			}
			Map cm = null;
			if (!Util.empty(cv)) {
				cm = om.readValue(cv, Map.class);
			}

			if (Debuger.isEnable()) {
				Debuger.log("QM : " + qm);
				Debuger.log("CM : " + cm);
			}

			Request qo = Request.create(qm);
			Context co = Context.create(cm);

			String aid = co.getString(PropertyConst.ASYNC_ID);
			if (!Util.empty(aid)) {
				PollAnswer pa = serv.pollAsync(aid);
				HttpServicePeer peer = null;
				Answer aa;
				Exception aerr = null;
				if (pa != null) {
					aa = pa.getAnswer();
					aerr = pa.getErr();
					peer = (HttpServicePeer) pa.getPeer();
				} else {
					aa = new Answer();
					aa.setStatus(StatusConst.ASYNC);
					aa.sureResult().put(PropertyConst.ASYNC_ID, aid);
				}
				doAnswer(peer, ctx, req, aa, aerr);
				return;
			}

			// sm.execute(peer, req, ctx)
			HttpServicePeer peer;
			String transId = co.getString(PropertyConst.TRANSACTION_ID);
			if (!Util.empty(transId)) {
				peer = getTrans(transId);
				if (peer == null) {
					throw new AppError("transaction '" + transId + " invalid");
				}
				peer.post(ctx, req, qo, co);
				return;
			}

			peer = new HttpServicePeer();
			peer.server = ServiceCallWebServer.this;
			peer.channelContext = ctx;
			peer.httpRequest = req;

			final ServiceMethod fsm = sm;
			final HttpServicePeer fpeer = peer;
			final Request fqo = qo;
			final Context fco = co;
			executorGroup.execute(new Runnable() {
				@Override
				public void run() {
					ChannelHandlerContext fctx = fpeer.channelContext;
					try {
						fsm.execute(fpeer, fqo, fco);
					} catch (Throwable t) {
						try {
							exceptionCaught(fctx, t);
						} catch (Exception e) {
							e.printStackTrace();
						}
					}
				}
			});

		}
	}

}
