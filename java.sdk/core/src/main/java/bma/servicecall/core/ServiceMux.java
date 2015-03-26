package bma.servicecall.core;

import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;

public class ServiceMux implements ServiceProvider {
	private Map<String, ServiceObject> services;
	private Map<String, Map<String, ServiceMethod>> methods;
	private ServiceProvider backend;

	public ServiceProvider getBackend() {
		return backend;
	}

	public void setBackend(ServiceProvider backend) {
		this.backend = backend;
	}

	public void setServiceObjects(Map<String, ServiceObject> map) {
		for (Entry<String, ServiceObject> e : map.entrySet()) {
			this.setServiceObject(e.getKey(), e.getValue());
		}
	}

	public void setServiceObject(String name, ServiceObject so) {
		if (this.services == null) {
			this.services = new HashMap<String, ServiceObject>();
		}
		this.services.put(name, so);
	}

	public void setServiceMethods(Map<String, ServiceMethod> map) {
		for (Entry<String, ServiceMethod> e : map.entrySet()) {
			String key = e.getKey();
			String[] sm = key.split(":", 2);
			this.setServiceMethod(sm[0], sm[1], e.getValue());
		}
	}

	public void setServiceMethod(String service, String method, ServiceMethod sm) {
		if (this.methods == null) {
			this.methods = new HashMap<String, Map<String, ServiceMethod>>();
		}
		Map<String, ServiceMethod> s = this.methods.get(service);
		if (s == null) {
			s = new HashMap<String, ServiceMethod>();
			this.methods.put(service, s);
		}
		s.put(method, sm);
	}

	public ServiceMethod find(String s, String m) {
		if (this.methods != null) {
			Map<String, ServiceMethod> ms = this.methods.get(s);
			if (ms != null) {
				ServiceMethod o = ms.get(m);
				if (o != null) {
					return o;
				}
			}
		}
		if (this.services != null) {
			ServiceObject ss = this.services.get(s);
			if (ss != null) {
				ServiceMethod o = ss.getMethod(m);
				if (o != null) {
					return o;
				}
			}
		}
		if (this.backend != null) {
			return this.backend.getServiceMethod(s, m);
		}
		return null;
	}

	public ServiceMethod getServiceMethod(String service, String method) {
		return this.find(service, method);
	}
}
