package main

type Filter struct {
	Tags     []string `json:"tags"`
	Easy     int      `json:"easy"`
	Medium   int      `json:"medium"`
	Hard     int      `json:"hard"`
	Accepted int      `json:"accepted"`
	Paid     int      `json:"paid"`
	Starred  int      `json:"starred"`
}

const accepted = "ac"

const (
	easy   = 1
	medium = 2
	hard   = 3
)

const (
	accept = iota
	deny   = iota
	none   = iota
)

func FilterProblem(p Problem, tags []Tag, f Filter) bool {
	fail := false

	d := f.Easy == accept || f.Medium == accept || f.Hard == accept
	switch p.Difficulty.Level {
	case easy:
		fail = ShouldFail(true, f.Easy, d)
	case medium:
		fail = ShouldFail(true, f.Medium, d)
	case hard:
		fail = ShouldFail(true, f.Hard, d)
	}

	fail = fail || ShouldFail(p.Starred, f.Starred, false)
	fail = fail || ShouldFail(p.PaidOnly, f.Paid, false)
	fail = fail || ShouldFail(p.Status == accepted, f.Accepted, false)
	fail = fail || !HasAnyTag(p, f.Tags, tags)

	return !fail
}

func ShouldFail(c bool, s int, d bool) bool {
	switch s {
	case accept:
		return !c
	case deny:
		return c
	case none:
		return d
	}

	return false
}

func HasAnyTag(p Problem, slugs []string, tags []Tag) bool {
	if len(slugs) == 0 {
		return true
	}

	for _, slug := range slugs {
		if HasTag(p, slug, tags) {
			return true
		}
	}

	return false
}

func HasTag(p Problem, slug string, tags []Tag) bool {
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
