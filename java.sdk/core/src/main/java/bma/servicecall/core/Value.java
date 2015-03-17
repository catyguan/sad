package bma.servicecall.core;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

public class Value {
	private int type;
	private Object value;

	public Value(int t, Object v) {
		super();
		this.type = t;
		this.value = v;
	}

	@SuppressWarnings("rawtypes")
	public static Value create(Object v) {
		if (v == null) {
			return new Value(TypeConst.NULL, null);
		}
		if (v instanceof Value) {
			return (Value) v;
		}
		if (v instanceof Boolean) {
			return new Value(TypeConst.BOOL, v);
		}
		if (v instanceof Byte) {
			return new Value(TypeConst.INT, ((Byte) v).intValue());
		}
		if (v instanceof Integer) {
			return new Value(TypeConst.INT, v);
		}
		if (v instanceof Short) {
			return new Value(TypeConst.INT, ((Short) v).intValue());
		}
		if (v instanceof Long) {
			return new Value(TypeConst.LONG, v);
		}
		if (v instanceof Float) {
			return new Value(TypeConst.FLOAT, v);
		}
		if (v instanceof Double) {
			return new Value(TypeConst.DOUBLE, v);
		}
		if (v instanceof String) {
			return new Value(TypeConst.STRING, v);
		}
		if (v instanceof ValueArray) {
			return new Value(TypeConst.ARRAY, v);
		}
		if (v instanceof List) {
			ValueArray a = ValueArray.create((List) v);
			return new Value(TypeConst.ARRAY, a);
		}
		if (v instanceof ValueMap) {
			return new Value(TypeConst.MAP, v);
		}
		if (v instanceof Map) {
			ValueMap a = ValueMap.create((Map) v);
			return new Value(TypeConst.MAP, a);
		}
		if (v instanceof byte[]) {
			return new Value(TypeConst.BINARY, v);
		}
		throw new AppError("unknow value(" + v.getClass().getName() + ")");
	}

	public int getType() {
		return type;
	}

	public Object getValue() {
		return value;
	}

	public Object toValue() {
		return this.convertValue(DataConverter.noConverter);
	}

	@SuppressWarnings({ "rawtypes", "unchecked" })
	public Object convertValue(DataConverter dc) {
		switch (this.type) {
		case TypeConst.NULL:
			return dc.convert(this.type, null);
		case TypeConst.ARRAY: {
			ValueArray o = (ValueArray) this.value;
			List r = new ArrayList();
			Iterator<Value> it = o.iterator();
			if (it != null) {
				while (it.hasNext()) {
					Value vv = it.next();
					Object v = vv.convertValue(dc);
					r.add(v);
				}
			}
			return r;
		}
		case TypeConst.MAP: {
			ValueMap o = (ValueMap) this.value;
			Map r = new TreeMap();
			Iterator<String> it = o.iterator();
			if (it != null) {
				while (it.hasNext()) {
					String key = it.next();
					Value vv = o.get(key);
					Object v = vv.convertValue(dc);
					r.put(key, v);
				}
			}
			return null;
		}
		default:
			return dc.convert(this.type, this.value);
		}
	}

	public Object as(int type) {
		switch (type) {
		case TypeConst.BOOL:
			return this.asBool();
		case TypeConst.INT:
			return this.asInt();
		case TypeConst.LONG:
			return this.asLong();
		case TypeConst.FLOAT:
			return this.asFloat();
		case TypeConst.DOUBLE:
			return this.asDouble();
		case TypeConst.STRING:
			return this.asString();
		case TypeConst.VAR:
			return this;
		case TypeConst.ARRAY:
			return this.asArray();
		case TypeConst.MAP:
			return this.asMap();
		case TypeConst.BINARY:
			return this.asBinary();
		}
		return null;
	}

	public boolean asBool() {
		if (this.value == null) {
			return false;
		}
		switch (this.type) {
		case TypeConst.BOOL:
			return (Boolean) this.value;
		case TypeConst.INT:
			return ((Number) this.value).intValue() != 0;
		case TypeConst.LONG:
			return ((Number) this.value).longValue() != 0;
		case TypeConst.FLOAT:
			return ((Number) this.value).floatValue() != 0;
		case TypeConst.DOUBLE:
			return ((Number) this.value).doubleValue() != 0;
		case TypeConst.STRING:
			String s = (String) this.value;
			if (s == "") {
				return false;
			}
			if (s.equalsIgnoreCase("true") || s.equalsIgnoreCase("yes")) {
				return true;
			} else if (s.equalsIgnoreCase("false") || s.equalsIgnoreCase("no")) {
				return false;
			} else {
				return false;
			}
		case TypeConst.ARRAY:
			ValueArray as = (ValueArray) this.value;
			return as.len() != 0;
		case TypeConst.MAP:
			ValueMap ms = (ValueMap) this.value;
			return ms.len() != 0;
		case TypeConst.BINARY:
			byte[] bs = (byte[]) this.value;
			return bs.length != 0;
		default:
			return false;
		}
	}

