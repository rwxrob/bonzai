package tk

// EOD is a special value that is returned when the end of data is
// reached enabling functional parser functions to look for it reliably
// no matter what is being parsed.
const EOD = 1<<31 - 1 // max int32
