package bma.servicecall.core;

public interface ValueMapWalker {
	/**
	 * 
	 * @param k
	 * @param v
	 * @return stop?
	 */
	public boolean walk(String k, Value v);
}
