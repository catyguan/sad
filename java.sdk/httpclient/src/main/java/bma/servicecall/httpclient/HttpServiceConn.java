package bma.servicecall.httpclient;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.client.methods.HttpUriRequest;
import org.apache.http.client.params.HttpClientParams;
import org.apache.http.client.utils.URLEncodedUtils;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.HttpConnectionParams;
import org.apache.http.params.HttpParams;
import org.apache.http.util.EntityUtils;
import org.codehaus.jackson.map.ObjectMapper;

import bma.servicecall.core.Address;
import bma.servicecall.core.Answer;
import bma.servicecall.core.AppError;
import bma.servicecall.core.Context;
import bma.servicecall.core.DataConverter;
import bma.servicecall.core.Debuger;
import bma.servicecall.core.InvokeContext;
import bma.servicecall.core.PropertyConst;
import bma.servicecall.core.Request;
import bma.servicecall.core.ServiceConn;
import bma.servicecall.core.Util;
import bma.servicecall.core.Value;
import bma.servicecall.core.ValueMap;
import bma.servicecall.core.ValueMapWalker;

public class HttpServiceConn implements ServiceConn {

	private static DataConverter jsonConverter = new JsonConverter();

	private String transId;

	public Answer invoke(InvokeContext ictx, Address addr, Request req,
			Context ctx) {
		try {
			return doInvoke(ictx, addr, req, ctx);
		} catch (Exception ex) {
			throw AppError.handle(ex);
		}
	}

	public HttpResp postContent(String url, List<NameValuePair> params,
			List<NameValuePair> headers, int timeoutMS)
			throws ClientProtocolException, IOException {
		String pstr = "";
		StringEntity post = null;
		if (params != null && !params.isEmpty()) {
			pstr = URLEncodedUtils.format(params, "UTF-8");
		}
		post = new StringEntity(pstr, "UTF-8");

		HttpUriRequest req = null;
		HttpPost preq = new HttpPost(url);
		preq.addHeader("Content-Type", "application/x-www-form-urlencoded");
		preq.setEntity(post);
		req = preq;

		HttpParams hparams = new BasicHttpParams();
		HttpConnectionParams.setConnectionTimeout(hparams, timeoutMS);
		HttpConnectionParams.setSoTimeout(hparams, timeoutMS);
		HttpClientParams.setRedirecting(hparams, false);

		DefaultHttpClient client = new DefaultHttpClient(hparams);
		HttpResponse resp = null;
		try {
			resp = client.execute(req);
			HttpResp r = new HttpResp();
			r.Status = resp.getStatusLine().getStatusCode();
			Header h = resp.getFirstHeader("Location");
			if (h != null) {
				r.Location = h.getValue();
			}
			HttpEntity entity = resp.getEntity();
			if (entity != null) {
				r.Content = EntityUtils.toString(entity);
			}
			return r;
		} catch (Exception e) {
			throw AppError.handle(e);
		}
	}

	@SuppressWarnings({ "rawtypes", "unchecked" })
	public Answer doInvoke(InvokeContext ictx, Address addr, Request req,
			Context ctx) throws Exception {
		String async = null;
		if (ctx != null) {
			async = ctx.getString(PropertyConst.ASYNC_MODE);
		}
		if (async != null && async.equals("push")) {
			throw new AppError("http not support AsyncMode(push)");
		}

		Map reqm;
		if (req == null) {
			reqm = new TreeMap();
		} else {
			reqm = req.convertMap(jsonConverter);
		}
		Map ctxm;
		if (ctx == null) {
			ctxm = new TreeMap();
		} else {
			ctxm = ctx.convertMap(jsonConverter);
		}
		if (!Util.empty(this.transId)) {
			ctxm.put(PropertyConst.TRANSACTION_ID, this.transId);
		}
		ValueMap opt = addr.getOption();

		ObjectMapper om = Driver.getDefaultMapper();
		String reqs = om.writeValueAsString(reqm);
		String ctxs = om.writeValueAsString(ctxm);

		List<NameValuePair> params = new ArrayList<NameValuePair>();
		params.add(new BasicNameValuePair("q", reqs));
		params.add(new BasicNameValuePair("c", ctxs));
		List<NameValuePair> headers = new ArrayList<NameValuePair>();
		headers.add(new BasicNameValuePair("Content-Type",
				"application/x-www-form-urlencoded"));
		if (opt != null) {
			String host = opt.getString("Host");
			if (!Util.empty(host)) {
				headers.add(new BasicNameValuePair("Host", host));
			}
			ValueMap hs = opt.getMap("Headers");
			if (hs != null) {
				final List<NameValuePair> fhs = headers;
				hs.walk(new ValueMapWalker() {

					@Override
					public boolean walk(String k, Value v) {
						fhs.add(new BasicNameValuePair(k, v.asString()));
						return false;
					}
				});
			}
		}

		String qurl = addr.getApi();
		long dlt = ctx.getLong(PropertyConst.DEADLINE);
		long now = new Date().getTime();
		Debuger.log("'" + qurl + "' start");
		HttpResp resp = postContent(qurl, params, headers,
				(int) ((dlt - now) * 1000));
		Debuger.log("'" + qurl + "' end '" + resp.Status + "' - '"
				+ resp.Content + "'");

		Answer a = new Answer();

		switch (resp.Status) {
		case 200: {
			Map m = om.readValue(resp.Content, Map.class);
			ValueMap mm = ValueMap.create(m);
			int sc = mm.getInt("Status");
			if (sc == 0) {
				sc = 200;
			}
			a.setStatus(sc);
			String msg = mm.getString("Message");
			if (Util.empty(msg) && sc == 200) {
				msg = "OK";
			}
			a.setMessage(msg);
			ValueMap rs = mm.getMap("Result");
			a.setResult(rs);
			ValueMap actx = mm.getMap("Context");
			a.setContext(actx);
		}
			break;
		case 301:
		case 302: {
			a.setStatus(302);
			String loc = resp.Location;
			if (Util.empty(loc)) {
				a.setStatus(502);
				a.setMessage("miss redirect location");
			} else {
				ValueMap rs = new ValueMap(null);
				rs.put("Type", Driver.NAME);
				rs.put("API", loc);
				a.setMessage("redirect");
				a.setResult(rs);
			}
		}
			break;
		case 400:
		case 404:
			a.setStatus(400);
			a.setMessage(resp.Content);
			break;
		case 403:
			a.setStatus(403);
			a.setMessage(resp.Content);
			break;
		case 504:
			a.setStatus(408);
			a.setMessage(resp.Content);
			break;
		case 500:
			a.setStatus(500);
			a.setMessage(resp.Content);
			break;
		default:
			a.setStatus(500);
			a.setMessage("unknow response code '" + resp.Status + "'");
			break;
		}
		if (a.getStatus() != 100) {
			this.transId = "";
		} else {
			ValueMap ctx2 = a.getContext();
			if (ctx2 != null && ctx2.has(PropertyConst.TRANSACTION_ID)) {
				this.transId = ctx2.getString(PropertyConst.TRANSACTION_ID);
			}
		}
		return a;
	}

	@Override
	public Answer waitAnswer(int timeoutMS) {
		throw new AppError("http not support waitAnswer");
	}

	public void clear() {
		this.transId = "";
	}

	@Override
	public void end() {
		this.clear();
	}

	@Override
	public void close() {
		this.clear();
	}

}
