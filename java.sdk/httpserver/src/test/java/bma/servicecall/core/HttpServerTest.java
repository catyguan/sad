package bma.servicecall.core;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;
import bma.servicecall.httpserver.ServiceCallWebServer;

public class HttpServerTest extends TestCase {
	/**
	 * Create the test case
	 * 
	 * @param testName
	 *            name of the test case
	 */
	public HttpServerTest(String testName) {
		super(testName);
	}

	/**
	 * @return the suite of tests being tested
	 */
	public static Test suite() {
		return new TestSuite(HttpServerTest.class);
	}

	@Override
	protected void setUp() throws Exception {
		Debuger.init(new JDKLogger("test"));
	}

	public void testServer() throws Exception {
		ServiceMux mux = new ServiceMux();
		mux.setServiceMethod("test", "hello", new ServiceMethod() {

			@Override
			public void execute(ServicePeer peer, Request req, Context ctx) {

			}
		});

		ServiceCallWebServer server = new ServiceCallWebServer();
		server.setServiceMux(mux);
		server.setPort(1080);
		server.setLog(true);
		server.startServer();
		try {
			Thread.sleep(60 * 1000);
		} finally {
			server.stopServer();
		}
	}
}
