package loaner

import (
	"log"
	"math"
	"sort"
)

// Config loaner configs
type Config struct {
	InPath  string
	OutPath string
}

// Loaner loans by input date
type Loaner struct {
	conf Config
}

// New inits new Loaner
func New(conf Config) *Loaner {
	l := &Loaner{
		conf: conf,
	}
	return l
}

// Loan make loan from stream input
func (l *Loaner) Loan() error {

	// Parse input data from csv
	inputs, err := l.parseInputCSV()
	if err != nil {
		log.Print("loaner: parse input csv err: ", err)
		return err
	}

	// handle loans
	assigns := Assignments{}
	facilYield := FacilYieldMap{}
	for _, loan := range inputs.Loans {

		// find best facils (facils sorted by rate)
		for _, f := range inputs.Facils {

			// validate facility
			if ok := l.validateFacility(loan, f, inputs); !ok {
				continue
			}

			// calc yield
			y := l.calcYield(loan, f)

			// validate by yield
			if y < 0.0 {
				continue
			}

			// set assignment
			assign := Assignment{FacilityID: f.ID, LoanID: loan.ID}
			assigns = append(assigns, assign)

			// recalc facility capacity
			inputs.FacilCap[f.ID] = inputs.FacilCap[f.ID] - float64(loan.Amount)
			// add facility yield
			facilYield[f.ID] = facilYield[f.ID] + y
			break
		}
	}

	// write result
	err = l.writeResult(assigns, facilYield)
	if err != nil {
		log.Print("loaner: write result csv err: ", err)
		return err
	}

	return nil
}

func (l *Loaner) parseInputCSV() (*Inputs, error) {

	// banks
	banks := Banks{}
	err := csvToObjs(l.conf.InPath+"/banks.csv", &banks)
	if err != nil {
		log.Print("banks csv parse to obj err: ", err)
		return nil, err
	}

	// facils
	facils := Facilities{}
	err = csvToObjs(l.conf.InPath+"/facilities.csv", &facils)
	if err != nil {
		log.Print("facilities csv parse to obj err: ", err)
		return nil, err
	}
	// facility actual capacity map
	facilCap := FacilCapMap{}
	for _, f := range facils {
		facilCap[f.ID] = f.Amount
	}
	// sort facils by rate
	sort.Sort(facils)

	// covens
	covens := Covenants{}
	err = csvToObjs(l.conf.InPath+"/covenants.csv", &covens)
	if err != nil {
		log.Print("covenants csv parse to obj err: ", err)
		return nil, err
	}
	// covens by facils/banks maps
	facilsCovensMap := FacilsCovensMap{}
	banksCovensMap := BanksCovensMap{}
	for _, c := range covens {
		if c.FacilityID != 0 {
			if fc, ok := facilsCovensMap[c.FacilityID]; ok {
				if c.MaxDefaultLikelihood != 0 {
					fc.MaxDefaultLikelihood = c.MaxDefaultLikelihood
				}
				if c.BannedState != "" {
					fc.BannedState[c.BannedState] = struct{}{}
				}
			} else {
				nfc := FacilCovens{
					BankID:               c.BankID,
					MaxDefaultLikelihood: c.MaxDefaultLikelihood,
					BannedState:          map[string]struct{}{c.BannedState: struct{}{}},
				}
				facilsCovensMap[c.FacilityID] = &nfc
			}
		} else {
			if bc, ok := banksCovensMap[c.BankID]; ok {
				if c.MaxDefaultLikelihood != 0 {
					bc.MaxDefaultLikelihood = c.MaxDefaultLikelihood
				}
				if c.BannedState != "" {
					bc.BannedState[c.BannedState] = struct{}{}
				}
			} else {
				nbc := BankCovens{
					MaxDefaultLikelihood: c.MaxDefaultLikelihood,
					BannedState:          map[string]struct{}{c.BannedState: struct{}{}},
				}
				banksCovensMap[c.FacilityID] = &nbc
			}
		}
	}

	// input loans
	loans := Loans{}
	err = csvToObjs(l.conf.InPath+"/loans.csv", &loans)
	if err != nil {
		log.Print("loans csv parse to obj err: ", err)
		return nil, err
	}

	in := &Inputs{
		Facils:          facils,
		FacilCap:        facilCap,
		FacilsCovensMap: facilsCovensMap,
		BanksCovensMap:  banksCovensMap,
		Loans:           loans,
	}
	return in, nil
}

func (l *Loaner) validateFacility(loan Loan, f Facility, in *Inputs) bool {

	// loan can't spend more then capacity of facility
	if float64(loan.Amount) > in.FacilCap[f.ID] {
		return false
	}

	// validate by covens of facility or bank
	fcs, ok := in.FacilsCovensMap[f.ID]
	if ok {
		// if banned state
		if _, ok := fcs.BannedState[loan.State]; ok {
			return false
		}
		// if out of likelihood max limit
		if fcs.MaxDefaultLikelihood != 0.0 &&
			loan.DefaultLikelihood > fcs.MaxDefaultLikelihood {
			return false
		}
	}
	bcs, ok := in.BanksCovensMap[f.BankID]
	if ok {
		// if banned state
		if _, ok := bcs.BannedState[loan.State]; ok {
			return false
		}
		// if out of likelihood max limit
		if bcs.MaxDefaultLikelihood != 0.0 &&
			loan.DefaultLikelihood > bcs.MaxDefaultLikelihood {
			return false
		}
	}
	return true
}

func (l *Loaner) calcYield(loan Loan, f Facility) float64 {
	y := (1.0 - loan.DefaultLikelihood) * loan.InterestRate * float64(loan.Amount)
	y = y - loan.DefaultLikelihood*float64(loan.Amount) - f.InterestRate*float64(loan.Amount)
	return y
}

func (l *Loaner) writeResult(assigns Assignments, facilYield FacilYieldMap) error {

	err := objsToCSV(l.conf.OutPath+"/assignments.csv", &assigns)
	if err != nil {
		log.Print("assigs to csv parse err: ", err)
		return err
	}
	yields := Yields{}
	for fid, y := range facilYield {
		yield := Yield{
			FacilityID:    fid,
			ExpectedYield: int(math.Round(y)),
		}
		yields = append(yields, yield)
	}
	err = objsToCSV(l.conf.OutPath+"/yields.csv", &yields)
	if err != nil {
		log.Print("yields to csv parse err: ", err)
		return err
	}

	return nil
}
