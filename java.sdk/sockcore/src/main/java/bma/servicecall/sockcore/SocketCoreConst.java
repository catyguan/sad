package bma.servicecall.sockcore;

public interface SocketCoreConst {
	public static final int MT_END = 0;
	public static final int MT_MESSAGE_ID = 1;
	public static final int MT_REQUEST = 2;
	public static final int MT_ADDRESS = 3;
	public static final int MT_DATA = 4;
	public static final int MT_CONTEXT = 5;
	public static final int MT_ANSWER = 6;
	public static final int MT_MLINE = 7;
	public static final int MT_PING = 9;

	public static final int DATA_MAXLEN = 256 * 256 * 256;
	public static final int HEADER_SIZE = 4;
}
