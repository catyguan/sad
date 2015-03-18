package bma.servicecall.core;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.GregorianCalendar;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;

public class Util {
	public static long currentUnixTimestamp() {
		return new GregorianCalendar().getTimeInMillis() / 1000;
	}

	@SuppressWarnings({ "rawtypes", "unchecked" })
	public static Object scopyv(Object v) {
		if (v == null) {
			return v;
		}
		if (v instanceof Map) {
			Map m = (Map) v;
			Map<String, Object> r = new HashMap<String, Object>();
			for (Object o : m.entrySet()) {
				Entry e = (Entry) o;
				Object vv = scopyv(e.getValue());
				r.put(e.getKey().toString(), vv);
			}
			return r;
		}
		if (v instanceof List) {
			List l = (List) v;
			List r = new ArrayList();
			for (Object o : l) {
				Object vv = scopyv(o);
				r.add(vv);
			}
			return r;
		}
		return v;
	}

	private final static char[] digit = { '0', '1', '2', '3', '4', '5', '6',
			'7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f' };

	public static String byte2Hex(byte[] data) {
		StringBuffer buf = new StringBuffer(data.length * 2);
		for (byte ib : data) {
			buf.append(digit[(ib >>> 4) & 0x0F]);
			buf.append(digit[ib & 0x0F]);
		}
		return buf.toString();
	}

	private static int hexValue(char c) {
		switch (c) {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
			return c - '0';
		case 'a':
		case 'A':
			return 10;
		case 'b':
		case 'B':
			return 11;
		case 'c':
		case 'C':
			return 12;
		case 'd':
		case 'D':
			return 13;
		case 'e':
		case 'E':
			return 14;
		case 'f':
		case 'F':
			return 15;
		default:
			throw new IllegalArgumentException("invalid HEX code '" + (int) c
					+ "'");
		}
	}

	public static byte[] hex2byte(String s) {
		if (s == null)
			return null;
		if (s.isEmpty())
			return null;
		int len = s.length();

		byte[] rb = new byte[len / 2];
		char[] rc = s.toCharArray();
		for (int i = 0, j = 0; i < rc.length; i += 2, j++) {
			rb[j] = (byte) (((hexValue(rc[i]) << 4) + hexValue(rc[i + 1])) & 0xFF);
		}
		return rb;
	}

	public static String md5(byte[] data) {
		MessageDigest md;
		try {
			md = MessageDigest.getInstance("MD5");
		} catch (NoSuchAlgorithmException e) {
			throw new UnsupportedOperationException("md5", e);
		}
		md.update(data);
		byte md5[] = md.digest();
		return byte2Hex(md5);
	}

	public static boolean empty(String e) {
		return e == null || e.length() == 0;
	}
}
