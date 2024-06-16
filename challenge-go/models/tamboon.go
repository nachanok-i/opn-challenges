package models

type Tamboon struct {
	Name           string
	AmountSubunits int
	CCNumber       string
	CVV            string
	ExpMonth       string
	ExpYear        string
}

type ChargeResponse struct {
	Status   string
	ChargeId string
}
