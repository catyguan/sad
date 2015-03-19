package bma.servicecall.usecase.invoke;

import java.util.Date;

import bma.servicecall.core.Address;
import bma.servicecall.core.AddressBuilder;
import bma.servicecall.core.Answer;
import bma.servicecall.core.Client;
import bma.servicecall.core.Context;
import bma.servicecall.core.Manager;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ValueMap;

public class SCIAsyncPoll {

	static public void invoke(Manager manager, AddressBuilder ab) {
		Client cl = manager.createClient();
		try {
			Address addr = ab.build("test", "async");
			Request req = new Request();
			Context ctx = new Context();
			ctx.put(PropertyConst.ASYNC_MODE, "poll");

			Answer answer = cl.invoke(addr, req, ctx);
			System.out.println(answer.toString());
			answer.checkError();
			if (!answer.isAsync()) {
				System.err.println("must answer async");
				return;
			}
			
			Date et = new Date(new Date().getTime()+10*1000);
			Answer answer2 = cl.PollAnswer(addr, answer, ctx, et, 500);
			if(answer2==null) {
				System.err.println("poll timeout");
				return;			
			}
			
			System.out.println(answer2.toString());
			answer2.checkError();			
			
			if(!answer2.isDone()) {
				System.err.println("Answer fail - "+answer2.getStatus());
				return;
			}

			ValueMap rs = answer2.getResult();
			if(rs != null) {
				System.out.println("Result === "+rs);
			}
		} finally {
			cl.close();
		}
	}
}
