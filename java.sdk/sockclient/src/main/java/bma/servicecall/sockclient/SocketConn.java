package bma.servicecall.sockclient;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.SocketTimeoutException;

public class SocketConn {

	private String key;
	private Socket socket;
	private OutputStream out;
	private InputStream in;

	public String getKey() {
		return key;
	}

	public void setKey(String key) {
		this.key = key;
	}

	public Socket getSocket() {
		return socket;
	}

	public void setSocket(Socket socket) {
		this.socket = socket;
	}

	public OutputStream getOut() throws IOException {
		if (out == null) {
			out = this.socket.getOutputStream();
		}
		return out;
	}

	public InputStream getIn() throws IOException {
		if (in == null) {
			in = this.socket.getInputStream();
		}
		return in;
	}

	public void close() {
		try {
			if (this.in != null) {
				this.in.close();
				this.in = null;
			}
		} catch (Exception e) {
		}
		try {
			if (this.out != null) {
				this.out.close();
				this.out = null;
			}
		} catch (Exception e) {
		}
		try {
			if (this.socket != null) {
				this.socket.close();
				this.socket = null;
			}
		} catch (Exception e) {
		}
	}

	public boolean check() {
		try {
			this.socket.setSoTimeout(1);
			getIn().read();
			return false;
		} catch (SocketTimeoutException te) {
			return true;
		} catch (Exception e) {
			return false;
		}
	}

}
