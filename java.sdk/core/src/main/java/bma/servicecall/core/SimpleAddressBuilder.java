package bma.servicecall.core;

public class SimpleAddressBuilder implements AddressBuilder {

	private String type;
	private String api;
	private ValueMap option;

	@Override
	public Address build(String service, String method) {
		String s = this.api.replaceAll("\\$SNAME\\$", service);
		s = s.replaceAll("\\$MNAME\\$", method);
		Address o = new Address();
		o.setType(this.type);
		o.setApi(s);
		o.setOption(option);
		return o;
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

}
