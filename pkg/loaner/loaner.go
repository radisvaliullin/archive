package loaner

import (
	"log"
	"sort"
	"math"
)

// Config loaner configs
type Config struct {
	InPath string
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

	// fill models from csv
	banks := Banks{}
	err := csvToObjs(l.conf.InPath+"/banks.csv", &banks)
	if err != nil {
		log.Print("banks csv parse to obj err: ", err)
		return err
	}

	facils := Facilities{}
	err = csvToObjs(l.conf.InPath+"/facilities.csv", &facils)
	if err != nil {
		log.Print("facilities csv parse to obj err: ", err)
		return err
	}
	facilCap := FacilCapMap{}
	for _, f := range facils {
		facilCap[f.ID]=f.Amount
	}
	sort.Sort(facils)
	
	covens := Covenants{}
	err = csvToObjs(l.conf.InPath+"/covenants.csv", &covens)
	if err != nil {
		log.Print("covenants csv parse to obj err: ", err)
		return err
	}
	facilsConvensMap := FacilsCovensMap{}
	for _, c := range covens {
		fc := FacilCoven{
			BankID:               c.BankID,
			MaxDefaultLikelihood: c.MaxDefaultLikelihood,
			BannedState:          c.BannedState,
		}
		if _, ok := facilsConvensMap[c.FacilityID]; ok {
			facilsConvensMap[c.FacilityID] = append(facilsConvensMap[c.FacilityID], fc)
		} else {
			facilsConvensMap[c.FacilityID] = []FacilCoven{fc}
		}
	}

	loans := Loans{}
	err = csvToObjs(l.conf.InPath+"/loans.csv", &loans)
	if err != nil {
		log.Print("loans csv parse to obj err: ", err)
		return err
	}

	// handle loans
	assigns := Assignments{}
	facilYield := FacilYieldMap{}
	// yields := Yields{}
	for _, l := range loans {

		loaned := false
		for _, f := range facils {
			if loaned {
				break
			}
			// TODO: need fix
			if float64(l.Amount) > facilCap[f.ID] {
				continue
			}

			cs, ok := facilsConvensMap[f.ID]
			if ok {

				for _, c := range cs {
					if l.State != c.BannedState &&
						c.MaxDefaultLikelihood != 0 &&
						l.DefaultLikelihood <= c.MaxDefaultLikelihood {

						y := (1.0-l.DefaultLikelihood)*l.InterestRate*float64(l.Amount) - l.DefaultLikelihood*float64(l.Amount) - f.InterestRate*float64(l.Amount)
						if y < 0.0 {
							continue
						}
						assign := Assignment{FacilityID: f.ID, LoanID: l.ID}
						assigns = append(assigns, assign)
						loaned = true
						facilCap[f.ID]= facilCap[f.ID]-float64(l.Amount)
						facilYield[f.ID]=facilYield[f.ID]+y
						break
					}
				}
			}
		}
	}
	// log.Printf("assignements %+v", assigns)
	// log.Printf("facil yilds %+v", facilYield)

	// write result
	err = objsToCSV(l.conf.OutPath+"/assignments.csv", &assigns)
	if err != nil {
		log.Print("assigs to csv parse err: ", err)
		return err
	}
	yields := Yields{}
	for fid, y := range facilYield {
		yield := Yield{
			FacilityID: fid,
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
