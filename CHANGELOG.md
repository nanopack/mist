## v1.0.2 (March 17, 2016)

BUG FIXES:
  - Updated how Mist parses configs.

IMPROVEMENTS:
  - Added additional options for specifying logging type (stdout, file) and an option
  for specifying the log path.

## v1.0.1 (March 10, 2016)

BUG FIXES:
  - Break if mist fails to write to a web socket client; it's presumably dead, and
  we don't want mist looping forever trying to read/write to something git will never
  be able to.

## Previous (March 9, 2016)

This change log began with version 1.0.0. Any prior changes can be seen by viewing
the commit history.
