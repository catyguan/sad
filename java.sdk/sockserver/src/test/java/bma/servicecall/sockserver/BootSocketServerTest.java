package bma.servicecall.sockserver;

import junit.framework.TestCase;

import org.springframework.context.support.FileSystemXmlApplicationContext;

import bma.servicecall.boot.SpringTestcaseUtil;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.JDKLogger;

public class BootSocketServerTest extends TestCase {

	FileSystemXmlApplicationContext context;

	@Override
	public void setUp() throws Exception {
		Debuger.init(new JDKLogger("test"));
		context = new SpringTestcaseUtil.ApplicationContextBuilder().project(
				"src/test/resources/spring_server.xml").build();
	}

	@Override
	public void tearDown() throws Exception {
		if (context != null)
			context.close();
	}

	/**
	 * 测试服务端启动
	 * 
	 * @throws Exception
	 */
	public void testServer() throws Exception {
		ServiceCallSocketServer server = context.getBean("server",
				ServiceCallSocketServer.class);
		server.startServer();
		try {
			Thread.sleep(60 * 1000);
		} finally {
			server.stopServer();
		}

	}

}
