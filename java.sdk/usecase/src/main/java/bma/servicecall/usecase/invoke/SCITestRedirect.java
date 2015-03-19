package bma.servicecall.usecase.invoke;

import bma.servicecall.core.Address;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Answer;
import bma.servicecall.core.Client;
import bma.servicecall.core.Context;
import bma.servicecall.core.Manager;
import bma.servicecall.core.Request;
import bma.servicecall.core.Util;
import bma.servicecall.core.ValueMap;

public class SCITestRedirect {

	static public void invoke(Manager manager, AddressBuilder ab, String word) {
		if (Util.empty(word)) {
			word = "Kitty";
		}
		Client cl = manager.createClient();
		try {
			Address addr = ab.build("test", "redirect");
			Request req = new Request();
			req.put("word", word);
			Context ctx = new Context();

			Answer answer = cl.invoke(addr, req, ctx);
			System.out.println(answer.toString());
			answer.checkError();
			if (answer.isDone()) {
				ValueMap rs = answer.getResult();
				if (rs != null) {
					System.out.println("Result === " + rs.toString());
				} else {
					System.out.println("Result NULL");
				}
			} else {
				System.out.println("not done");
			}
		} finally {
			cl.close();
		}
	}
}
