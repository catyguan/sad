package bma.servicecall.usecase.service;

import java.util.Map;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;

public class SMEcho extends BaseSM {

	@SuppressWarnings("rawtypes")
	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		Map reqm = req.toMap();
		Answer a = new Answer();
		a.sureResult().put("Data", reqm);
		peer.writeAnswer(a, null);
	}
}
