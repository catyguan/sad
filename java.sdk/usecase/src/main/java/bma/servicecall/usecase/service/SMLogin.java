package bma.servicecall.usecase.service;

import java.util.Random;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;
import bma.servicecall.core.ServiceRequest;
import bma.servicecall.core.StatusConst;
import bma.servicecall.core.ValueMap;

public class SMLogin extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		peer.beginTransaction();
		String key = "";
		if (true) {
			dumpSM(req, ctx);

			Random ro = new Random(System.currentTimeMillis());
			key = "" + ro.nextInt(99999999 - 10000000) + 10000000;

			Answer a = new Answer();
			a.setStatus(StatusConst.CONTINUE);
			ValueMap rs = a.sureResult();
			rs.put("Key", key);

			peer.writeAnswer(a, null);
		}

		if (true) {
			ServiceRequest sr = peer.readRequest(30 * 1000);
			Request req2 = sr.getRequest();
			Context ctx2 = sr.getContext();
			dumpSM(req2, ctx2);

			String user = req2.getString("User");
			String pass = req2.getString("Password");
			System.out.println("param " + user + "," + pass);
			Answer a = new Answer();
			ValueMap rs = a.sureResult();
			if (user != null && pass != null && user.equals("test")
					&& pass.equals(key)) {
				rs.put("Done", true);
				rs.put("UID", 12345);
			} else {
				rs.put("Done", false);
				rs.put("Why", "user or pass invalid");
			}
			peer.writeAnswer(a, null);
		}
	}
}
