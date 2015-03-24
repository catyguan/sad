package bma.servicecall.core;

import java.util.Iterator;
import java.util.Map;
import java.util.Map.Entry;
import java.util.TreeMap;

public class ValueMap {

	private Map<String, Value> data;

	public ValueMap(Map<String, Value> d) {
		super();
		this.data = d;
	}

	@SuppressWarnings("rawtypes")
	public static ValueMap create(Map data) {
		ValueMap o = new ValueMap(null);
		o.initValueMap(data);
		return o;
	}

	@SuppressWarnings("rawtypes")
	public void initValueMap(Map data) {
		if (data != null) {
			Map<String, Value> r = null;
			r = new TreeMap<String, Value>();
			for (Object o : data.entrySet()) {
				Entry e = (Entry) o;
				r.put(e.getKey().toString(), Value.create(e.getValue()));
			}
			this.data = r;
		}
	}

	public Iterator<String> iterator() {
		if (this.data == null)
			return null;
		return this.data.keySet().iterator();
	}

	public Map<String, Object> toMap() {
		return this.convertMap(DataConverter.noConverter);
	}

	public Map<String, Object> convertMap(DataConverter dc) {
		Map<String, Object> r = new TreeMap<String, Object>();
		if (this.data != null) {
			for (Entry<String, Value> e : this.data.entrySet()) {
				Value v = e.getValue();
				if (v != null) {
					r.put(e.getKey(), v.convertValue(dc));
				}
			}
		}
		return r;
	}

	public boolean has(String key) {
		return this.data == null ? false : this.data.containsKey(key);
	}

	public Value get(String key) {
		if (this.data == null)
			return null;
		return this.data.get(key);
	}

	public int len() {
		if (this.data == null)
			return 0;
		return this.data.size();
	}

	public boolean getBool(String key) {
		Value v = this.get(key);
		if (v == null) {
			return false;
		}
		return v.asBool();
	}

	public int getInt(String key) {
		Value v = this.get(key);
		if (v == null) {
			return 0;
		}
		return v.asInt();
	}

	public long getLong(String key) {
		Value v = this.get(key);
		if (v == null) {
			return 0;
		}
		return v.asLong();
	}

	public float getFloat(String key) {
		Value v = this.get(key);
		if (v == null) {
			return 0;
		}
		return v.asFloat();
	}

	public double getDouble(String key) {
		Value v = this.get(key);
		if (v == null) {
			return 0;
		}
		return v.asDouble();
	}

	public String getString(String key) {
		Value v = this.get(key);
		if (v == null) {
			return "";
		}
		return v.asString();
	}

	public ValueArray getArray(String key) {
		Value v = this.get(key);
		if (v == null) {
			return null;
		}
		return v.asArray();
	}

	public ValueMap getMap(String key) {
		Value v = this.get(key);
		if (v == null) {
			return null;
		}
		return v.asMap();
	}

	public byte[] getBinary(String key) {
		Value v = this.get(key);
		if (v == null) {
			return null;
		}
		return v.asBinary();
	}

	public void set(String key, Value val) {
		if (this.data == null) {
			this.data = new TreeMap<String, Value>();
		}
		this.data.put(key, val);
	}

	public void put(String key, Object val) {
		Value v = Value.create(val);
		this.set(key, v);
	}

	public ValueMap createMap(String k) {
		ValueMap v = this.getMap(k);
		if (v != null) {
			return v;
		}
		v = new ValueMap(null);
		this.set(k, new Value(TypeConst.MAP, v));
		return v;
	}

	public ValueArray createArray(String k) {
		ValueArray v = this.getArray(k);
		if (v != null) {
			return v;
		}
		v = new ValueArray(null);
		this.set(k, new Value(TypeConst.ARRAY, v));
		return v;
	}

	public void remove(String k) {
		if (this.data == null) {
			return;
		}
		this.data.remove(k);
	}

	public void walk(ValueMapWalker w) {
		if (this.data != null) {
			for (Entry<String, Value> e : this.data.entrySet()) {
				if (w.walk(e.getKey(), e.getValue())) {
					return;
				}
			}
		}
	}

	public Map<String, Value> theData() {
		return this.data;
	}

	@Override
	public String toString() {
		Object m = this.toMap();
		if (m == null) {
			return "";
		}
		return m.toString();
	}
}
