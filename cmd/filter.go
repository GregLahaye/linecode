package cmd

import (
	"flag"
	"strings"
)

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *ArrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}


// this will be moved to different package
type Status int

const (
	Accept = iota
	Deny = iota
	None = iota
)

type Filter struct {
	Tags []string
	Easy Status
	Medium Status
	Hard Status
	Accepted Status
	Starred Status
	Paid Status
}
//

type filterHolder struct {
	Tags ArrayFlags
	Easy *bool
	NotEasy *bool
	Medium *bool
	NotMedium *bool
	Hard *bool
	NotHard *bool
	Accepted *bool
	NotAccepted *bool
	Starred *bool
	NotStarred *bool
	Paid *bool
	NotPaid *bool
}

func (fh filterHolder) parse() Filter {
	var f Filter
	f.Tags = fh.Tags
	f.Easy = filterStatus(fh.Easy, fh.NotEasy)
	f.Medium = filterStatus(fh.Medium, fh.NotMedium)
	f.Hard = filterStatus(fh.Hard, fh.Hard)
	f.Accepted = filterStatus(fh.Accepted, fh.NotAccepted)
	f.Starred = filterStatus(fh.Starred, fh.NotStarred)
	f.Paid = filterStatus(fh.Paid, fh.NotPaid)
	return f
}

func filterStatus(yes, no *bool) Status {
	if *yes {
		return Accept
	} else if *no {
		return Deny
	} else {
		return None
	}
}

func filterFlags(name string, fh *filterHolder) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Var(fh.Tags,"t", "tags")
	fh.Easy = fs.Bool("e", false, "easy")
	fh.NotEasy = fs.Bool("E", false, "not easy")
	fh.Medium = fs.Bool("m", false, "medium")
	fh.NotMedium = fs.Bool("M", false, "not medium")
	fh.Hard = fs.Bool("h", false, "hard")
	fh.NotHard = fs.Bool("H", false, "not hard")
	fh.Accepted = fs.Bool("a", false, "accepted")
	fh.NotAccepted = fs.Bool("A", false, "not accepted")
	fh.Starred = fs.Bool("s", false, "starred")
	fh.NotStarred = fs.Bool("S", false, "not starred")
	fh.Paid = fs.Bool("p", false, "paid")
	fh.NotPaid = fs.Bool("P", false, "not paid")
	return fs
}

var fh filterHolder
