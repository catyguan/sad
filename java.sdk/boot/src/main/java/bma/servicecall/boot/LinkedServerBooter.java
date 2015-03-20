package bma.servicecall.boot;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;

import bma.servicecall.core.Debuger;

public class LinkedServerBooter implements ServerBooter {

	private List<Object> mainList;

	public void setMain(Object obj) {
		if (this.mainList == null)
			this.mainList = new LinkedList<Object>();
		this.mainList.add(obj);
	}

	public void setMainList(List<Object> list) {
		this.mainList = list;
	}

	@Override
	public void startServer() {
		if (this.mainList != null) {
			for (Object obj : this.mainList) {
				if (obj instanceof ServerBooter) {
					if (Debuger.isEnable()) {
						Debuger.log("startup - " + obj);
					}
					ServerBooter sb = (ServerBooter) obj;
					sb.startServer();
				}
			}
		}
	}

	@Override
	public void stopServer() {
		if (this.mainList != null) {
			List<Object> tmp = new ArrayList<Object>(this.mainList);
			Collections.reverse(tmp);
			for (Object obj : this.mainList) {
				if (obj instanceof ServerBooter) {
					if (Debuger.isEnable()) {
						Debuger.log("stopdown - " + obj);
					}
					ServerBooter sb = (ServerBooter) obj;
					try {
						sb.stopServer();
					} catch (Exception e) {
					}
				}
			}
		}
	}

}
