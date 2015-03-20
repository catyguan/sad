package bma.servicecall.boot;

import java.lang.management.ManagementFactory;
import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;

import org.springframework.context.support.FileSystemXmlApplicationContext;

import bma.servicecall.core.Debuger;

public class ServerBoot {

	public static final String SERVER_XML = "classpath:spring_server.xml";
	public static final String BOOTER = "booter";

	public static void main(String[] args) {
		new ServerBoot()._main(args);
	}

	protected FileSystemXmlApplicationContext context;

	protected void _main(String[] args) {
		System.out.println(getClass().getName() + " start");

		String serverXml = SERVER_XML;
		String xml = System.getProperty("spring_server_xml");
		if (xml != null && !xml.isEmpty()) {
			serverXml = xml;
		}

		String booter = BOOTER;
		String tmp = System.getProperty("spring_server_booter");
		if (tmp != null && !tmp.isEmpty()) {
			booter = tmp;
		}
		try {
			boot(serverXml, booter);
		} catch (Throwable e) {
			System.err.println("xml=>" + serverXml + "  booter=>" + booter);
			e.printStackTrace(System.err);
		} finally {
			System.out.println(getClass().getName() + " exit");
		}

		// check threads
		try {
			ThreadMXBean tb = ManagementFactory.getThreadMXBean();
			long thisId = Thread.currentThread().getId();
			String[] names = new String[] { "Finalizer", "Reference Handler",
					"Signal Dispatcher" };
			long tids[] = tb.getAllThreadIds();
			for (long l : tids) {
				if (l == thisId)
					continue;
				ThreadInfo info = tb.getThreadInfo(l, 10);
				if (info == null)
					continue;
				boolean m = false;
				for (int i = 0; i < names.length; i++) {
					String n = names[i];
					if (n.equals(info.getThreadName())) {
						m = true;
						break;
					}
				}
				if (m)
					continue;

				System.out.println("LIVING THREAD!!\n" + info.toString());
			}
		} catch (Exception err) {
			System.out.println("check living thread fail");
			err.printStackTrace(System.out);
		}

		System.exit(1);
	}

	public void boot(String serverXml, String booter) throws Throwable {

		if (serverXml.isEmpty()) {
			throw new Exception("serverXml is miss!!");
		}
		System.out.println("loading spring application");
		context = new FileSystemXmlApplicationContext(serverXml);
		ServerBooter sbooter = null;
		if (context.containsBean(booter)) {
			sbooter = context.getBean(booter, ServerBooter.class);
			if (Debuger.isEnable()) {
				Debuger.log("" + sbooter);
			}
		}

		// 监听关闭
		final Thread t = Thread.currentThread();
		Runtime.getRuntime().addShutdownHook(new Thread() {
			@Override
			public void run() {
				if (t.isAlive()) {
					t.interrupt();
				}
			}
		});
		// 启动
		try {
			beforeStartBooter();
			if (sbooter != null) {
				try {
					System.out.println("start bootServers");
					sbooter.startServer();
				} catch (Throwable err) {
					System.err.println("start booter fail");
					err.printStackTrace(System.err);
					throw err;
				}
			}
			System.out.println("run  mainLoop");
			mainLoop();
			System.out.println("exit mainLoop");
		} catch (InterruptedException e) {
			System.out.println("booter mainLoop Interrupted!");
		} finally {
			if (sbooter != null) {
				try {
					System.out.println("stop bootServers");
					sbooter.stopServer();
				} catch (Exception e) {
				}
			}
			try {
				afterStopBooter();
			} catch (Exception e) {
			}
			System.out.println("close spring application");
			context.close();
			try {
				afterCloseContext();
			} catch (Exception e) {
			}
		}

	}

	protected void afterStopBooter() {

	}

	protected void afterCloseContext() {

	}

	protected void mainLoop() throws InterruptedException {
		while (true) {
			Thread.sleep(Long.MAX_VALUE);
		}
	}

	protected void beforeStartBooter() {

	}
}
