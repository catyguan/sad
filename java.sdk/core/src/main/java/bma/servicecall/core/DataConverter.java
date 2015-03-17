package bma.servicecall.core;

public interface DataConverter {
	public Object convert(int type, Object val);
	
	public final static DataConverter noConverter = new DataConverter() {

		public Object convert(int type, Object val) {
			return val;
		}
	};

}
