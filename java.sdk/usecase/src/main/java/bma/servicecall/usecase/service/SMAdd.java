package bma.servicecall.usecase.service;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;

public class SMAdd extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		int pa = req.getInt("a");
		int pb = req.getInt("b");
		int pc = pa + pb;

		System.out.println("a + b = " + pa + " + " + pb + " = " + pc);

		Answer a = new Answer();
		a.sureResult().put("Data", pc);
		peer.writeAnswer(a, null);
	}
}
