package incapsula

import (
	"testing"
)

func TestSuppressEquivalentStringDiffsSameStringsInOrder(t *testing.T) {
	old := "a,b,c"
	new := "a,b,c"

	if !suppressEquivalentStringDiffs("", old, new, nil) {
		t.Errorf("Should be equivalent")
	}
}

func TestSuppressEquivalentStringDiffsSameStringsOutOfOrder(t *testing.T) {
	old := "a,b,c"
	new := "c,b,a"

	if !suppressEquivalentStringDiffs("", old, new, nil) {
		t.Errorf("Should be equivalent")
	}
}

func TestSuppressEquivalentStringDiffsSameStringsDifferent(t *testing.T) {
	old := "a,b,c"
	new := "x,y,z"

	if suppressEquivalentStringDiffs("", old, new, nil) {
		t.Errorf("Should not be equivalent")
	}
}
