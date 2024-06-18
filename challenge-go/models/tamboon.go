package models

type Tamboon struct {
	Name           string
	AmountSubunits int
	CCNumber       string
	CVV            string
	ExpMonth       string
	ExpYear        string
}

type Report struct {
	TotalReceived   int
	Success         int
	Failed          int
	Average         float32
	TopDonors       []string
	TopDonateAmount int
	TotalDonator    int
}
