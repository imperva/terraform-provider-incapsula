package incapsula

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}
func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}
