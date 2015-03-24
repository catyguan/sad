package bma.servicecall.sockcore;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import bma.servicecall.core.AppError;
import bma.servicecall.core.TypeConst;
import bma.servicecall.core.Value;
import bma.servicecall.core.ValueArray;
import bma.servicecall.core.ValueMap;

public class Coder {

	public static final Charset UTF8 = Charset.forName("UTF-8");

	public static final byte[] endData = { 0, 0, 0, 0 };
	public static final byte[] flagData = { 2, 0, 0, 0 };

	public static class AddressInfo {
		public String service;
		public String method;
	}

	public static class DataInfo {
		public String key;
		public Value value;
	}

	public static class AnswerInfo {
		public int status;
		public String message;
	}

	public static void readAll(InputStream in, byte[] bs) throws IOException {
		readAll(in, bs, 0, bs.length);
	}

	public static void readAll(InputStream in, byte[] bs, int off, int len)
			throws IOException {
		int pos = off;
		int epos = off + len;
		while (pos < epos) {
			int n = in.read(bs, pos, epos - pos);
			pos += n;
		}
	}

	public static int encodeBool(OutputStream w, boolean v) throws IOException {
		int b = v ? 1 : 0;
		w.write(b);
		return 1;
	}

	public static boolean decodeBool(InputStream r) throws IOException {
		int v = r.read();
		return v != 0;
	}

	public static int encodeInt32(OutputStream w, int v) throws IOException {
		long l = ((long) v) << 1;
		if (v < 0) {
			l = ~l;
		}
		return encodeUInt64(w, l);
	}

	public static int decodeInt32(InputStream r) throws IOException {
		long l = decodeUInt64(r);
		long l2 = (long) (l >> 1);
		if ((l & 1) != 0) {
			l2 = ~l2;
		}
		return Long.valueOf(l2).intValue();
	}

	public static int encodeInt64(OutputStream w, long v) throws IOException {
		long l1 = ((long) v) << 1;
		if (v < 0) {
			l1 = ~l1;
		}
		return encodeUInt64(w, l1);
	}

	public static long decodeInt64(InputStream r) throws IOException {
		long l = decodeUInt64(r);
		long l2 = (long) (l >> 1);
		if ((l & 1) != 0) {
			l2 = ~l2;
		}
		return l2;
	}

	public static int encodeUInt64(OutputStream w, long v) throws IOException {
		int i = 0;
		while (v >= 0x80) {
			w.write((byte) v | 0x80);
			v >>= 7;
			i++;
		}
		w.write((byte) v);
		return i + 1;
	}

	public static long decodeUInt64(InputStream r) throws IOException {
		long s = 0;
		int b;
		int w = 0;
		int i = 0;
		while (true) {
			b = r.read();
			if (b < 0x80) {
				if (i > 9 || i == 9 && b > 1) {
					return 0; // overflow
				}
				return (s | (long) (b) << w);

			}
			s |= (long) (b & 0x7f) << w;
			w += 7;
			i++;
		}
	}

	public static int encodeFixInt8(OutputStream w, byte v) throws IOException {
		w.write(v);
		return 1;
	}

	public static byte decodeFixInt8(InputStream r) throws IOException {
		return (byte) (r.read() & 0xFF);
	}

	public static int encodeFixInt16(OutputStream w, short v)
			throws IOException {
		w.write((byte) (v >> 8));
		w.write((byte) (v >> 0));
		return 2;
	}

	public static short decodeFixInt16(InputStream r) throws IOException {
		int s = 0;
		s += (r.read() & 0xff) << 8;
		s += (r.read() & 0xff);
		return (short) s;
	}

	public static int encodeFixInt32(OutputStream w, int v) throws IOException {
		w.write((byte) (v >> 24));
		w.write((byte) (v >> 16));
		w.write((byte) (v >> 8));
		w.write((byte) (v >> 0));
		return 4;
	}

	public static int decodeFixInt32(InputStream r) throws IOException {
		int s = 0;
		s += (r.read() & 0xff) << 24;
		s += (r.read() & 0xff) << 16;
		s += (r.read() & 0xff) << 8;
		s += (r.read() & 0xff);
		return s;
	}

