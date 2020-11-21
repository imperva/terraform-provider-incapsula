package incapsula

import (
	"encoding/json"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentStringDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldSlice := strings.Split(old, ",")
	newSlice := strings.Split(new, ",")
	sort.Strings(oldSlice)
	sort.Strings(newSlice)

	return reflect.DeepEqual(oldSlice, newSlice)
}

func suppressEquivalentJSONStringDiffs(k, old, new string, d *schema.ResourceData) bool {
	var o1 interface{}
	var o2 interface{}

	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)

	if old == "" && new == "" {
		return true
	}
	if old == "" && new != "" {
		return false
	}
	if old != "" && new == "" {
		return false
	}

	var err error
	err = json.Unmarshal([]byte(old), &o1)
	if err != nil {
		log.Panicf("Invalid JSON (current value): %s", err.Error())
		return false
	}
	err = json.Unmarshal([]byte(new), &o2)
	if err != nil {
		log.Panicf("Invalid JSON (new value): %s", err.Error())
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
