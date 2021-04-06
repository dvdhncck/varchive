# varchive




Installing goncurses

export CGO_CFLAGS_ALLOW=".*"
export CGO_LDFLAGS_ALLOW=".*"
go get github.com/rthornton128/goncurses
go install ~/projects/go/src/davidhancock.com/varchive

(see https://github.com/rthornton128/goncurses/wiki/KnownIssues)

actually, rather horribly after the update the go1.16, I had to 

go get github.com/rthornton128/goncurses

sudo ln -s /home/dave/projects/go/pkg/mod/pkg/mod/github.com/rthornton128/goncurses@v0.0.0-20210302221415-1355ee05acae /usr/local/go/src/goncurses
