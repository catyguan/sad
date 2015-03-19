package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
)

func SCIHello(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "hello")
	req := sccore.NewRequest()
	req.Put("word", "Kitty")
	ctx := sccore.NewContext()

	answer, err := cl.Invoke(addr, req, ctx)
	if err != nil {
		return fmt.Errorf("invoke fail - %s", err)
	}
	fmt.Println(answer.Dump())

	if answer.IsDone() {
		rs := answer.GetResult()
		if rs != nil {
			fmt.Println("RESULT ===", rs.Dump())
		} else {
			fmt.Println("RESULT NULL")
		}
	} else {
		fmt.Println("not done")
	}
	return nil
}

func SCIBinary(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "echo")
	req := sccore.NewRequest()
	req.Put("binary", []byte("Kitty"))
	ctx := sccore.NewContext()

	answer, err := cl.Invoke(addr, req, ctx)
	if err != nil {
		return fmt.Errorf("invoke fail - %s", err)
	}
	fmt.Println(answer.Dump())

	if answer.IsDone() {
		rs := answer.GetResult()
		if rs != nil {
			dat := rs.GetMap("Data")
			fmt.Println("RESULT ===", dat.GetBinary("binary"))
		} else {
			fmt.Println("RESULT NULL")
		}
	} else {
		fmt.Println("not done")
	}
	return nil
}

func SCIAdd(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	c := int32(0)
	if true {
		addr := ab("test", "add")
		req := sccore.NewRequest()
		req.Put("a", 1)
		req.Put("b", 2)
		ctx := sccore.NewContext()

		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		fmt.Println(answer.Dump())

		if answer.IsDone() {
			rs := answer.SureResult()
			c = rs.GetInt("Data")
		} else {
			fmt.Println("not done")
			return nil
		}
	}

	if true {
		addr := ab("test", "add")
		req := sccore.NewRequest()
		req.Put("a", c)
		req.Put("b", 3)
		ctx := sccore.NewContext()

		answer, err := cl.Invoke(addr, req, ctx)
		if err != nil {
			return fmt.Errorf("invoke fail - %s", err)
		}
		fmt.Println(answer.Dump())
	}
	return nil
}
