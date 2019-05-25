package loaner

//
// Simple input cvs structs
//

// Banks list of banks
type Banks []Bank

// Bank model
type Bank struct {
	ID   int    `csv:"id"`
	Name string `csv:"name"`
}

// Facilities list of facilities
type Facilities []Facility

// Len -
func (f Facilities) Len() int {
	return len(f)
}

// Less -
func (f Facilities) Less(i, j int) bool {
	return f[i].InterestRate < f[j].InterestRate
}

// Swap -
func (f Facilities) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Facility model
type Facility struct {
	ID           int     `csv:"id"`
	BankID       int     `csv:"bank_id"`
	InterestRate float64 `csv:"interest_rate"`
	Amount       float64 `csv:"amount"`
}

// FacilCapMap stores capacity of facility
type FacilCapMap map[int]float64

// FacilYieldMap stores total yield of facility
type FacilYieldMap map[int]float64

// Covenants list of covenants
type Covenants []Covenant

// Covenant model
type Covenant struct {
	BankID               int     `csv:"bank_id"`
	FacilityID           int     `csv:"facility_id"`
	MaxDefaultLikelihood float64 `csv:"max_default_likelihood"`
	BannedState          string  `csv:"banned_state"`
}

// FacilsCovensMap facilities covenants map
type FacilsCovensMap map[int][]FacilCoven

// FacilCoven model
type FacilCoven struct {
	BankID  int
	MaxDefaultLikelihood float64
	BannedState string
}

// Loans list of loans
type Loans []Loan

// Loan model
type Loan struct {
	ID                int     `csv:"id"`
	Amount            int     `csv:"amount"`
	InterestRate      float64 `csv:"interest_rate"`
	DefaultLikelihood float64 `csv:"default_likelihood"`
	State             string  `csv:"state"`
}

// Assignments list of assingments
type Assignments []Assignment

// Assignment model
type Assignment struct {
	LoanID int `csv:"loan_id"`
	FacilityID int `csv:"facility_id"`
}

// Yields list of yields
type Yields []Yield

// Yield model
type Yield struct {
	FacilityID int `csv:"facility_id"`
	ExpectedYield int `csv:"expected_yield"`
}
