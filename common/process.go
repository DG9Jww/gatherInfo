package common

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func StringToSlice(target string, sep string) []string {
	t := strings.TrimSpace(target)
	return strings.Split(t, sep)
}

//covert text row by row to string string
func FileToSlice(f *os.File) []string {
	scanner := bufio.NewScanner(f)
	var s []string
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	return s
}

func MatchStr(substr string, str string) bool {
	ok, _ := regexp.MatchString(substr, str)
	if ok {
		return true
	} else {
		return false
	}
}

func MatchInt(i int, list []int) bool {
	for _, v := range list {
		if i == v {
			return true
		}
	}
	return false
}

//iterate the slice and check whether the slice item equal "str"
func IsStringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

//iterate the slice and check whether item within str
func IsSliceWithinStr(str string, list []string) bool {
	for _, item := range list {
		if strings.Contains(str, item) {
			return true
		}
	}
	return false
}

func PortToList(s string) []int {
	if MatchStr("-", s) {
		list := strings.Split(s, "-")
		min, err := strconv.Atoi(list[0])
		if err != nil {
			return nil
		}
		max, err := strconv.Atoi(list[1])
		if err != nil {
			return nil
		}

		var ports []int
		for min <= max {
			ports = append(ports, min)
			min++
		}
		return ports
	} else {
		list := strings.Split(s, ",")

		var ports []int
		for _, port := range list {
			p, err := strconv.Atoi(port)
			if err != nil {
				return nil
			}
			ports = append(ports, p)
		}
		return ports
	}
}

func GetMinInt(list []int) int {
	var min int
	if len(list) != 0 {
		min = list[0]
		for _, v := range list {
			if v < min {
				min = v
			}
		}
	}
	return min
}

func GetMaxInt(list []int) int {
	var max int
	if len(list) != 0 {
		max = list[0]
		for _, v := range list {
			if v > max {
				max = v
			}
		}
		return max
	}
	return max
}

func RandomInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

//delete item according index
func DeleteStringFromSlice(strSlice []string, index int) []string {
	tmp1 := strSlice[:index]
	tmp2 := strSlice[index+1:]
	tmp3 := append(tmp1, tmp2...)
	return tmp3
}

//process ip and domain according to regular expression
func ProRegularExp(tmpResSlice *[]string, exp string) []string {
	var tmp []string
	re, err := regexp.Compile(exp)
	if err != nil {
		return nil
	}

	for _, res := range *tmpResSlice {
		s := re.FindAllString(res, -1)
		for _, i := range s {
			tmp = append(tmp, i)
		}
	}
	return tmp
}


//get subdomain regular expression
func GetExp(field string) string {
	return fmt.Sprintf(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z]{0,62})*\.(%s)$?`, field)
}
