package network

import (
	"fmt"
)

// We might want support another protocol, so leave our options open with
// adding an interface, and using that.
type GwApi interface {
	Connect(address string) error
	Authorise(pHash []byte) ([]byte, error)
	GetBalance(userToken []byte) (uint, error)
	WithdrawRequest(userToken []byte, amount uint) (bool, error)
	Disconnect() error
}

type MagicApi struct {}

func NewMagic() *MagicApi {
	m := &MagicApi{}
	return m
}

func (m *MagicApi) Connect(address string) error {
	fmt.Println("#connecting to:", address)
	return nil
}

func (m* MagicApi) Authorise(_pHash []byte) ([]byte, error) {
	token := []byte("token")
	return token, nil
}

func (m* MagicApi) GetBalance(userToken []byte) (uint, error) {
	return 200, nil
}


// Return whether a
func (m* MagicApi) WithdrawRequest(userToken []byte, amount uint) (bool, error) {
	return amount < 200, nil
}

func (m *MagicApi) Disconnect() error {
	fmt.Println("#disconnecting")
	return nil
}
