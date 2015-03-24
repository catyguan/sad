package bma.servicecall.sockclient;

import bma.servicecall.core.ServiceConn;

public class Driver implements bma.servicecall.core.Driver {
	public static final String NAME = "socket";

	@Override
	public ServiceConn createConn(String type, String api) {
		SocketServiceConn o = new SocketServiceConn();
		return o;
	}
}