	public static int encodeFixInt64(OutputStream w, long v) throws IOException {
		w.write((byte) (v >> 56));
		w.write((byte) (v >> 48));
		w.write((byte) (v >> 40));
		w.write((byte) (v >> 32));
		w.write((byte) (v >> 24));
		w.write((byte) (v >> 16));
		w.write((byte) (v >> 8));
		w.write((byte) (v >> 0));
		return 8;
	}

	public static long decodeFixInt64(InputStream r) throws IOException {
		long s = 0;
		s += (long) (r.read() & 0xff) << 56;
		s += (long) (r.read() & 0xff) << 48;
		s += (long) (r.read() & 0xff) << 40;
		s += (long) (r.read() & 0xff) << 32;
		s += (long) (r.read() & 0xff) << 24;
		s += (long) (r.read() & 0xff) << 16;
		s += (long) (r.read() & 0xff) << 8;
		s += (long) (r.read() & 0xff);
		return s;
	}

	public static int encodeFloat32(OutputStream w, float v) throws IOException {
		int l = Float.floatToIntBits(v);
		return encodeFixInt32(w, l);
	}

	public static float decodeFloat32(InputStream r) throws IOException {
		int l = decodeFixInt32(r);
		return Float.intBitsToFloat(l);
	}

	public static int encodeFloat64(OutputStream w, double v)
			throws IOException {
		long l = Double.doubleToLongBits(v);
		return encodeFixInt64(w, l);
	}

	public static double decodeFloat64(InputStream r) throws IOException {
		long l = decodeFixInt64(r);
		return Double.longBitsToDouble(l);
	}

	public static int encodeLenBytes(OutputStream w, byte[] bs)
			throws IOException {
		int l = 0;
		if (bs != null) {
			l = bs.length;
		}
		int n = encodeInt32(w, l);
		if (l > 0) {
			w.write(bs);
			n += l;
		}
		return n;
	}

	public static byte[] decodeLenBytes(InputStream r, int maxlen)
			throws IOException {
		int l = decodeInt32(r);
		if (maxlen <= 0) {
			maxlen = SocketCoreConst.DATA_MAXLEN;
		}
		if (l > maxlen) {
			throw new AppError("too large bytes block - " + l + "/" + maxlen);
		}
		if (l == 0) {
			return new byte[0];
		}
		byte[] p = new byte[l];
		if (l > 0) {
			readAll(r, p);
		}
		return p;
	}

	public static int encodeLenString(OutputStream w, String s)
			throws IOException {
		if (s == null) {
			return encodeInt32(w, 0);
		}
		byte[] bs = s.getBytes(UTF8);
		return encodeLenBytes(w, bs);
	}

	public static String decodeLenString(InputStream r, int maxlen)
			throws IOException {
		byte[] bs = decodeLenBytes(r, maxlen);
		if (bs == null) {
			return null;
		}
		return new String(bs, UTF8);
	}

	@SuppressWarnings("rawtypes")
	public static int encodeMap(OutputStream w, Map obj) throws IOException {
		if (obj == null) {
			return encodeInt32(w, 0);
		}
		// map长度
		int n = 0;
		int len = obj.size();
		n += encodeInt32(w, len);
		for (Object o : obj.keySet()) {
			n += encodeLenString(w, o.toString());
			n += encodeVar(w, obj.get(o));
		}
		return n;
	}

	public static Map<String, Value> decodeMap(InputStream r)
			throws IOException {
		Map<String, Value> obj = new HashMap<String, Value>();
		int size = decodeInt32(r);
		int mark = 1;
		while (mark <= size) {
			String key = decodeLenString(r, 0);
			Value val = decodeVar(r);
			obj.put(key, val);
			mark++;
		}
		return obj;
	}

	@SuppressWarnings("rawtypes")
	public static int encodeList(OutputStream w, List obj) throws IOException {
		if (obj == null) {
			return encodeInt32(w, 0);
		}
		int len = obj.size();
		int n = 0;
		n += encodeInt32(w, len);
		for (Object o : obj) {
			n += encodeVar(w, o);
		}
		return n;
	}

