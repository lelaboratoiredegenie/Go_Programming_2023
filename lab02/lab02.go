package main

import "fmt"
import "strconv"

func main() {
	var n int64

	fmt.Print("Enter a number: ")
	fmt.Scanln(&n)

	result := Sum(n)
	fmt.Println(result)
}

func Sum(n64 int64) string {
	// TODO: Finish this function
	ans := 0
	ans_str := ""
	n := int(n64)
	for i := 1; i < n-1; i++ {
		if i % 7 != 0{
			ans += i
			strNum := strconv.Itoa(i)
			ans_str += strNum
			ans_str += "+"
		}
	}

	if (n-1) % 7 != 0 && n % 7 != 0 {
		ans += n-1
		ans += n
		strNum := strconv.Itoa(n-1)
		ans_str += strNum
		ans_str += "+"
		strNum = strconv.Itoa(n)
		ans_str += strNum
	} else if (n-1) % 7 != 0 {
		ans += n-1
		strNum := strconv.Itoa(n-1)
		ans_str += strNum
	} else{
		ans += n
		strNum := strconv.Itoa(n)
		ans_str += strNum
	}
	
	strNum := strconv.Itoa(ans)
	ans_str += "="
	ans_str += strNum

	return ans_str
}
