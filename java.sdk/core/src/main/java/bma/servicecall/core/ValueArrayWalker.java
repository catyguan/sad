package bma.servicecall.core;

public interface ValueArrayWalker {
	/**
	 * 
	 * @param idx
	 * @param v
	 * @return stop
	 */
	public boolean walk(int idx, Value v);
}
