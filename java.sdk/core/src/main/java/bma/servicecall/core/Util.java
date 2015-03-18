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
}
