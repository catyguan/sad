package bma.servicecall.httpclient;

import bma.servicecall.core.DataConverter;
import bma.servicecall.core.TypeConst;
import bma.servicecall.core.Util;

public class JsonConverter implements DataConverter {
	@Override
	public Object convert(int type, Object val) {
		if(type==TypeConst.BINARY) {
			if (val instanceof byte[]) {
				byte[] bs = (byte[]) val;
				return Util.byte2Hex(bs);
			}
		}
		return val;
	}
}
