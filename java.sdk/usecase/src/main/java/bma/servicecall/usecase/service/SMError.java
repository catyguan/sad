package bma.servicecall.usecase.service;

import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;
import bma.servicecall.core.Util;

public class SMError extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		String emsg = req.getString("Error");
		if (Util.empty(emsg)) {
			emsg = "<test error>";
		}
		peer.writeAnswer(null, new AppError(emsg));
	}
}
