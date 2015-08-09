package main

type params struct {
	params  map[string]string
	remote  string
	org     string
	user    string
	repo    string
	goos    string
	goarch  string
	version string
}

func newParams(p map[string]string) *params {
	return &params{
		params:  p,
		remote:  p["remote"],
		org:     p["org"],
		user:    p["user"],
		repo:    p["repo"],
		goos:    p["goos"],
		goarch:  p["goarch"],
		version: p["version"],
	}
}

func (p params) owner() string {
	var owner string
	if p.org != "" {
		owner = p.org + "/"
	}
	return owner + p.user
}
