package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
)

func SCITrans(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "login")

	key := ""
	if true {
		req := sccore.NewRequest()
		ctx := sccore.NewContext()
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		fmt.Println(answer.Dump())

		if answer.IsContinue() {
			rs := answer.GetResult()
			if rs != nil {
				key = rs.GetString("Key")
				fmt.Println("get login key -> ", key)
			}
		} else {
			fmt.Println("Invoke fail", answer.GetStatus())
			return nil
		}
	}

	if true {
		req := sccore.NewRequest()
		req.Put("User", "test")
		req.Put("Password", key)
		ctx := sccore.NewContext()
		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		fmt.Println(answer.Dump())
	}
	return nil
}
