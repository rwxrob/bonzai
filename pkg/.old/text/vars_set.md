The {{aka}} command writes the changes to the specified cached variable in a way that is reasonably safe for system-wide concurrent writes by checking the file for any changes since last right and refusing to overwrite if so (much like editing from a Vim session). If no name is passed will throw an error. If no new value arguments are passed will behave as if {{cmd "get"}} was called instead.

The exact process is as follows:

1. Save the current time in nanoseconds
2. Load and parse {{ execachefile "vars" }} into vars.Map
3. Change the specified value
4. Check file for changes since saved time, error if changed
5. Marshal vars.Map and atomically write to file

***Multiple argument fields are joined with space.*** In UNIX tradition, multiple arguments are assumed to be a part of a single string argument to be joined with spaces. This saves users from having to quote everything when it is not needed.

After setting the value, the new value is printed as if the {{cmd "get"}} was called.
