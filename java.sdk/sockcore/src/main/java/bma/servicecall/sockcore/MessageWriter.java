package bma.servicecall.sockcore;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;

import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.Value;
import bma.servicecall.core.ValueMap;
import bma.servicecall.core.ValueMapWalker;

public class MessageWriter extends OutputStream {
	private OutputStream w;
	private byte[] hbuf = new byte[SocketCoreConst.HEADER_SIZE];
	private ByteArrayOutputStream buf;

	public MessageWriter(OutputStream w) {
		super();
		this.w = w;
	}

	@Override
	public void close() throws IOException {
		super.close();
		if (this.buf != null) {
			this.buf.close();
			this.buf = null;
		}
	}

	@Override
	public void write(int b) throws IOException {
		this.w.write(b);
	}

	@Override
	public void write(byte[] b, int off, int len) throws IOException {
		this.w.write(b, off, len);
	}

	public void writeHeader(int mt, int sz) throws IOException {
		this.hbuf[0] = (byte) (mt & 0xFF);
		this.hbuf[1] = (byte) (sz >> 16 & 0xFF);
		this.hbuf[2] = (byte) (sz >> 8 & 0xFF);
		this.hbuf[3] = (byte) (sz & 0xFF);
		this.write(this.hbuf);
	}

	public void writeEnd() throws IOException {
		this.write(Coder.endData);
	}

	public void writeMessageId(int mid) throws IOException {
		this.writeHeader(SocketCoreConst.MT_MESSAGE_ID, 4);
		Coder.encodeFixInt32(this, mid);
	}

	public void writeFlag() throws IOException {
		this.write(Coder.flagData);
	}

	protected ByteArrayOutputStream sbuf() {
		if (this.buf == null) {
			this.buf = new ByteArrayOutputStream();
		}
		this.buf.reset();
		return this.buf;
	}

	public void writeAddress(String s, String m) throws IOException {
		ByteArrayOutputStream buf = this.sbuf();
		int n = 0;
		n += Coder.encodeLenString(buf, s);
		n += Coder.encodeLenString(buf, m);
		this.writeHeader(SocketCoreConst.MT_ADDRESS, n);
		this.write(buf.toByteArray());
	}

	public void writeData(String name, Value val) throws IOException {
		ByteArrayOutputStream buf = this.sbuf();
		int n = 0;
		n += Coder.encodeLenString(buf, name);
		n += Coder.encodeVar(buf, val);
		this.writeHeader(SocketCoreConst.MT_DATA, n);
		this.write(buf.toByteArray());
	}

	public void writeContext(String name, Value val) throws IOException {
		ByteArrayOutputStream buf = this.sbuf();
		int n = 0;
		n += Coder.encodeLenString(buf, name);
		n += Coder.encodeVar(buf, val);
		this.writeHeader(SocketCoreConst.MT_CONTEXT, n);
		this.write(buf.toByteArray());
	}

	public void writeAnswer(int st, String msg) throws IOException {
		ByteArrayOutputStream buf = this.sbuf();
		int n = 0;
		n += Coder.encodeInt32(buf, st);
		n += Coder.encodeLenString(buf, msg);
		this.writeHeader(SocketCoreConst.MT_ANSWER, n);
		this.write(buf.toByteArray());
	}

	public void sendRequest(int mid, String s, String m, Request req,
			Context ctx) throws IOException {
		this.writeMessageId(mid);
		this.writeFlag();
		this.writeAddress(s, m);
		if (req != null) {
			req.walk(new ValueMapWalker() {

				@Override
				public boolean walk(String k, Value v) {
					try {
						writeData(k, v);
					} catch (IOException e) {
						throw AppError.handle(e);
					}
					return false;
				}
			});
		}

		if (ctx != null) {
			ctx.walk(new ValueMapWalker() {

				@Override
				public boolean walk(String k, Value v) {
					try {
						writeContext(k, v);
					} catch (IOException e) {
						throw AppError.handle(e);
					}
					return false;
				}
			});
		}
		this.writeEnd();
	}

	public void sendAnswer(int mid, Answer an) throws IOException {
		this.writeMessageId(mid);
		this.writeAnswer(an.getStatus(), an.getMessage());

		ValueMap rs = an.getResult();
		if (rs != null) {
			rs.walk(new ValueMapWalker() {

				@Override
				public boolean walk(String k, Value v) {
					try {
						writeData(k, v);
					} catch (IOException e) {
						throw AppError.handle(e);
					}
					return false;
				}
			});
		}

		ValueMap ctx = an.getContext();
		if (ctx != null) {
			ctx.walk(new ValueMapWalker() {

				@Override
				public boolean walk(String k, Value v) {
					try {
						writeContext(k, v);
					} catch (IOException e) {
						throw AppError.handle(e);
					}
					return false;
				}
			});
		}
		this.writeEnd();
	}
}
