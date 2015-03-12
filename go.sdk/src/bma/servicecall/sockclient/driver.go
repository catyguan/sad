package sockclient

import sccore "bma/servicecall/core"

type SocketDriverFactory struct {
}

func (this *SocketDriverFactory) GetDriver(typ string) (sccore.Driver, error) {
	o := new(SocketDriver)
	return o, nil
}

func init() {
	df := new(SocketDriverFactory)
	sccore.InitDriverFactory(NAME_DRIVER, df)
}

type SocketDriver struct {
}

func (this *SocketDriver) Invoke(ictx sccore.InvokeContext, addr *sccore.Address, req *sccore.Request, ctx *sccore.Context) (*sccore.Answer, error) {
	// async := ctx.GetString(constv.KEY_ASYNC_MODE)
	// if async == "push" {
	// 	return nil, fmt.Errorf("http not support AsyncMode(%s)", async)
	// }
	// var reqm map[string]interface{}
	// if req == nil {
	// 	reqm = make(map[string]interface{})
	// } else {
	// 	reqm = req.ConvertMap(jsonConverter)
	// }

	// var ctxm map[string]interface{}
	// if ctx == nil {
	// 	ctxm = make(map[string]interface{})
	// } else {
	// 	ctxm = ctx.ConvertMap(jsonConverter)
	// }
	// opt := addr.GetOption()

	// reqbs, errE0 := json.Marshal(reqm)
	// if errE0 != nil {
	// 	return nil, errE0
	// }
	// ctxbs, errE1 := json.Marshal(ctxm)
	// if errE1 != nil {
	// 	return nil, errE1
	// }

	// var body io.Reader
	// qurl := addr.GetAPI()

	// data := make(url.Values)
	// data.Add("q", string(reqbs))
	// data.Add("c", string(ctxbs))

	// method := "POST"
	// body = strings.NewReader(data.Encode())

	// hreq, err2 := http.NewRequest(method, qurl, body)
	// if err2 != nil {
	// 	return nil, err2
	// }
	// hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// if opt != nil {
	// 	host := opt.GetString("Host")
	// 	if host != "" {
	// 		hreq.Header.Set("Host", host)
	// 	}
	// 	hs := opt.GetMap("Headers")
	// 	if hs != nil {
	// 		hs.Walk(func(k string, v *sccore.Value) bool {
	// 			s := v.AsString()
	// 			if s != "" {
	// 				hreq.Header.Set(k, s)
	// 			}
	// 			return false
	// 		})
	// 	}
	// }
	// dltm := ctx.GetLong(constv.KEY_DEADLINE)
	// dl := time.Unix(dltm, 0)
	// client := newHttpClient(dl)
	// sccore.DoLog("'%s' start", qurl)

	// hresp, err3 := client.Do(hreq)
	// if err3 != nil {
	// 	sccore.DoLog("'%s' fail '%s'", qurl, err3)
	// 	return nil, err3
	// }
	// sccore.DoLog("'%s' end '%d'", qurl, hresp.StatusCode)
	// defer hresp.Body.Close()
	// respBody, err4 := ioutil.ReadAll(hresp.Body)
	// if err4 != nil {
	// 	return nil, err4
	// }
	// content := string(respBody)
	// sccore.DoLog("'%s' --> %s", qurl, content)

	// a := sccore.NewAnswer()

	// switch hresp.StatusCode {
	// case 200:
	// 	m := make(map[string]interface{})
	// 	err5 := json.Unmarshal(respBody, &m)
	// 	if err5 != nil {
	// 		return nil, fmt.Errorf("decode response content fail - %s", content)
	// 	}
	// 	mm := sccore.CreateValueMap(m)
	// 	sc := mm.GetInt("Status")
	// 	if sc == 0 {
	// 		sc = 200
	// 	}
	// 	a.SetStatus(int(sc))
	// 	msg := mm.GetString("Message")
	// 	if msg == "" && sc == 200 {
	// 		msg = "OK"
	// 	}
	// 	a.SetMessage(msg)
	// 	rs := mm.GetMap("Result")
	// 	a.SetResult(rs)
	// 	actx := mm.GetMap("Context")
	// 	a.SetContext(actx)
	// case 301, 302:
	// 	a.SetStatus(302)
	// 	loc := hresp.Header.Get("Location")
	// 	if loc == "" {
	// 		a.SetStatus(502)
	// 		a.SetMessage("miss redirect location")
	// 	} else {
	// 		rs := sccore.NewValueMap(nil)
	// 		rs.Put("Type", NAME_DRIVER)
	// 		rs.Put("API", loc)
	// 		a.SetMessage("redirect")
	// 		a.SetResult(rs)
	// 	}
	// case 400, 404:
	// 	a.SetStatus(400)
	// 	a.SetMessage(content)
	// case 403:
	// 	a.SetStatus(403)
	// 	a.SetMessage(content)
	// case 504:
	// 	a.SetStatus(408)
	// 	a.SetMessage(content)
	// case 500:
	// 	a.SetStatus(500)
	// 	a.SetMessage(content)
	// default:
	// 	a.SetStatus(500)
	// 	a.SetMessage(fmt.Sprintf("unknow response code '%d'", hresp.StatusCode))
	// }
	// return a, nil
	return nil, nil
}