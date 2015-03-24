package bma.servicecall.sockclient;

public class SocketPoolConfig {

	public void setPoolSize(int ps) {
		SocketPool.pool().config.poolSize = ps;
	}

	public void setPoolTimeoutMS(int ms) {
		SocketPool.pool().config.poolSize = ms;
	}

	public void setPoolIdleTimeMS(int ms) {
		SocketPool.pool().config.idleMS = ms;
	}

	public void close() {
		SocketPool.pool().close();
	}
}
