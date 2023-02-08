package main

import (
	"os"
	"fmt"
	"magic-atm/config"
	"magic-atm/cash"
	"magic-atm/network"
	"magic-atm/session"
	"magic-atm/magiccrypto"
)

// A real system would have the bootstrap sequence in the main. It would
// define the config, init the network, cashStore, and the userSession Manager.
// (Also, make sure that all peripherals are working fine )
// Once all started correctly, it would start the UI system, and the UI system
// would be linked to the control the system.
// However, we are going to just start the network, cashStore, and the manager, and
// play out a typical user scenario in the main.
func main() {
	conf := config.GetConfig()

	cashStore, err := cash.New()
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "INIT CASH error: %v\n", err)
		os.Exit(1)
	}

	err = cashStore.Refill(cash.Stacks{Fivers: uint16(100)})
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "FILL error: %v\n", err)
		os.Exit(1)
	}

	n := network.NewMagic()
	mgr, err := session.NewManager(n, cashStore, conf.Retries)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "MGR error: %v\n", err)
		os.Exit(1)
	}

	mgr.InsertCard("cardId")

	fmt.Println(" <<< typing pin")
	err = mgr.Authenticate(magiccrypto.HashPin("1234"))

	fmt.Println(" <<< user checking balance")
	bal, err := mgr.Balance()
	fmt.Println(" <<< user seeing balance:", bal)

	fmt.Println(" <<< UI request withdraw")
	err = mgr.Withdraw(100)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Couldn't withdraw due to error: %v\n", err)
		os.Exit(1)
	}

	mgr.Cancel()

	fmt.Println("Thank you for using MagicATM!.")

}
