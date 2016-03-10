## v1.0.1 (March 10, 2016)

BUG FIXES:
  - Break if mist fails to write to a web socket client; it's presumably dead, and
  we don't want mist looping forever trying to read/write to something git will never
  be able to.

## Previous (March 9, 2016)

This change log began with version 1.0.0. Any prior changes can be seen by viewing
the commit history.
