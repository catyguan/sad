package bma.servicecall.core;

import java.util.Map;

public class Address {
	private String type;
	private String api;
	ValueMap option;

	public Address() {
		super();
	}

	public String getType() {
		return type;
	}

	public void setType(String type) {
		this.type = type;
	}

	public String getApi() {
		return api;
	}

	public void setApi(String api) {
		this.api = api;
	}

	public ValueMap getOption() {
		return option;
	}

	public void setOption(ValueMap option) {
		this.option = option;
	}

	public ValueMap sureOption() {
		if (this.option == null) {
			this.option = new ValueMap(null);
		}
		return this.option;
	}

	@SuppressWarnings("rawtypes")
	public static Address createAddress(String type, String api, Map opts) {
		Address o = new Address();
		o.type = type;
		o.api = api;
		o.option = ValueMap.create(opts);
		return o;
	}

	@SuppressWarnings("rawtypes")
	public static Address createAddressFromMap(Map vm) {
		return createAddressFromValue(ValueMap.create(vm));
	}

	public static Address createAddressFromValue(ValueMap vm) {
		Address o = new Address();
		o.type = vm.getString("Type");
		o.api = vm.getString("API");
		o.option = vm.getMap("Option");
		return o;
	}

	public void Valid() {
		if (this.type == "") {
			throw new AppError("address type empty");
		}
		if (this.api == "") {
			throw new AppError("address api empty");
		}
		return;
	}

	public ValueMap toValueMap() {
		ValueMap vm = new ValueMap(null);
		vm.put("Type", this.type);
		vm.put("API", this.api);
		if (this.option != null && this.option.len() > 0) {
			vm.put("Option", this.option);
		}
		return vm;
	}

	@Override
	public String toString() {
		StringBuffer buf = new StringBuffer();
		buf.append("Address(");
		buf.append("Type:").append(this.type).append(";");
		buf.append("API:").append(this.api).append(";");
		buf.append(")");
		return buf.toString();
	}
}
