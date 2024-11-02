# Easy Input-Output Test Specifications in YAML

Often testing requires a large number of inputs with specific outputs --- particularly when testing parsers and grammar implementations. Such tests are also usually required across different language implementations with different testing suites and approaches. For testing such specifications it makes sense to maintain inputs and outputs in a language-agnostic way. This is one such approach.

YAML is used as the format for the test file because of its strong support for multi-line data. Specifications --- which are maintained by human hands --- are a particularly good use case for YAML structured data.

Any number or fields may be included, but all key names must begin with an initial capital letter. The following are common and expected:

* `Name` - name or title of the spec
* `Version` - semantic version of the spec being tested
* `Source` - URL of YAML files that define the spec
* `Issues` - URL to report issues with the spec
* `Discuss` - URL of place to discuss the spec
* `Notes` - any notes about the specification
* `Date` - date of the specification
* `License` - license of the spec

One of the following is required although both are allowed:

* `Tests` - array of tests
* `Sections` - array of sections containing tests

If both are included the non-grouped tests will be checked first.

The content of `Tests` can be a simple array of tests, or an array of
test groups. All tests must have an `I` and `O` property and can
optionally add an `N` property as well:

Must:

* `I`: Any string, number, or boolean containing the input. 
* `O`: Any string, number, or boolean containing the required output. 

Optional:

* `N`: Short note about what is being tests, often from the specification.

Here's an example from the CommonMark specification. The first uses the array format:

```yaml
Name: CommonMark
Version: v1.0.0
Source: https://gitlab.com/commonmark/commonmark-spec
Issues: https://gitlab.com/commonmark/commonmark-spec/issues
Discuss: https://talk.commonmark.org/
Notes: CommonMark is a clarification version of Markdown.
Tests:
- I: '\tfoo\tbaz\t\tbim\n'
  O: '<pre><code>foo\tbaz\t\tbim\n</code></pre>\n'
  N: 'tabs are equal to four spaces' 
- I: '  \tfoo\tbaz\t\tbim\n'
  O: '<pre><code>foo\tbaz\t\tbim\n</code></pre>\n'
```

And the same as a map with keys matching the sections:

```yaml
Sections:
- Name: Tabs
  Notes: Tests related to use of the tab character.
  Tests:
  - I: '\tfoo\tbaz\t\tbim\n'
    O: '<pre><code>foo\tbaz\t\tbim\n</code></pre>\n'
    N: 'tabs are equal to four spaces' 
  - I: '  \tfoo\tbaz\t\tbim\n'
    O: '<pre><code>foo\tbaz\t\tbim\n</code></pre>\n'
```

See [tests](tinout_test.go) for usage examples.

## Go (golang) Reference Implementation


A reference implementation in Go has been included for reference and
import inclusion:

```go
package mytest

import (
  "testing"

  "github.com/rwxrob/tinout"
)

// ...
```
