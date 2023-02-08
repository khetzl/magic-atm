// IRL this package would implement the API for managing the physical cash store of an ATM machine.
// Ideally, this could be a pro-active component: letting the operator know when supplies are low, etc.
// for this excercise it's going to be a simple storage.

package cash

import("fmt")


// Stacks describes a state of the actual physical cash store, what store
// for IRL system this probably would replaced by an interface+map in their own module,
// since different machines might have different storages.
// we are also going to do assume the solution is for a single currency (GBP).
// building a general approach would need this structure to be a bit more flexible, since different
// currencies require different denominations, sometime quite funky numbers
// e.g. hungary -> 500, 1000, 2000, 5000, 10000, 20000. (500Ft is not stored in ATMs though)
//      trinidad and tobago -> 1, 5, 10, 20, 50, 100 (and usually stores $50 and $100, or often only$100)
type Stacks struct {
	Fivers uint16 // a fiver is a bit of an overkill, the only ATM paying out £5 bills
	Tenners uint16 // 16bit will be enough, there won't be 10k+ bills stored.
	Twenties uint16
	Fifties uint16
}

type CashStore struct {
	stacks Stacks
	access bool
}

func New() (*CashStore, error) {
	initStacks := Stacks{
		Fivers: 0,
		Tenners: 0,
		Twenties: 0,
		Fifties: 0,
	};
	c := &CashStore{
		stacks: initStacks,
		access: true,
	}
	return c, nil
}

// Operator physically refills, and rebalances the stores
func (c *CashStore) Refill(s Stacks) error {
	return nil
}

// Pay out a certain amount of cash.
func (c *CashStore) Payout(amount uint) (Stacks, error) {
	// Reuse validation to get the amount of different denominations to be released.
	// It might seem a waste (as no paralel user sessions will be using our system locally),
	// but to validate, practically, we will have to figure out how it's going to be paid out. i
	// t's going to be a relatively cheap operation, so why not return this information from
	// the validation, and have a bit of extra safety in this function call.
	s, error := c.ValidatePayout(amount)
	if (error != nil) {
		return s, error
	}
	error = c.doPayOut(s)
	return s, error
}

// Check if we can pay the amount out. Return an error if for some reason it's not going to
// work, tell what denominations can be used to pay out.
func (c *CashStore) ValidatePayout(amount uint) (Stacks, error) {
	// this is pretty much dummy for now
	if (amount % 5 == 0) {
		return Stacks{Fivers: uint16(amount/5)}, nil
	}
	return Stacks{}, fmt.Errorf("not divisor")
}

// Return the available smallest denomination. - This is informative so might just use numbers
func (c *CashStore) SmallestDenomination() (uint) {
	return 5
}

func (c *CashStore) doPayOut(s Stacks) error {
	return nil
}
