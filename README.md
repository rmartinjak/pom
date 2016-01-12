# pom
Pomodoro timer

# Installation
```
go get github.com/rmartinjak/pom/pom
go get github.com/rmartinjak/pom/pomd
```

# Usage
Add `pomd &` (with appropriate flags, see `pomd help`) to your `.xinitrc` or
whatever. You can use the `-script=SCRIPT` option to make `pomd` run
`SCRIPT $state` whenever a state change occurs, e.g. to send a notification to your desktop.

`pomd` cycles the following 4 states:

* work pending
* work
* pause pending
* pause

The transitions *work pending* → *work* and *pause pending* → *pause* require
manual intervention in the form of invoking `pom next`. The other transitions
are done after a certain amount of time has elapsed, by default 25min work
time and 5min pause time.

Run `pom` to print the current state.
