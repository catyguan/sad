package bma.servicecall.usecase.service;

import bma.servicecall.core.Answer;
import bma.servicecall.core.Context;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServicePeer;

public class SMAsync extends BaseSM {

	@Override
	public void execute(ServicePeer peer, Request req, Context ctx) {
		dumpSM(req, ctx);

		int sleepTime = req.getInt("sleep");
		if (sleepTime <= 0) {
			sleepTime = 3;
		}
		peer.sendAsync(ctx, null, 10 * 1000);
		final int fSleepTime = sleepTime;
		final ServicePeer fpeer = peer;
		new Thread(new Runnable() {

			@Override
			public void run() {
				try {
					Thread.sleep(fSleepTime * 1000);
				} catch (InterruptedException e) {
				}
				Answer a = new Answer();
				a.sureResult().put("Word", "Hello Kitty");
				fpeer.writeAnswer(a, null);
			}
		}).run();
	}
}
