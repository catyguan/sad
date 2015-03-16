package core

type ServiceMux struct {
	services map[string]ServiceObject
	methods  map[string]map[string]ServiceMethod
	backend  ServiceProvider
}

func (this *ServiceMux) SetBackend(v ServiceProvider) {
	this.backend = v
}

func (this *ServiceMux) SetServiceObject(name string, so ServiceObject) {
	if this.services == nil {
		this.services = make(map[string]ServiceObject)
	}
	this.services[name] = so
}

func (this *ServiceMux) SetServiceMethod(service string, method string, sm ServiceMethod) {
	if this.methods == nil {
		this.methods = make(map[string]map[string]ServiceMethod)
	}
	s, ok := this.methods[service]
	if !ok {
		s = make(map[string]ServiceMethod)
		this.methods[service] = s
	}
	s[method] = sm
}

func (this *ServiceMux) Find(s, m string) (ServiceMethod, error) {
	if ms, ok := this.methods[s]; ok {
		r := ms[m]
		if r != nil {
			return r, nil
		}
	}
	if ss, ok := this.services[s]; ok {
		return ss.GetMethod(m), nil
	}
	if this.backend != nil {
		return this.backend(s, m)
	}
	return nil, nil
}
