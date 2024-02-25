package entities

import "fmt"

type Currency struct {
	Whole   int
	Decimal int
}

func (c *Currency) Validate() error {
	if c.Whole < 0 || c.Decimal < 0 {
		return ValidationError{Err: fmt.Errorf("invalid currency: %d.%d", c.Whole, c.Decimal)}
	}
	return nil
}

func (c *Currency) Add(add *Currency) {
	c.Whole += add.Whole
	c.Decimal += add.Decimal

	if c.Decimal >= 100 {
		c.Whole++
		c.Decimal -= 100
	}
}

func (c *Currency) Sub(sub *Currency) error {
	whole := c.Whole
	decimal := c.Decimal

	whole -= sub.Whole
	decimal -= sub.Decimal

	if decimal < 0 {
		whole--
		decimal += 100
	}

	if whole < 0 {
		return OutOfBalanceError{}
	}

	c.Whole = whole
	c.Decimal = decimal

	return nil
}
