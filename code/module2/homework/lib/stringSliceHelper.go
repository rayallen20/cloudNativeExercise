package lib

func Slice2Str(slice []string) (res string) {
	for k, v := range slice {
		res += v
		if k != len(slice)-1 {
			res += ", "
		}
	}
	
	return res
}
