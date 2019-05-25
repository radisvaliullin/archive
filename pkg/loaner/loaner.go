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

	// Parse input data
	// fill models from csv

	// banks
	banks := Banks{}
	err := csvToObjs(l.conf.InPath+"/banks.csv", &banks)
	if err != nil {
		log.Print("banks csv parse to obj err: ", err)
		return err
	}

	// facils
	facils := Facilities{}
	err = csvToObjs(l.conf.InPath+"/facilities.csv", &facils)
	if err != nil {
		log.Print("facilities csv parse to obj err: ", err)
		return err
	}
	// facility actual capacity map
	facilCap := FacilCapMap{}
	for _, f := range facils {
		facilCap[f.ID] = f.Amount
	}
	// sort facils
	sort.Sort(facils)

	// covens
	covens := Covenants{}
	err = csvToObjs(l.conf.InPath+"/covenants.csv", &covens)
	if err != nil {
		log.Print("covenants csv parse to obj err: ", err)
		return err
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
		return err
	}

	// handle loans
	assigns := Assignments{}
	facilYield := FacilYieldMap{}
	for _, l := range loans {

		// if loaned go to next loan
		loaned := false

		// find best facils
		for _, f := range facils {
			if loaned {
				break
			}
			// loan can't spend more then capacity of facility
			if float64(l.Amount) > facilCap[f.ID] {
				continue
			}

			// validate by covens of facility or bank
			fcs, ok := facilsCovensMap[f.ID]
			if ok {
				// if banned state
				if _, ok := fcs.BannedState[l.State]; ok {
					continue
				}
				// if out of likelihood max limit
				if fcs.MaxDefaultLikelihood != 0.0 &&
					l.DefaultLikelihood > fcs.MaxDefaultLikelihood {
					continue
				}
			}
			bcs, ok := banksCovensMap[f.BankID]
			if ok {
				// if banned state
				if _, ok := bcs.BannedState[l.State]; ok {
					continue
				}
				// if out of likelihood max limit
				if bcs.MaxDefaultLikelihood != 0.0 &&
					l.DefaultLikelihood > bcs.MaxDefaultLikelihood {
					continue
				}
			}

			// calc yield
			y := (1.0-l.DefaultLikelihood)*l.InterestRate*float64(l.Amount) - l.DefaultLikelihood*float64(l.Amount) - f.InterestRate*float64(l.Amount)
			// validate by yield
			if y < 0.0 {
				continue
			}
			assign := Assignment{FacilityID: f.ID, LoanID: l.ID}
			assigns = append(assigns, assign)
			loaned = true
			facilCap[f.ID] = facilCap[f.ID] - float64(l.Amount)
			facilYield[f.ID] = facilYield[f.ID] + y
			break
		}
	}

	// write result
	err = objsToCSV(l.conf.OutPath+"/assignments.csv", &assigns)
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
