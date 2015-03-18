package bma.servicecall.httpclient;

import org.codehaus.jackson.JsonParser;
import org.codehaus.jackson.map.ObjectMapper;

import bma.servicecall.core.ServiceConn;

public class Driver implements bma.servicecall.core.Driver {
	public static final String NAME = "http";

	private static final ObjectMapper mapper;

	static {
		mapper = new ObjectMapper();
		mapper.configure(JsonParser.Feature.ALLOW_SINGLE_QUOTES, true);
		mapper.configure(JsonParser.Feature.ALLOW_COMMENTS, true);
		mapper.configure(JsonParser.Feature.ALLOW_UNQUOTED_FIELD_NAMES, true);
	}

	public static ObjectMapper getDefaultMapper() {
		return mapper;
	}

	@Override
	public ServiceConn createConn(String type, String api) {
		HttpServiceConn o = new HttpServiceConn();
		return o;
	}
}