	public int asInt() {
		if (this.value == null) {
			return 0;
		}
		switch (this.type) {
		case TypeConst.BOOL:
			return ((Boolean) this.value).booleanValue() ? 1 : 0;
		case TypeConst.INT:
			return ((Number) this.value).intValue();
		case TypeConst.LONG:
			return ((Number) this.value).intValue();
		case TypeConst.FLOAT:
			return ((Number) this.value).intValue();
		case TypeConst.DOUBLE:
			return ((Number) this.value).intValue();
		case TypeConst.STRING:
			String s = (String) this.value;
			try {
				return Integer.parseInt(s);
			} catch (Exception e) {
				return 0;
			}
		case TypeConst.ARRAY:
			return 0;
		case TypeConst.MAP:
			return 0;
		case TypeConst.BINARY:
			return 0;
		default:
			return 0;
		}
	}

	public long asLong() {
		if (this.value == null) {
			return 0;
		}
		switch (this.type) {
		case TypeConst.BOOL:
			return ((Boolean) this.value).booleanValue() ? 1 : 0;
		case TypeConst.INT:
			return ((Number) this.value).longValue();
		case TypeConst.LONG:
			return ((Number) this.value).longValue();
		case TypeConst.FLOAT:
			return ((Number) this.value).longValue();
		case TypeConst.DOUBLE:
			return ((Number) this.value).longValue();
		case TypeConst.STRING:
			String s = (String) this.value;
			try {
				return Long.parseLong(s);
			} catch (Exception e) {
				return 0;
			}
		case TypeConst.ARRAY:
			return 0;
		case TypeConst.MAP:
			return 0;
		case TypeConst.BINARY:
			return 0;
		default:
			return 0;
		}
	}

	public float asFloat() {
		if (this.value == null) {
			return 0;
		}
		switch (this.type) {
		case TypeConst.BOOL:
			return ((Boolean) this.value).booleanValue() ? 1 : 0;
		case TypeConst.INT:
			return ((Number) this.value).floatValue();
		case TypeConst.LONG:
			return ((Number) this.value).floatValue();
		case TypeConst.FLOAT:
			return ((Number) this.value).floatValue();
		case TypeConst.DOUBLE:
			return ((Number) this.value).floatValue();
		case TypeConst.STRING:
			String s = (String) this.value;
			try {
				return Float.parseFloat(s);
			} catch (Exception e) {
				return 0;
			}
		case TypeConst.ARRAY:
			return 0;
		case TypeConst.MAP:
			return 0;
		case TypeConst.BINARY:
			return 0;
		default:
			return 0;
		}
	}

	public double asDouble() {
		if (this.value == null) {
			return 0;
		}
		switch (this.type) {
		case TypeConst.BOOL:
			return ((Boolean) this.value).booleanValue() ? 1 : 0;
		case TypeConst.INT:
			return ((Number) this.value).doubleValue();
		case TypeConst.LONG:
			return ((Number) this.value).doubleValue();
		case TypeConst.FLOAT:
			return ((Number) this.value).doubleValue();
		case TypeConst.DOUBLE:
			return ((Number) this.value).doubleValue();
		case TypeConst.STRING:
			String s = (String) this.value;
			try {
				return Double.parseDouble(s);
			} catch (Exception e) {
				return 0;
			}
		case TypeConst.ARRAY:
			return 0;
		case TypeConst.MAP:
			return 0;
		case TypeConst.BINARY:
			return 0;
		default:
			return 0;
		}
	}

	public String asString() {
		if (this.value == null)
			return "";
		return this.value.toString();
	}

	public ValueArray asArray() {
		if (this.type == TypeConst.ARRAY) {
			return (ValueArray) this.value;
		}
		return null;
	}

	public ValueMap asMap() {
		if (this.type == TypeConst.MAP) {
			return (ValueMap) this.value;
		}
		return null;
	}

	public byte[] asBinary() {
		if (this.type == TypeConst.BINARY) {
			return (byte[]) this.value;
		}
		return null;
	}

	@Override
	public String toString() {
		return "V(" + this.type + "," + this.value + ")";
	}
}
