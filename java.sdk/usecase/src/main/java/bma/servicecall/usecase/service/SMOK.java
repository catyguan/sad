package bma.servicecall.usecase.service;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;
import bma.servicecall.core.StatusConst;

public class SMOK extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		Answer a = new Answer();
		a.setStatus(StatusConst.OK);
		peer.writeAnswer(a, null);
	}
}