	/**
	 * List解码
	 * 
	 * @param buf
	 * @return
	 * @throws IOException
	 */
	public static List<Value> decodeList(InputStream r) throws IOException {
		List<Value> obj = new ArrayList<Value>();
		int size = decodeInt32(r);
		int mark = 1;
		while (mark <= size) {
			obj.add(decodeVar(r));
			mark++;
		}
		return obj;
	}

	@SuppressWarnings("rawtypes")
	public static int doEncodeVar(OutputStream buf, int type, Object obj)
			throws IOException {
		switch (type) {
		case TypeConst.NULL:
			buf.write(type);
			return 1;
		case TypeConst.INT:// int32
			if (obj instanceof Number) {
				Number i = (Number) obj;
				buf.write(type);
				return encodeInt32(buf, i.intValue()) + 1;
			}
			throw new IllegalArgumentException("not Int type");
		case TypeConst.LONG:// int64
			if (obj instanceof Number) {
				Number l = (Number) obj;
				buf.write(type);
				return encodeInt64(buf, l.longValue()) + 1;
			}
			throw new IllegalArgumentException("not long type");
		case TypeConst.FLOAT:// float32
			if (obj instanceof Number) {
				Number f = (Number) obj;
				buf.write(type);
				return encodeFloat32(buf, f.floatValue()) + 1;
			}
			throw new IllegalArgumentException("not float type");
		case TypeConst.DOUBLE:// float64
			if (obj instanceof Number) {
				Number d = (Number) obj;
				buf.write(type);
				return encodeFloat64(buf, d.doubleValue()) + 1;
			}
			throw new IllegalArgumentException("not double type");
		case TypeConst.STRING:// lenString
			if (true) {
				String s = obj.toString();
				buf.write(type);
				return encodeLenString(buf, s) + 1;
			}
		case TypeConst.MAP:// map
			if (obj instanceof ValueMap) {
				obj = ((ValueMap) obj).theData();
			}
			if (obj instanceof Map) {
				Map m = (Map) obj;
				buf.write(type);
				return encodeMap(buf, m) + 1;
			}
			throw new IllegalArgumentException("not map type");
		case TypeConst.ARRAY:// list
			if (obj instanceof ValueArray) {
				obj = ((ValueArray) obj).theData();
			}
			if (obj instanceof List) {
				List l = (List) obj;
				buf.write(type);
				return encodeList(buf, l) + 1;
			}
			throw new IllegalArgumentException("not list type");
		case TypeConst.BINARY:// list
			if (obj instanceof byte[]) {
				byte[] s = (byte[]) obj;
				buf.write(type);
				return encodeLenBytes(buf, s) + 1;
			}
			throw new IllegalArgumentException("not list type");
		default:
			throw new IllegalArgumentException("invalid type(" + type + ")");
		}
	}

	public static int encodeVar(OutputStream w, Object v) throws IOException {
		if (v == null) {
			return doEncodeVar(w, TypeConst.NULL, null);
		}
		Value val;
		if (v instanceof Value) {
			val = (Value) v;
		} else {
			val = Value.create(v);
		}
		return doEncodeVar(w, val.getType(), val.getValue());
	}

	public static Object doDecodeVar(InputStream buf, int type)
			throws IOException {
		switch (type) {
		case TypeConst.NULL:
			return null;
		case TypeConst.BOOL:
			return decodeBool(buf);
		case TypeConst.INT:// int32
			return decodeInt32(buf);
		case TypeConst.LONG:// int64
			return decodeInt64(buf);
		case TypeConst.FLOAT:// float32
			return decodeFloat32(buf);
		case TypeConst.DOUBLE:// float64
			return decodeFloat64(buf);
		case TypeConst.STRING:// lenString
			return decodeLenString(buf, 0);
		case TypeConst.BINARY:// lenString
			return decodeLenBytes(buf, 0);
		case TypeConst.MAP:// lenString
			return new ValueMap(decodeMap(buf));
		case TypeConst.ARRAY:// lenString
			return new ValueArray(decodeList(buf));
		default:
			break;
		}
		return null;
	}

	public static Value decodeVar(InputStream r) throws IOException {
		int type = r.read();
		Object o = doDecodeVar(r, type);
		return new Value(type, o);
	}
}
