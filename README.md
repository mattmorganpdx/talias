# talias
easier linux aliases

## Building and using
Download go run `go build talias.go`

Add the talias bin dir to your path
eg `vi ~/.bashrc` and add this at the end
    export PATH=~/.talias/bin:$PATH
This will keep it set across shell sessions.
Then resource the file to get started
    source ~/.bashrc

## Talias commands
talias -l | list your recent history
talias -L | list your aliases
talias -a <alias> | show history so you can add alias

## Files
~/.talias/bin/      | location of alias scripts
~/.talias/alias.db  | json meta data for aliases
