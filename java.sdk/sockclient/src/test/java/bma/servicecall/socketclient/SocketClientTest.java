package bma.servicecall.socketclient;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.JDKLogger;
import bma.servicecall.core.Manager;
import bma.servicecall.core.SimpleAddressBuilder;
import bma.servicecall.sockclient.SocketPoolConfig;
import bma.servicecall.usecase.invoke.SCIAsyncCallback;
import bma.servicecall.usecase.invoke.SCIAsyncPoll;
import bma.servicecall.usecase.invoke.SCIAsyncPush;
import bma.servicecall.usecase.invoke.SCITestAdd;
import bma.servicecall.usecase.invoke.SCITestBinary;
import bma.servicecall.usecase.invoke.SCITestHello;
import bma.servicecall.usecase.invoke.SCITestRedirect;
import bma.servicecall.usecase.invoke.SCITestTransaction;

/**
 * Unit test for simple App.
 */
public class SocketClientTest extends TestCase {
	/**
	 * Create the test case
	 * 
	 * @param testName
	 *            name of the test case
	 */
	public SocketClientTest(String testName) {
		super(testName);
	}

	/**
	 * @return the suite of tests being tested
	 */
	public static Test suite() {
		return new TestSuite(SocketClientTest.class);
	}

	public AddressBuilder builder() {
		SimpleAddressBuilder ab = new SimpleAddressBuilder();
		ab.setType("socket");
		ab.setApi("tcp:localhost:1080:$SNAME$:$MNAME$");
		return ab;
	}

	private Manager m;

	@Override
	protected void setUp() throws Exception {
		super.setUp();
		Debuger.init(new JDKLogger("test"));
		m = new Manager("test");
		SocketPoolConfig cfg = new SocketPoolConfig();
		cfg.setPoolSize(3);
	}

	@Override
	protected void tearDown() throws Exception {
		super.tearDown();
		new SocketPoolConfig().close();
	}

	public void testHello() throws Exception {
		SCITestHello.invoke(m, builder(), null);
	}

	public void testIdle() throws Exception {
		SCITestHello.invoke(m, builder(), null);
		Thread.sleep(61 * 1000);
	}

	public void testBinary() throws Exception {
		SCITestBinary.invoke(m, builder(), null);
	}

	public void testAdd() throws Exception {
		SCITestAdd.invoke(m, builder(), 1, 2, 3);
	}

	public void testTrasaction() throws Exception {
		SCITestTransaction.invoke(m, builder(), "test");
	}

	public void testRedirect() throws Exception {
		SCITestRedirect.invoke(m, builder(), null);
	}

	public void testAsyncPoll() throws Exception {
		SCIAsyncPoll.invoke(m, builder());
	}

	public void testAsyncPush() throws Exception {
		SCIAsyncPush.invoke(m, builder());
	}

	public void testAsyncCallback() throws Exception {
		SCIAsyncCallback.invoke(m, builder(), null);
	}
}
