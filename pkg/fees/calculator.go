package fees

import (
	"github.com/shopspring/decimal"
)

func CalculateFees(stockPrice, quantity decimal.Decimal) decimal.Decimal {
	transactionValue := stockPrice.Mul(quantity)

	brokerage := transactionValue.Mul(decimal.NewFromFloat(0.0003))
	if brokerage.LessThan(decimal.NewFromInt(20)) {
		brokerage = decimal.NewFromInt(20)
	}

	stt := transactionValue.Mul(decimal.NewFromFloat(0.00025))

	gst := brokerage.Mul(decimal.NewFromFloat(0.18))

	exchangeCharges := transactionValue.Mul(decimal.NewFromFloat(0.0000325))

	sebiCharges := transactionValue.Mul(decimal.NewFromFloat(0.000001))

	stampDuty := transactionValue.Mul(decimal.NewFromFloat(0.00003))

	totalFees := brokerage.Add(stt).Add(gst).Add(exchangeCharges).Add(sebiCharges).Add(stampDuty)

	return totalFees.Round(2)
}

