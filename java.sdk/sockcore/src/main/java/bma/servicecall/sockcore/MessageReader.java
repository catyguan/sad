package bma.servicecall.sockcore;

import java.io.EOFException;
import java.io.IOException;
import java.io.InputStream;

import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.Request;
import bma.servicecall.core.Value;
import bma.servicecall.sockcore.Coder.AddressInfo;
import bma.servicecall.sockcore.Coder.AnswerInfo;
import bma.servicecall.sockcore.Coder.DataInfo;

public class MessageReader extends InputStream {
	private InputStream r;
	private boolean h;
	private int l;
	private int sz;
	private byte mt;
	private byte[] hbs = new byte[SocketCoreConst.HEADER_SIZE];

	public MessageReader(InputStream r) {
		super();
		this.r = r;
	}

	public byte next() throws IOException {
		this.h = true;
		this.readHeader();
		this.h = false;
		this.l = this.sz;
		return this.mt;
	}

	public int len() {
		return this.sz;
	}

	@Override
	public int read(byte[] b, int off, int len) throws IOException {
		if (!this.h) {
			if (this.l <= 0) {
				throw new EOFException();
			}
			if (len > this.l) {
				len = this.l;
			}
		}
		Coder.readAll(this.r, b, off, len);
		int n = len;
		if (this.h) {
			this.l -= n;
		}
		// String sss = "";
		// for (int i = off; i < off + len; i++) {
		// sss += b[i] + "; ";
		// }
		// System.out.println("READ -> " + sss);
		return n;
	}

	@Override
	public int read() throws IOException {
		if (!this.h) {
			if (this.l <= 0) {
				throw new EOFException();
			}
		}
		int b = this.r.read();
		if (this.h) {
			this.l -= 1;
		}
		// System.out.println("READ -> " + b);
		return b;
	}

	public void readHeader() throws IOException {
		Coder.readAll(this, this.hbs);
		this.decodeHeader(this.hbs);
	}

	public static int decodeHeaderSize(byte[] b) {
		return ((int) b[3]) | ((int) b[2]) << 8 | ((int) b[1]) << 16;
	}

	public void decodeHeader(byte[] b) {
		this.mt = (byte) b[0];
		this.sz = decodeHeaderSize(b);
	}

	// //////
	public int readMessageId() throws IOException {
		if (this.mt != SocketCoreConst.MT_MESSAGE_ID) {
			throw new AppError("MT(" + this.mt + ") invalid MT_MESSAGE_ID");
		}
		int v = Coder.decodeFixInt32(this);
		return v;
	}

	public AddressInfo readAddress() throws IOException {
		if (this.mt != SocketCoreConst.MT_ADDRESS) {
			throw new AppError("MT(" + this.mt + ") invalid MT_ADDRESS");
		}
		String s = Coder.decodeLenString(this, 0);
		String m = Coder.decodeLenString(this, 0);
		AddressInfo addr = new AddressInfo();
		addr.service = s;
		addr.method = m;
		return addr;
	}

	public DataInfo readData() throws IOException {
		if (this.mt != SocketCoreConst.MT_DATA) {
			new AppError("MT(" + this.mt + ") invalid MT_DATA");
		}
		String key = Coder.decodeLenString(this, 0);
		Value val = Coder.decodeVar(this);
		DataInfo info = new DataInfo();
		info.key = key;
		info.value = val;
		return info;
	}

	public DataInfo readContext() throws IOException {
		if (this.mt != SocketCoreConst.MT_CONTEXT) {
			throw new AppError("MT(" + this.mt + ") invalid MT_CONTEXT");
		}
		String key = Coder.decodeLenString(this, 0);
		Value val = Coder.decodeVar(this);
		DataInfo info = new DataInfo();
		info.key = key;
		info.value = val;
		return info;
	}

	public AnswerInfo readAnswer() throws IOException {
		if (this.mt != SocketCoreConst.MT_ANSWER) {
			throw new AppError("MT(" + this.mt + ") invalid MT_ANSWER");
		}
		int status = Coder.decodeInt32(this);
		String message = Coder.decodeLenString(this, 0);
		AnswerInfo info = new AnswerInfo();
		info.status = status;
		info.message = message;
		return info;
	}

	public byte nextMessage(Message msg) throws IOException {
		msg.reset();
		while (true) {
			if (processFrame(null, msg)) {
				return msg.type;
			}
		}
	}

	public boolean processFrame(InputStream in, Message msg) throws IOException {
		if (in != null) {
			this.r = in;
		}
		byte mt = this.next();
		// System.out.println("read line - " + mt);
		switch (mt) {
		case SocketCoreConst.MT_END:
			if (msg.type == SocketCoreConst.MT_ANSWER) {
				Answer an = msg.answer;
				if (msg.request != null) {
					an.setResult(msg.request);
					msg.request = null;
				}
				if (msg.context != null) {
					an.setContext(msg.context);
					msg.context = null;
				}
			}
			if (Debuger.isEnable()) {
				Debuger.log("read message -> " + msg.type + ", " + msg);
			}
			return true;
		case SocketCoreConst.MT_MESSAGE_ID:
			int v = this.readMessageId();
			msg.id = v;
			break;
		case SocketCoreConst.MT_PING:
			boolean bv = Coder.decodeBool(this);
			msg.boolFlag = bv;
			msg.type = mt;
			break;
		case SocketCoreConst.MT_REQUEST:
			msg.type = mt;
			break;
		case SocketCoreConst.MT_ADDRESS:
			AddressInfo info1 = this.readAddress();
			msg.service = info1.service;
			msg.method = info1.method;
			break;
		case SocketCoreConst.MT_DATA:
			DataInfo info2 = this.readData();
			if (msg.request == null) {
				msg.request = new Request();
			}
			msg.request.set(info2.key, info2.value);
			break;
		case SocketCoreConst.MT_CONTEXT:
			DataInfo info3 = this.readContext();
			if (msg.context == null) {
				msg.context = new Context();
			}
			msg.context.set(info3.key, info3.value);
			break;
		case SocketCoreConst.MT_ANSWER:
			AnswerInfo info4 = this.readAnswer();
			msg.type = mt;
			if (msg.answer == null) {
				msg.answer = new Answer();
			}
			msg.answer.setStatus(info4.status);
			msg.answer.setMessage(info4.message);
			break;
		default:
			throw new AppError("unknow MessageType(" + mt + ")");
		}
		return false;
	}
}
