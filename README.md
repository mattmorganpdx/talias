# Talias
easier linux aliases

The purpose of talias is to allow users to create temporary aliases based on commands in their shell history. This works by creating a bash script in a known location in the users path. By default Talias will remove expired commands when it is run, but they can always be restored as long as you don't explicitly delete them.

## Quick Start

* Download the talias binary from the [release](https://github.com/mattmorganpdx/talias/releases/download/1.0-1/talias) page
* Place talias somewhere in your path and make sure it is executable, e.g.
  * `sudo cp ~/Downloads/talias /usr/local/bin/`
  * `sudo chmod +x /usr/local/bin/talias`
* Add these lines to your `~/.bashrc` so that the talias bin dir is in your path and your shell history is updated after each command.
    `export PATH=~/.talias/bin:$PATH`
    `export PROMPT_COMMAND="${PROMPT_COMMAND:+$PROMPT_COMMAND$'\n'}history -a; history -c; history -r"`
* Re-source your rc file with the command `source ~/.bashrc`
* You're ready to start taliasing - `talias -h`


## Building and using
Download go run `go build talias.go`

Add the talias bin dir to your path
eg `vi ~/.bashrc` and add this at the end
    export PATH=~/.talias/bin:$PATH
This will keep it set across shell sessions.
Then resource the file to get started
    source ~/.bashrc

## Usage
>Usage: talias \[OPTION]... \[ALIAS]
>> 	 -l 	 list aliases
> 	 -a 	 add or extend an alias
> 	 -d 	 delete an alias
> 	 -h 	 print usage message
> 	 -v 	 display version string

## Files Locations
- ~/.talias/bin/      - home of alias scripts
- ~/.talias/alias.db  - json meta data for aliases
- ~/.talias/talis/conf - talis config overrides
