package usecase

import (
	sccore "bma/servicecall/core"
	"fmt"
)

func SCIRedirect(m *sccore.Manager, ab sccore.AddressBuilder) error {
	cl := m.CreateClient()
	defer cl.Close()

	addr := ab("test", "redirect")
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
