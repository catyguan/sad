package bma.servicecall.usecase.service;

import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServiceMethod;

public abstract class BaseSM implements ServiceMethod {

	public static void dumpSM(Request req, Context ctx) {
		System.out.println("Request : " + req);
		System.out.println("Content : " + ctx);
		System.out.println("Deadline : " + ctx.getLong("Deadline"));
	}

}
