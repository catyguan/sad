package bma.servicecall.sockserver;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufInputStream;
import io.netty.buffer.ByteBufOutputStream;
import io.netty.buffer.Unpooled;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.ByteToMessageDecoder;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.util.concurrent.DefaultEventExecutorGroup;

import java.io.IOException;
import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.LinkedBlockingDeque;

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
import bma.servicecall.sockcore.Message;
import bma.servicecall.sockcore.MessageReader;
import bma.servicecall.sockcore.MessageWriter;
import bma.servicecall.sockcore.SocketCoreConst;

public class ServiceCallSocketServer implements ServerBooter {

	private static final byte[] pingRData = { 9, 0, 0, 1, 1, 0, 0, 0, 0 };

	private EventLoopGroup bossGroup;
	private EventLoopGroup workerGroup;
	private DefaultEventExecutorGroup executorGroup;
	private Channel listener;
	private Timer timer;

	private int executors = 10;
	private int port;
	private int maxContentLength = 10 * 1024 * 1024;
	private boolean debugLog;
	private ServiceMux serviceMux;
	protected BaseServiceServ serv = new BaseServiceServ();

	public int getMaxContentLength() {
		return maxContentLength;
	}

	public void setMaxContentLength(int maxContentLength) {
		this.maxContentLength = maxContentLength;
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

	public int getPort() {
		return port;
	}

	public void setPort(int port) {
		this.port = port;
	}

	public void start() throws Exception {
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
				p.addLast(new MessageDecoder());
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

	public static void doAnswer(ChannelHandlerContext ctx, int mid, Answer aa,
			Throwable aerr) {
		Answer a = Answer.error2Answer(aa, aerr);
		if (a.getStatus() == 0) {
			a.setStatus(200);
		}
		ByteBuf buf = Unpooled.buffer();
		ByteBufOutputStream os = new ByteBufOutputStream(buf);
		MessageWriter mw = new MessageWriter(os);
		try {
			mw.sendAnswer(mid, a);
		} catch (IOException e) {
			throw AppError.handle(e);
		}
		ctx.writeAndFlush(buf);
	}

	public class MessageDecoder extends ByteToMessageDecoder {

		private int total;
		private Message message;
		private MessageReader mr;
		private byte[] hbuf = new byte[SocketCoreConst.HEADER_SIZE];

		@Override
		protected void decode(ChannelHandlerContext ctx, ByteBuf in,
				List<Object> out) throws Exception {
			while (true) {
				if (in.readableBytes() < SocketCoreConst.HEADER_SIZE) {
					return;
				}
				in.markReaderIndex();
				in.readBytes(hbuf);
				int sz = MessageReader.decodeHeaderSize(hbuf);
				if (in.readableBytes() < sz) {
					in.resetReaderIndex();
					return;
				}
				in.resetReaderIndex();
				if (total + sz + SocketCoreConst.HEADER_SIZE > maxContentLength) {
					int tsz = total + sz + SocketCoreConst.HEADER_SIZE;
					System.err.println("max content " + tsz + "/"
							+ maxContentLength);
					ctx.close();
					return;
				}
				ByteBufInputStream ins = new ByteBufInputStream(in);
				if (this.message == null) {
					this.message = new Message();
				}
				if (this.mr == null) {
					this.mr = new MessageReader(ins);
				}
				boolean done = this.mr.processFrame(ins, this.message);
				this.total = this.total + sz + SocketCoreConst.HEADER_SIZE;
				if (done) {
					this.total = 0;
					out.add(this.message);
					this.message = null;
				} else {
					if (Debuger.isEnable()) {
						Debuger.log(ctx.channel() + " read frame ==> " + total);
					}
				}
			}
		}
	}

	public class ServiceCallServerHandler extends
			SimpleChannelInboundHandler<Message> {

		private int mid;

		@Override
		public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
				throws Exception {
			if (!ctx.channel().isOpen()) {
				return;
			}
			if (debugLog) {
				System.err.println("exceptionCaught - ");
				cause.printStackTrace(System.err);
			} else {
				if (Debuger.isEnable()) {
					Debuger.log("exceptionCaught - " + cause);
					// cause.printStackTrace();
				}
			}

			if (mid != 0) {
				doAnswer(ctx, mid, null, cause);
				mid = 0;
			}
			// super.exceptionCaught(ctx, cause);
		}

		private SocketServicePeer workingPeer;
		private LinkedBlockingDeque<Message> wqueue = new LinkedBlockingDeque<Message>();

		@Override
		protected void channelRead0(ChannelHandlerContext ctx, Message msg)
				throws Exception {
			if (msg.getType() == SocketCoreConst.MT_PING) {
				if (!msg.isBoolFlag()) {
					ctx.writeAndFlush(Unpooled.copiedBuffer(pingRData));
					return;
				}
			}

			if (msg.getType() != SocketCoreConst.MT_REQUEST) {
				System.err.println("can't handle not Request message("
						+ msg.getType() + ")");
				ctx.close();
				return;
			}
			try {
				if (serviceMux == null) {
					throw new AppError("serviceMux invalid");
				}
				ServiceMethod sm = serviceMux.find(msg.getService(),
						msg.getMethod());
				if (sm == null) {
					Answer a = new Answer();
					a.setStatus(StatusConst.INVALID);
					a.setMessage("service(" + msg.getService() + ":"
							+ msg.getMethod() + ") not found");
					doAnswer(ctx, msg.getId(), a, null);
					return;
				}

				Request qo = msg.getRequest();
				if (qo == null) {
					qo = new Request();
					msg.setRequest(qo);
				}
				Context co = msg.getContext();
				if (co == null) {
					co = new Context();
					msg.setContext(co);
				}

				if (Debuger.isEnable()) {
					Debuger.log("QM : " + qo.toMap());
					Debuger.log("CM : " + co.toMap());
				}

				String aid = co.getString(PropertyConst.ASYNC_ID);
				if (!Util.empty(aid)) {
					PollAnswer pa = serv.pollAsync(aid);
					int mid = msg.getId();
					Answer aa;
					Exception aerr = null;
					if (pa != null) {
						aa = pa.getAnswer();
						aerr = pa.getErr();
					} else {
						aa = new Answer();
						aa.setStatus(StatusConst.ASYNC);
						aa.sureResult().put(PropertyConst.ASYNC_ID, aid);
					}
					doAnswer(ctx, mid, aa, aerr);
					return;
				}
				wqueue.add(msg);
				processMessages(ctx);
			} catch (Throwable t) {
				doAnswer(ctx, msg.getId(), null, t);
			}

		}

		@Override
		public void userEventTriggered(ChannelHandlerContext ctx, Object evt)
				throws Exception {
			super.userEventTriggered(ctx, evt);
			if (evt instanceof PeerDone) {
				this.workingPeer = null;
				processMessages(ctx);

			}
		}

		private void processMessages(ChannelHandlerContext ctx) {
			SocketServicePeer peer = this.workingPeer;
			if (peer != null) {
				return;
			}
			if (this.wqueue.isEmpty()) {
				return;
			}

			peer = new SocketServicePeer();
			peer.server = ServiceCallSocketServer.this;
			peer.wqueue = this.wqueue;
			peer.channelContext = ctx;
			this.workingPeer = peer;

			final ChannelHandlerContext fctx = ctx;
			final SocketServicePeer fpeer = peer;
			executorGroup.execute(new Runnable() {
				@Override
				public void run() {
					try {
						Message msg = wqueue.poll();
						if (msg == null) {
							return;
						}
						ServiceMethod sm = serviceMux.find(msg.getService(),
								msg.getMethod());
						try {
							fpeer.msg = msg;
							sm.execute(fpeer, msg.getRequest(),
									msg.getContext());
						} catch (Throwable t) {
							try {
								doAnswer(fctx, msg.getId(), null, t);
							} catch (Exception e) {
								e.printStackTrace();
							}
						}
					} finally {
						fpeer.channelContext = null;
						fpeer.msg = null;

						Object evt = new PeerDone();
						fctx.pipeline().fireUserEventTriggered(evt);
					}
				}
			});

		}
	}

}
