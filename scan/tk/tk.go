/*
Package tk contains tokens that when detected in the Expect expression stream cause specific behavior or change the state of the scan.R.
*/
package tk

const (

	// EOD is a special value that is returned when the end of data is
	// reached enabling functional parser functions to look for it reliably
	// no matter what is being parsed. Since rune is alias for int32 and
	// Unicode (currently) ends at \U+FFFD we are safe to use the largest
	// possible valid rune value. Passing EOD directly to Expect always
	// stops the scan where it is.
	EOD rune = 1<<31 - (iota + 1) // max int32

	// ANY represents any possible single rune similar to question mark (?)
	// when globing and dot (.) in most regular expression syntaxes.
	ANY
)
