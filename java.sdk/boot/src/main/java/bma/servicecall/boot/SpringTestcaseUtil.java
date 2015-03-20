package bma.servicecall.boot;

import java.io.File;
import java.net.URL;
import java.util.LinkedList;
import java.util.List;

import org.springframework.context.ApplicationContext;
import org.springframework.context.support.FileSystemXmlApplicationContext;

import bma.servicecall.core.AppError;

/**
 * Spring框架的测试用例工具
 * 
 * @author 关中
 * @since 1.0
 * 
 */
public class SpringTestcaseUtil {

	public static ApplicationContext projectContext(String[] names) {
		ApplicationContextBuilder b = new ApplicationContextBuilder();
		for (int i = 0; i < names.length; i++) {
			b.project(names[i]);
		}
		return b.build();
	}

	public static class ApplicationContextBuilder {

		private List<String> urlList = new LinkedList<String>();

		public ApplicationContextBuilder classpath(String classpath) {
			URL url = Thread.currentThread().getContextClassLoader()
					.getResource(classpath);
			urlList.add(url.toString());
			return this;
		}

		private URL getResource(String name) {
			URL url = Thread.currentThread().getContextClassLoader()
					.getResource(name);
			if (url == null) {
				throw new AppError("resource(" + name + ") invalid");
			}
			return url;
		}

		public ApplicationContextBuilder resource(String name) {
			URL url = getResource(name);
			urlList.add(url.toString());
			return this;
		}

		public ApplicationContextBuilder resource(Class<?> cls, String[] names) {
			for (String name : names) {
				URL url = getResource(name);
				urlList.add(url.toString());
			}
			return this;
		}

		protected File getUserDirFile(String subName) {
			return new File(System.getProperty("user.dir"), subName);
		}

		public ApplicationContextBuilder project(String name) {
			File file = getUserDirFile(name);
			if (file == null) {
				throw new AppError("file(" + name + ") invalid");
			}
			urlList.add(file.toURI().toString());
			return this;
		}

		public FileSystemXmlApplicationContext build() {
			return new FileSystemXmlApplicationContext(
					urlList.toArray(new String[0]));
		}
	}
}
