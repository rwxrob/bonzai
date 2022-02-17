/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
