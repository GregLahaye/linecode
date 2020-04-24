package filter

import (
	"github.com/GregLahaye/linecode/linecode"
	"sort"
)

type Status int

const (
	Accept = 0
	Deny   = 1
	None   = 2
)

type Filter struct {
	Tags     []string
	Sort     bool
	Limit    int
	Easy     Status
	Medium   Status
	Hard     Status
	Accepted Status
	Starred  Status
	Paid     Status
}

func Check(p linecode.Problem, tags []linecode.Tag, f Filter) bool {
	fail := false

	// check if there is a positive difficulty check filter
	d := f.Easy == Accept || f.Medium == Accept || f.Hard == Accept

	switch p.Difficulty.Level {
	case linecode.Easy:
		fail = shouldFail(true, f.Easy, d)
	case linecode.Medium:
		fail = shouldFail(true, f.Medium, d)
	case linecode.Hard:
		fail = shouldFail(true, f.Hard, d)
	}

	fail = fail || shouldFail(p.Starred, f.Starred, false)
	fail = fail || shouldFail(p.PaidOnly, f.Paid, false)
	fail = fail || shouldFail(p.Status == linecode.Accepted, f.Accepted, false)
	fail = fail || !hasAnyTag(p, f.Tags, tags)

	return !fail
}

func shouldFail(c bool, s Status, d bool) bool {
	switch s {
	case Accept:
		return !c
	case Deny:
		return c
	case None:
		return d
	}

	return false
}

func hasAnyTag(p linecode.Problem, slugs []string, tags []linecode.Tag) bool {
	if len(slugs) == 0 {
		return true
	}

	for _, slug := range slugs {
		if hasTag(p, slug, tags) {
			return true
		}
	}

	return false
}

func hasTag(p linecode.Problem, slug string, tags []linecode.Tag) bool {
	for _, tag := range tags {
		if tag.Slug == slug {
			for _, id := range tag.Questions {
				if id == p.Stat.ID {
					return true
				}
			}
		}
	}

	return false
}

func Array(a []linecode.Problem, t []linecode.Tag, f Filter) []linecode.Problem {
	var problems []linecode.Problem
	for _, p := range a {
		if Check(p, t, f) {
			problems = append(problems, p)
		}
	}

	if f.Sort {
		s := func(i, j int) bool {
			a := problems[i]
			b := problems[j]
			x := float64(a.Stat.TotalAccepted) / float64(a.Stat.TotalSubmitted)
			y := float64(b.Stat.TotalAccepted) / float64(b.Stat.TotalSubmitted)
			return x < y
		}

		sort.Slice(problems, s)
	}

	if f.Limit > -1 {
		i := len(problems) - f.Limit
		if i < 0 {
			i = 0
		}
		problems = problems[i:]
	}

	return problems
}
