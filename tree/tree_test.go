package tree_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rwxrob/bonzai/tree"
)

func ExampleTree() {
	t := tree.New([]string{"foo", "bar"})
	t.Print()
	// Output:
	// {"Trunk":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleTree_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	t := tree.New([]string{"foo", "bar"})
	t.Log()
	// Output:
	// {"Trunk":[1],"Types":["UNKNOWN","foo","bar"]}
}

func ExampleTree_Seed() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	t := tree.New([]string{"foo", "bar"})
	s := t.Seed(2, "")
	fmt.Println(s.IsRoot(), s.IsDetached(), s.IsNull())
	// Output:
	// true true true
}

func ExampleNode_Log() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Log()
	// Output:
	// [1]
}

func ExampleNode_Info() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Info()

	// Type:       1
	// Value:      ""
	// IsRoot:     true
	// IsDetached: true
	// IsLeaf:     false
	// IsBranch:   false
	// IsNull:     true
}

func ExampleNode_leaf() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	n.V = "something"
	n.Print()
	n.T = 0
	n.V = ""
	n.Print()
	// Output:
	// [1]
	// [1,"something"]
	// []
}

func ExampleNode_branch() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.NewUnder(2, "something")
	n.Print()
	n.PrettyPrint()
	// Output:
	// [1,[[2,"something"]]]
	// ["foo", [
	//   ["bar", "something"]
	// ]]
}

func ExampleNode_NewRight() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewRight(2, "some")
	t.Print()
	n.Right().Print()
	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
}

func ExampleNode_NewLeft() {

	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewLeft(2, "some")
	t.Print()
	n.Left().Print()

	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
}

func ExampleNode_NewUnder() {
	n := tree.New([]string{"foo", "bar"}).Trunk
	n.Print()
	t := n.NewUnder(2, "some")
	t.Print()
	n.FirstUnder().Print()
	n.LastUnder().Print()
	n.Print()

	//Output:
	// [1]
	// [2,"some"]
	// [2,"some"]
	// [2,"some"]
	// [1,[[2,"some"]]]
}

func ExampleNode_Graft() {
	t := tree.New([]string{"foo", "😭", "💔"})
	n := t.Trunk
	u := n.NewUnder(2, "graftme")
	u.NewUnder(3, "under1") // should go with orig
	u.NewUnder(3, "under2") // should go with orig
	u.NewLeft(3, "left")
	u.NewRight(3, "right")
	x := u.Graft(t.Seed(3, "newbie"))
	fmt.Println(u.AllUnder()) // still has leaves and branches
	fmt.Println(x.AllUnder()) // never had them
	fmt.Println(x.Left())     // new left
	fmt.Println(x.Right())    // new right
	fmt.Println(u.Left())     // detached
	fmt.Println(u.Right())    // detached
	fmt.Println(x.Branch())

	//Output:
	// [[3,"under1"] [3,"under2"]]
	// []
	// [3,"left"]
	// [3,"right"]
	// <nil>
	// <nil>
	// [1,[[3,"left"],[3,"newbie"],[3,"right"]]]

}

func ExampleNode_GraftRight() {
	t := tree.New([]string{"foo", "😭", "💔"})
	n := t.Trunk
	u := n.NewUnder(3, "under")
	x := u.GraftRight(t.Seed(3, "newbie"))
	n.Print()
	u.Print()
	x.Left().Print()

	//Output:
	// [1,[[3,"under"],[3,"newbie"]]]
	// [3,"under"]
	// [3,"under"]

}

func ExampleNode_GraftLeft() {
	t := tree.New([]string{"foo", "😭", "💔"})
	n := t.Trunk
	u := n.NewUnder(3, "under")
	x := u.GraftLeft(t.Seed(3, "newbie"))
	n.Print()
	u.Print()
	x.Right().Print()

	//Output:
	// [1,[[3,"newbie"],[3,"under"]]]
	// [3,"under"]
	// [3,"under"]

}

func ExampleNode_GraftUnder() {
	t := tree.New([]string{"foo", "😭", "💔"})
	n := t.Trunk
	u := n.NewUnder(3, "under")
	x := u.GraftUnder(t.Seed(3, "newbie"))
	n.Print()
	u.Print()
	x.Branch().Print()

	//Output:
	// [1,[[3,[[3,"newbie"]]]]]
	// [3,[[3,"newbie"]]]
	// [3,[[3,"newbie"]]]

}

