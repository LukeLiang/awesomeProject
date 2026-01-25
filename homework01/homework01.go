package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// SingleNumber (只出现一次的数字)
func SingleNumber(nums *[]uint) uint {

	if nums == nil || len(*nums) == 0 {
		fmt.Printf("数组为空...\n")
		return 0
	}

	givenMap := make(map[uint]uint)
	for _, value := range *nums {
		v1, ok := givenMap[value]
		if ok {
			givenMap[value] = v1 + 1
		} else {
			givenMap[value] = 1
		}
	}

	for key, value := range givenMap {
		if value == 1 {
			fmt.Printf("find it, the key is %d\n", key)
			return key
		}
	}

	fmt.Printf("没有找到...\n")
	return 0

	//异或的特性：
	//- a ^ a = 0（相同的数异或为0）
	//- a ^ 0 = a（任何数与0异或为它本身）
	//if nums == nil || len(*nums) == 0 {
	//	return 0
	//}
	//
	//var result uint = 0
	//for _, value := range *nums {
	//	result ^= value
	//}
	//return result
}

// 回文数
func isPalindromeString(x int) bool {
	// 负数不可能是回文数
	if x < 0 {
		return false
	}

	// 将整数转换为字符串
	s := strconv.Itoa(x)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}

	return true
}

// IsValidParentheses 有效的括号
func IsValidParentheses(kh string) bool {
	length := len(kh)
	if strconv.Itoa(length%2) != "0" {
		fmt.Printf("奇数 invalid")
		return false
	}
	for i := 0; i < len(kh)-1; i++ {
		if strconv.Itoa(i%2) == "0" {
			if kh[i] != '(' && kh[i] != '{' && kh[i] != '[' {
				fmt.Printf("invalid")
				return false
			}

			switch kh[i] {
			case '(':
				if ')' != kh[i+1] {
					fmt.Printf("invalid")
					return false
				}
			case '{':
				if '}' != kh[i+1] {
					fmt.Printf("invalid")
					return false
				}
			case '[':
				if ']' != kh[i+1] {
					fmt.Printf("invalid")
					return false
				}
			default:
				fmt.Printf("invalid")
				return false
			}
		}

		if strconv.Itoa(i%2) != "0" {
			if kh[i] != ')' && kh[i] != '}' && kh[i] != ']' {
				fmt.Printf("invalid")
				return false
			}
		}

	}
	return true
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	// 1. 特殊情况处理：如果切片为空，返回空字符串
	if len(strs) == 0 {
		return ""
	}

	// 2. 以第一个字符串作为基准，逐个字符向后遍历
	firstStr := strs[0]
	for i := 0; i < len(firstStr); i++ {
		// 取出基准字符串第 i 个位置的字符
		char := firstStr[i]

		// 3. 遍历剩余的字符串，检查第 i 个位置
		for j := 1; j < len(strs); j++ {
			// 判定条件：
			// - 如果当前字符串长度已经到了 i (说明它比基准串短)
			// - 或者当前字符串第 i 个位置的字符与基准串不同
			if i == len(strs[j]) || strs[j][i] != char {
				// 只要发现不匹配，立即截取并返回基准串的前 i 个字符
				return firstStr[:i]
			}
		}
	}

	// 4. 如果循环正常结束，说明基准串本身就是最短的公共前缀
	return firstStr
}

// 加一
func plusOne(digits *[]int) {
	digitStr := strings.Builder{}
	for _, value := range *digits {
		digitStr.WriteString(strconv.Itoa(value))
	}
	atoi, _ := strconv.Atoi(digitStr.String())
	itoa := strconv.Itoa(atoi + 1)
	fmt.Printf(itoa)
	//for index, _ := range itoa {
	//	fmt.Printf("%c", itoa[index])
	//}
}

// 删除有序数组中的重复项
func removeDuplicates(nums *[]int) int {
	if len(*nums) == 0 {
		return 0
	}

	slow := 0

	for fast := 1; fast < len(*nums); fast++ {
		if (*nums)[fast] != (*nums)[slow] {
			slow++
			(*nums)[slow] = (*nums)[fast]
		}
	}

	// 截断切片，只保留去重后的元素
	*nums = (*nums)[:slow+1]
	return slow + 1
}

// 合并区间
func mergeIntervals(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	// 按区间起始位置排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 存储合并后的区间
	result := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		last := result[len(result)-1]
		current := intervals[i]

		// 如果当前区间与最后一个区间重叠
		if current[0] <= last[1] {
			// 合并区间，更新结束位置为两者中较大的
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			// 没有重叠，直接添加
			result = append(result, current)
		}
	}

	return result
}

// 两数之和
func findSumTarget(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-1; j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func main() {

	// 1. 只出现一次的数字
	nonInt := []uint{1, 2, 3, 4, 3, 2, 1}
	SingleNumber(&nonInt)

	// 2. 回文数
	palindrome := 12321
	fmt.Printf("数字%d是回文数？%t\n", palindrome, isPalindromeString(palindrome))

	// 3. 有效括号
	kh := "(){}[]"
	fmt.Printf("%s是有效括号？%t\n", kh, IsValidParentheses(kh))

	// 4. 最长公共前缀
	strs := []string{"flower", "flow", "flight"}
	fmt.Printf("最长公共前缀是%s\n", longestCommonPrefix(strs))

	// 5. 加一
	digits := []int{1, 2, 3}
	plusOne(&digits)
	for index := range digits {
		fmt.Printf("%c", digits[index])
	}
	fmt.Println()

	// 6. 删除有序数组中的重复项
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	nums2 := removeDuplicates(&nums)
	fmt.Printf("%d\n", nums2)
	fmt.Println("去重之后数组的值", nums)

	// 7. 合并区间
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Println("合并前:", intervals)
	fmt.Println("合并后:", mergeIntervals(intervals))

	// 8. 两数之和
	nums = []int{2, 7, 11, 15}
	target := 9
	if result := findSumTarget(nums, target); result != nil {
		fmt.Println(result)
	} else {
		fmt.Printf("没有找到合适的值\n")
	}

}
