// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package bonzai provides a rooted node tree of commands and singular
parameters making tab completion a breeze and complicated applications
much easier to intuit without reading all the docs. Documentation is
embedded with each command removing the need for separate man pages and
such and can be viewed as text or a locally served web page.

Rooted Node Tree

Commands and parameters are linked to create a rooted node tree of the
following types of nodes:

    * Leaves with a method and optional parameters
		* Branches with leaves, other branches, and a optional method
		* Parameters, single words that are passed to a leaf command

*/
package bonzai