func ExampleNode_Prune() {
	t := tree.New([]string{"foo", "😭", "💔"})
	n := t.Trunk
	u := n.NewUnder(2, "pruneme")
	u1 := u.NewUnder(3, "under") // should go with orig
	left := u.NewLeft(3, "left")
	right := u.NewRight(3, "right")
	x := u.Prune()
	x.Print()            // should take stuff under with it
	u1.Branch().Print()  // same as previous
	left.Right().Print() // should be joined to the left now
	right.Left().Print() // should be joined to right now

	//Output:
	// [2,[[3,"under"]]]
	// [2,[[3,"under"]]]
	// [3,"right"]
	// [3,"left"]

}

func ExampleNode_Take() {
	n := tree.New([]string{"foo", "😭", "💔"}).Trunk
	u1 := n.NewUnder(2)
	u2 := n.NewUnder(3)
	u2.NewUnder(3)
	u2.NewUnder(3)
	u1.Print()
	u2.Print()
	u1.Take(u2)
	u1.Print()
	u2.Print()

	// Output:
	// [2]
	// [3,[[3],[3]]]
	// [2,[[3],[3]]]
	// [3]
}

func ExampleNode_Visit() {

	n := tree.New([]string{"🌴", "😭", "💔", "💀"}).Trunk
	l2 := n.NewUnder(2, "two")
	l_2 := l2.NewUnder(2, "_two")
	l3 := n.NewUnder(3, "three")
	l_3 := l3.NewUnder(3, "_three")
	l4 := n.NewUnder(4, "four")
	l_4 := l4.NewUnder(4, "_four")

	n.Print()

	l2.Print()
	l_2.Print()
	l3.Print()
	l_3.Print()
	l4.Print()
	l_4.Print()

	n.Visit(
		func(n *tree.Node) any {
			n.Print()
			return nil
		}, nil)

	// Output:
	// [1,[[2,[[2,"_two"]]],[3,[[3,"_three"]]],[4,[[4,"_four"]]]]]
	// [2,[[2,"_two"]]]
	// [2,"_two"]
	// [3,[[3,"_three"]]]
	// [3,"_three"]
	// [4,[[4,"_four"]]]
	// [4,"_four"]
	// [1,[[2,[[2,"_two"]]],[3,[[3,"_three"]]],[4,[[4,"_four"]]]]]
	// [2,[[2,"_two"]]]
	// [2,"_two"]
	// [3,[[3,"_three"]]]
	// [3,"_three"]
	// [4,[[4,"_four"]]]
	// [4,"_four"]

}

func ExampleNode_VisitAsync() {

	n := tree.New([]string{"🌴", "😭", "💔", "💀"}).Trunk
	n.NewUnder(2, "two").NewUnder(2, "_two")
	n.NewUnder(3, "three").NewUnder(3, "_three")
	n.NewUnder(4, "four").NewUnder(4, "_four")

	n.VisitAsync(
		func(n *tree.Node) any {
			n.Print()
			time.Sleep(30 * time.Millisecond)
			return nil
		}, 3, nil)

	// note that "unordered" output is required

	// Unordered Output:
	// [1,[[2,[[2,"_two"]]],[3,[[3,"_three"]]],[4,[[4,"_four"]]]]]
	// [2,[[2,"_two"]]]
	// [2,"_two"]
	// [3,[[3,"_three"]]]
	// [3,"_three"]
	// [4,[[4,"_four"]]]
	// [4,"_four"]

}

func ExampleNode_AllUnder() {
	n := tree.New([]string{"🌴", "😭", "💔", "💀"}).Trunk
	n.NewUnder(2)
	n.NewUnder(3)
	n.NewUnder(4)
	n.Print()
	// Output:
	// [1,[[2],[3],[4]]]
}

func ExampleNode_JSON() {
	n := tree.New([]string{"🌴", "😭", "💔", "💀"}).Trunk
	n.NewUnder(2)
	n.NewUnder(3)
	n.NewUnder(4)
	n.NewUnder(5)
	n.Print()
	// Output:
	// [1,[[2],[3],[4],[5]]]
}
