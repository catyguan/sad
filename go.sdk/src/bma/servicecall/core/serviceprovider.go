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
	if this.methods != nil {
		if ms, ok := this.methods[s]; ok {
			r := ms[m]
			if r != nil {
				return r, nil
			}
		}
	}
	if this.services != nil {
		if ss, ok := this.services[s]; ok {
			o := ss.GetMethod(m)
			if o != nil {
				return o, nil
			}
		}
	}
	if this.backend != nil {
		m, err := this.backend(s, m)
		if err != nil {
			return nil, err
		}
		if m != nil {
			return m, nil
		}
	}
	return nil, nil
}
