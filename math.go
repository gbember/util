// math.go
package util

func Abs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}
