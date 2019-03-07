package incapsula

import (
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func suppressEquivalentStringDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldSlice := strings.Split(old, ",")
	newSlice := strings.Split(new, ",")
	sort.Strings(oldSlice)
	sort.Strings(newSlice)

	return reflect.DeepEqual(oldSlice, newSlice)
}
