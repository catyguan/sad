package bma.servicecall.core;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class ValueArray {

	private List<Value> data;

	public ValueArray(List<Value> d) {
		super();
		this.data = d;
	}

	public int len() {
		if (this.data == null)
			return 0;
		return this.data.size();
	}

	@SuppressWarnings("rawtypes")
	public static ValueArray create(List list) {
		List<Value> o = null;
		if (list != null) {
			o = new ArrayList<Value>();
			for (Object value : list) {
				o.add(Value.create(value));
			}
		}
		return new ValueArray(o);
	}

	public Iterator<Value> iterator() {
		if (this.data == null)
			return null;
		return this.data.iterator();
	}

	public List<Object> toArray() {
		return this.convertArray(DataConverter.noConverter);
	}

	public List<Object> convertArray(DataConverter dc) {
		List<Object> r = new ArrayList<Object>();
		if (this.data != null) {
			for (Value val : this.data) {
				if (val == null) {
					r.add(null);
				} else {
					r.add(val.convertValue(dc));
				}
			}
		}
		return r;
	}

	public Value get(int idx) {
		if (this.data == null) {
			return null;
		}
		if (idx < this.data.size()) {
			return this.data.get(idx);
		}
		return null;
	}

	public boolean set(int idx, Value v) {
		if (this.data == null) {
			return false;
		}
		if (idx < this.data.size()) {
			this.data.set(idx, v);
			return true;
		}
		return false;
	}

	public void add(Value v) {
		if (this.data == null) {
			this.data = new ArrayList<Value>();
		}
		this.data.add(v);
	}

	public void remove(int idx) {
		if (this.data == null) {
			return;
		}
		if (idx < 0 || idx >= this.data.size()) {
			return;
		}
		this.data.remove(idx);
	}

	public void walk(ValueArrayWalker w) {
		if (this.data == null) {
			int idx = 0;
			for (Value v : this.data) {
				if (w.walk(idx, v)) {
					return;
				}
				idx++;
			}

		}
	}

	public List<Value> theData() {
		return this.data;
	}
}
