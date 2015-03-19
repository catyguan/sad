package bma.servicecall.usecase.invoke;

import bma.servicecall.core.Address;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Answer;
import bma.servicecall.core.Client;
import bma.servicecall.core.Context;
import bma.servicecall.core.Manager;
import bma.servicecall.core.Request;
import bma.servicecall.core.ValueMap;

public class SCITestAdd {

	static public void invoke(Manager manager, AddressBuilder ab, int a, int b,
			int c) {
		int z = 0;
		Client cl = manager.createClient();
		try {
			if (true) {
				Address addr = ab.build("test", "add");
				Request req = new Request();
				req.put("a", a);
				req.put("b", b);
				Context ctx = new Context();

				Answer answer = cl.invoke(addr, req, ctx);
				System.out.println(answer.toString());
				answer.checkError();
				if (answer.isDone()) {
					ValueMap rs = answer.getResult();
					if (rs != null) {
						z = rs.getInt("Data");
					} else {
						System.out.println("Result NULL");
						return;
					}
				} else {
					System.out.println("not done");
				}
			}
			if (true) {
				Address addr = ab.build("test", "add");
				Request req = new Request();
				req.put("a", c);
				req.put("b", z);
				Context ctx = new Context();

				Answer answer = cl.invoke(addr, req, ctx);
				System.out.println(answer.toString());
				answer.checkError();
				if (answer.isDone()) {
					ValueMap rs = answer.getResult();
					if (rs != null) {
						System.out.println("Result ==="+rs.getInt("Data"));
					} else {
						System.out.println("Result NULL");
						return;
					}
				} else {
					System.out.println("not done");
				}
			}
		} finally {
			cl.close();
		}
	}
}
