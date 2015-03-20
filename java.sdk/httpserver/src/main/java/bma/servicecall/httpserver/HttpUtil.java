package bma.servicecall.httpserver;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpHeaders;
import io.netty.handler.codec.http.HttpHeaders.Names;
import io.netty.handler.codec.http.HttpHeaders.Values;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.util.CharsetUtil;

public class HttpUtil {
	public static void sendError(ChannelHandlerContext ctx, HttpRequest req,
			HttpResponseStatus status, String msg, Throwable err) {
		if (status == null) {
			status = HttpResponseStatus.INTERNAL_SERVER_ERROR;
		}
		if (msg == null && err != null) {
			msg = err.getMessage();
		}
		sendResponse(ctx, req, status, msg);
	}

	public static void sendResponse(ChannelHandlerContext ctx, HttpRequest req,
			HttpResponseStatus status, String msg) {
		boolean keepAlive = false;
		if (req != null) {
			keepAlive = HttpHeaders.isKeepAlive(req);
		}
		ByteBuf buf = Unpooled.copiedBuffer(msg, CharsetUtil.UTF_8);
		FullHttpResponse response = new DefaultFullHttpResponse(
				HttpVersion.HTTP_1_1, status, buf);
		response.headers().set(Names.CONTENT_TYPE, "text/plain; charset=UTF-8");
		if (keepAlive) {
			response.headers().set(Names.CONTENT_LENGTH, buf.readableBytes());
			response.headers().set(Names.CONNECTION, Values.KEEP_ALIVE);
		}
		ctx.writeAndFlush(response);
		if (!keepAlive) {
			ctx.writeAndFlush(Unpooled.EMPTY_BUFFER).addListener(
					ChannelFutureListener.CLOSE);
		}
	}
}
