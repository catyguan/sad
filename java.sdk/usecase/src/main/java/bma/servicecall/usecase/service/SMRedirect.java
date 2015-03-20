package bma.servicecall.usecase.service;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;

public class SMRedirect extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		String word = req.getString("word");

		String r = "Hello " + word;
		System.out.println(r);

		Answer a = new Answer();
		a.sureResult().put("Data", r);
		peer.writeAnswer(a, null);
	}
}
