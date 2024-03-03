package entities

import (
	"fmt"
	"strconv"
	"strings"
)

const DecimalPartLength = 100

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

//func (c *Currency) String() string {
//	return fmt.Sprintf("%d.%d", c.Whole, c.Decimal)
//}

func (c *Currency) Float() float32 {
	return float32(c.Whole) + float32(c.Decimal)/100
}

//func CurrencyFromString(s string) Currency {
//	whole, decimal := 0, 0
//
//	parts := strings.Split(s, ".")
//	if len(parts) == 1 {
//		whole, _ = strconv.Atoi(parts[0])
//	} else if len(parts) == 2 {
//		whole, _ = strconv.Atoi(parts[0])
//		decimal, _ = strconv.Atoi(parts[1])
//	}
//
//	return Currency{Whole: whole, Decimal: decimal}
//}

func CurrencyFromFloat(v float32) Currency {
	if v < 0 {
		return Currency{Whole: 0, Decimal: 0}
	}

	s := fmt.Sprintf("%.2f", v)
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		whole, _ := strconv.Atoi(parts[0])
		return Currency{Whole: whole, Decimal: 0}
	} else if len(parts) == 2 {
		whole, _ := strconv.Atoi(parts[0])
		decimal, _ := strconv.Atoi(parts[1])
		return Currency{Whole: whole, Decimal: decimal}
	}

	return Currency{Whole: 0, Decimal: 0}
}
