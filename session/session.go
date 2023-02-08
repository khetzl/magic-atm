package session

import(
	"fmt"
	"magic-atm/config"
	"magic-atm/network"
	"magic-atm/cash"
)

// No Interface needed for this, the user session is a fundamental part of the system
type FSMState int

const (
	Init         = 0
	Authenticate = 1
	Menu         = 2
	GoodBye      = 3
	Error        = 4
)

// Implements the user session FSM - pretty dumb
// IRL system probably would implement a separate process, that communicates
// with messages, and the public struct only holds onto the message queues.
// Ideally we can find a nice fitting library
type Manager struct {
	network network.GwApi
	cashStore *cash.CashStore
	state FSMState
	retries uint
	error string
	userToken []byte
}

//FIXME: error return should be refactored to also disconnect the network
func NewManager(network network.GwApi, c *cash.CashStore, retries uint) (*Manager, error) {
	m := &Manager{
		network: network,
		cashStore: c,
		state: Init,
		retries: retries,
		error: "",
	}
	return m, nil
}

func (m *Manager)OkError() error {
	m.state = Init
	return nil
}

func (m *Manager)InsertCard(_cardId string) error {
	err := m.requireState(Init)
	if (err != nil) { return err }

	fmt.Println("Processing card")
	conf := config.GetConfig()

	// Connect when new card is inserted
	err = m.network.Connect(conf.GatewayEndpoint)
	if (err != nil) {
		fmt.Println("ERROR: can't connect to network")
		return err
	}
	m.state = Authenticate
	return nil
}

func (m *Manager)Authenticate(pHash []byte) error {
	err := m.requireState(Authenticate)
	if (err != nil) { return err }

	token, err := m.network.Authorise(pHash)
	if (err != nil) {
		if (m.retries <= 1) {
			fmt.Println("Failed: out of tries")
			m.error = "out of tries"
			m.state = Error
		} else {
			fmt.Println("Failed Try again")
			m.retries -= 1
		}
		return err
	}
	fmt.Println("Welcome!")
	m.userToken = token
	m.state = Menu

	return nil
}

func (m *Manager)Withdraw(amount uint) error {
	err := m.requireState(Menu)
	if (err != nil) { return err }

	res, err := m.network.WithdrawRequest(m.userToken, amount)
	if (err != nil) { return err }

	if (res) {
		// success - return card - could tri
		_, cashError := m.cashStore.Payout(amount)
		if (cashError != nil) {
			fmt.Println("Error: couldn't sispensing cash... KKRRRBRRRzzzzz")
			return cashError
		}
		fmt.Println("Dispensing cash... BRRRRRRR")
		m.state = Init
	} else {
		fmt.Println("Sorry not enough balance.")
		// tell user balance wasn't enough. (even if we get balance
		// design choice: we let user to try again
	}
	m.network.Disconnect()
	return nil
}

func (m *Manager)Balance() (uint, error) {
	err := m.requireState(Menu)
	if (err != nil) { return 0, err }

	balance, err := m.network.GetBalance(m.userToken)
	if (err != nil) { return 0, err }

	return balance, nil
}

func (m *Manager) Cancel() error {
	// from any state, user can get their card back
	fmt.Println("Please take your card")
	m.network.Disconnect()
	m.state = Init
	return nil
}

func (m *Manager) requireState(s FSMState) error {
	if (m.state != s) {
		return fmt.Errorf("invalid state")
	}
	return nil
}
