package bma.servicecall.usecase.invoke;

import bma.servicecall.core.Address;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Answer;
import bma.servicecall.core.Client;
import bma.servicecall.core.Context;
import bma.servicecall.core.Manager;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;

public class SCIAsyncCallback {

	static public void invoke(Manager manager, AddressBuilder ab,
			Address callback) {
		if (callback == null) {
			callback = ab.build("test", "ok");
		}
		Client cl = manager.createClient();
		try {
			Address addr = ab.build("test", "async");
			Request req = new Request();
			Context ctx = new Context();
			ctx.put(PropertyConst.ASYNC_MODE, "callback");
			ctx.put(PropertyConst.CALLBACK, callback.toValueMap());

			Answer answer = cl.invoke(addr, req, ctx);
			System.out.println(answer.toString());
			answer.checkError();
			if (!answer.isAsync()) {
				System.err.println("must answer async");
				return;
			}
		} finally {
			cl.close();
		}
	}
}
