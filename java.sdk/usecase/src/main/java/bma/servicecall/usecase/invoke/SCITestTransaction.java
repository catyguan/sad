package bma.servicecall.usecase.invoke;

import bma.servicecall.core.Address;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Answer;
import bma.servicecall.core.Client;
import bma.servicecall.core.Context;
import bma.servicecall.core.Manager;
import bma.servicecall.core.Request;
import bma.servicecall.core.ValueMap;

public class SCITestTransaction {

	static public void invoke(Manager manager, AddressBuilder ab, String user) {
		String key = "";
		Client cl = manager.createClient();
		try {
			if (true) {
				Address addr = ab.build("test", "login");
				Request req = new Request();
				Context ctx = new Context();

				Answer answer = cl.invoke(addr, req, ctx);
				System.out.println(answer.toString());
				answer.checkError();
				if (answer.isDone()) {
					ValueMap rs = answer.getResult();
					if (rs != null) {
						key = rs.getString("Key");
					} else {
						System.out.println("Result NULL");
						return;
					}
				} else {
					System.out.println("not done");
					return;
				}
			}
			if (true) {
				Address addr = ab.build("test", "login");
				Request req = new Request();
				req.put("User", user);
				req.put("Password", key);
				Context ctx = new Context();

				Answer answer = cl.invoke(addr, req, ctx);
				System.out.println(answer.toString());
				answer.checkError();
			}
		} finally {
			cl.close();
		}
	}
}
