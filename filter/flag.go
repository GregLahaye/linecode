package filter

import (
	"flag"
	"strings"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

type Holder struct {
	Tags        arrayFlags
	Sort        *bool
	Limit       *int
	Easy        *bool
	NotEasy     *bool
	Medium      *bool
	NotMedium   *bool
	Hard        *bool
	NotHard     *bool
	Accepted    *bool
	NotAccepted *bool
	Starred     *bool
	NotStarred  *bool
	Paid        *bool
	NotPaid     *bool
}

func (h Holder) Parse() Filter {
	var f Filter
	f.Tags = h.Tags
	f.Sort = *h.Sort
	f.Limit = *h.Limit
	f.Easy = status(h.Easy, h.NotEasy)
	f.Medium = status(h.Medium, h.NotMedium)
	f.Hard = status(h.Hard, h.NotHard)
	f.Accepted = status(h.Accepted, h.NotAccepted)
	f.Starred = status(h.Starred, h.NotStarred)
	f.Paid = status(h.Paid, h.NotPaid)
	return f
}

func status(yes, no *bool) Status {
	if *yes {
		return Accept
	} else if *no {
		return Deny
	} else {
		return None
	}
}

func DefaultFlags(name string, h *Holder) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Var(&h.Tags, "t", "tags")
	h.Sort = fs.Bool("r", false, "sort by acceptance")
	h.Limit = fs.Int("l", -1, "limit, -1 for all")
	h.Easy = fs.Bool("e", false, "easy")
	h.NotEasy = fs.Bool("E", false, "not easy")
	h.Medium = fs.Bool("m", false, "medium")
	h.NotMedium = fs.Bool("M", false, "not medium")
	h.Hard = fs.Bool("h", false, "hard")
	h.NotHard = fs.Bool("H", false, "not hard")
	h.Accepted = fs.Bool("a", false, "accepted")
	h.NotAccepted = fs.Bool("A", false, "not accepted")
	h.Starred = fs.Bool("s", false, "starred")
	h.NotStarred = fs.Bool("S", false, "not starred")
	h.Paid = fs.Bool("p", false, "paid")
	h.NotPaid = fs.Bool("P", true, "not paid")
	return fs
}

func PreferencesFlags(name string, h *Holder) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Var(&h.Tags, "t", "tags")
	h.Sort = fs.Bool("r", true, "sort by acceptance")
	h.Limit = fs.Int("l", 20, "limit, -1 for all")
	h.Easy = fs.Bool("e", false, "easy")
	h.NotEasy = fs.Bool("E", false, "not easy")
	h.Medium = fs.Bool("m", false, "medium")
	h.NotMedium = fs.Bool("M", false, "not medium")
	h.Hard = fs.Bool("h", false, "hard")
	h.NotHard = fs.Bool("H", true, "not hard")
	h.Accepted = fs.Bool("a", false, "accepted")
	h.NotAccepted = fs.Bool("A", true, "not accepted")
	h.Starred = fs.Bool("s", false, "starred")
	h.NotStarred = fs.Bool("S", false, "not starred")
	h.Paid = fs.Bool("p", false, "paid")
	h.NotPaid = fs.Bool("P", true, "not paid")
	return fs
}
