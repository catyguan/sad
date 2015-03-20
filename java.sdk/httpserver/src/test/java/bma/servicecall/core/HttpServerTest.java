package bma.servicecall.core;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;
import bma.servicecall.httpserver.ServiceCallWebServer;
import bma.servicecall.usecase.service.SMAdd;
import bma.servicecall.usecase.service.SMAsync;
import bma.servicecall.usecase.service.SMEcho;
import bma.servicecall.usecase.service.SMError;
import bma.servicecall.usecase.service.SMHello;
import bma.servicecall.usecase.service.SMLogin;
import bma.servicecall.usecase.service.SMOK;
import bma.servicecall.usecase.service.SMRedirect;

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
		Manager manager = new Manager("callback");

		ServiceMux mux = new ServiceMux();
		mux.setServiceMethod("test", "echo", new SMEcho());
		mux.setServiceMethod("test", "ok", new SMOK());
		mux.setServiceMethod("test", "hello", new SMHello());
		mux.setServiceMethod("test", "add", new SMAdd());
		mux.setServiceMethod("test", "error", new SMError());
		mux.setServiceMethod("test", "redirect", new SMRedirect());
		mux.setServiceMethod("test", "login", new SMLogin());
		mux.setServiceMethod("test", "async", new SMAsync());

		ServiceCallWebServer server = new ServiceCallWebServer();
		server.setServiceMux(mux);
		server.setClientFactory(manager);
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
