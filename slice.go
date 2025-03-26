// 参数 https://github.com/AlistGo/alist/blob/main/pkg/utils/slice.go

package ji

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// SliceEqual 检查两个切片是否相等
// 如果两个切片的长度不同或内容不一致，返回 false
// 否则返回 true
func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// SliceContains 检查切片是否包含指定元素
// 如果切片中包含元素 v，返回 true；否则返回 false
func SliceContains[T comparable](arr []T, v T) bool {
	return slices.Contains(arr, v)
}

// SliceAllContains 检查切片是否完全包含指定元素集合
// 如果切片中包含所有指定元素，返回 true；否则返回 false
func SliceAllContains[T comparable](arr []T, vs ...T) bool {
	vsMap := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		vsMap[v] = struct{}{}
	}
	for _, v := range vs {
		if _, ok := vsMap[v]; !ok {
			return false
		}
	}
	return true
}

// SliceConvert 将一个类型的切片转换为另一个类型的切片
// 根据给定的转换函数进行转换，转换过程中遇到错误时返回错误
func SliceConvert[S any, D any](srcS []S, convert func(src S) (D, error)) ([]D, error) {
	res := make([]D, len(srcS))
	for i, src := range srcS {
		dst, err := convert(src)
		if err != nil {
			return nil, err
		}
		res[i] = dst
	}
	return res, nil
}

// MustSliceConvert 将一个类型的切片转换为另一个类型的切片，忽略错误
// 如果转换过程中出错，直接忽略错误进行转换
func MustSliceConvert[S any, D any](srcS []S, convert func(src S) D) []D {
	res := make([]D, len(srcS))
	for i, src := range srcS {
		res[i] = convert(src)
	}
	return res
}

// MergeErrors 合并多个错误信息为一个错误
// 将所有错误的字符串信息拼接起来，并返回一个新的错误
func MergeErrors(_errs ...error) error {
	errStr := strings.Join(MustSliceConvert(_errs, func(err error) string {
		return err.Error()
	}), "\n")
	if errStr != "" {
		return Miss(errStr)
	}
	return nil
}

// SliceMeet 检查切片中是否存在满足条件的元素
// 根据给定的条件函数，如果切片中存在满足条件的元素，返回 true
// 否则返回 false
func SliceMeet[T1, T2 any](arr []T1, v T2, meet func(item T1, v T2) bool) bool {
	for _, item := range arr {
		if meet(item, v) {
			return true
		}
	}
	return false
}

// SliceFilter 根据指定条件过滤切片中的元素
// 返回一个新的切片，包含所有满足条件的元素
func SliceFilter[T any](arr []T, filter func(src T) bool) []T {
	res := make([]T, 0, len(arr))
	for _, src := range arr {
		if filter(src) {
			res = append(res, src)
		}
	}
	return res
}

// SliceReplace 用指定函数替换切片中的所有元素
// 根据给定的替换函数，对切片中的每个元素进行替换
func SliceReplace[T any](arr []T, replace func(src T) T) {
	for i, src := range arr {
		arr[i] = replace(src)
	}
}

// SplitStringToIntSlice 将逗号分隔的字符串拆分成整数切片
// 将传入的逗号分隔字符串拆分为整数切片，返回切片和可能的错误
func SplitStringToIntSlice(record string) ([]int, error) {
	strSlice := strings.Split(record, ",")
	intSlice := make([]int, len(strSlice))
	for i, s := range strSlice {
		num, err := strconv.Atoi(s)
		// 如果遇到空字符串，返回错误
		if s == "" {
			return nil, fmt.Errorf("输入的字符串在索引 %d 处为空", i)
		}
		if err != nil {
			return nil, err
		}
		intSlice[i] = num
	}
	return intSlice, nil
}

// SplitStringToIntSliceIgnoringEmpty 将逗号分隔的字符串拆分成整数切片，跳过空字符串
// 将传入的逗号分隔字符串拆分为整数切片，跳过空字符串并返回切片和可能的错误
func SplitStringToIntSliceIgnoringEmpty(record string) ([]int, error) {
	strSlice := strings.Split(record, ",")
	intSlice := make([]int, 0, len(strSlice))
	for _, s := range strSlice {
		// 跳过空字符串
		if s == "" {
			continue
		}
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num)
	}
	return intSlice, nil
}

// SplitStringToUintSlice 将逗号分隔的字符串拆分成 uint 类型的整数切片
// 将传入的逗号分隔字符串拆分为无符号整数切片，返回切片和可能的错误
func SplitStringToUintSlice(record string) ([]uint, error) {
	strSlice := strings.Split(record, ",")
	uintSlice := make([]uint, len(strSlice))
	for i, s := range strSlice {
		// 如果遇到空字符串，返回错误
		if s == "" {
			return nil, fmt.Errorf("输入的字符串在索引 %d 处为空", i)
		}
		num, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}
		uintSlice[i] = uint(num)
	}
	return uintSlice, nil
}

// SplitStringToUintSliceIgnoringEmpty 将逗号分隔的字符串拆分成 uint 类型的整数切片，跳过空字符串
// 将传入的逗号分隔字符串拆分为无符号整数切片，跳过空字符串并返回切片和可能的错误
func SplitStringToUintSliceIgnoringEmpty(record string) ([]uint, error) {
	strSlice := strings.Split(record, ",")
	uintSlice := make([]uint, 0, len(strSlice))
	for _, s := range strSlice {
		// 跳过空字符串
		if s == "" {
			continue
		}
		num, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}
		uintSlice = append(uintSlice, uint(num))
	}
	return uintSlice, nil
}

// UintSliceToString 将整数切片转换为用逗号连接的字符串，并去重
// 将传入的整数切片转换为字符串，并使用逗号连接返回
// 同时会去重重复的元素
func UintSliceToString(intSlice []uint) string {
	// 使用 map 去重
	idSet := make(map[uint]struct{}, len(intSlice))
	var uniqueIntSlice []uint
	for _, id := range intSlice {
		if _, exists := idSet[id]; !exists {
			idSet[id] = struct{}{}
			uniqueIntSlice = append(uniqueIntSlice, id)
		}
	}

	// 转换为字符串切片
	strSlice := make([]string, len(uniqueIntSlice))
	for i, v := range uniqueIntSlice {
		strSlice[i] = fmt.Sprint(v)
	}

	// 返回逗号连接的字符串
	return strings.Join(strSlice, ",")
}

// MergeSlices 合并两个切片，不去重
// 将两个切片合并成一个新的切片，返回合并后的切片
func MergeSlices[T any](a, b []T) []T {
	return append(a, b...)
}

// UniqueSlice 去重操作
// 对传入的切片进行去重，返回一个去重后的切片
func UniqueSlice[T comparable](arr []T) []T {
	seen := make(map[T]struct{}, len(arr))
	result := make([]T, 0, len(arr))
	for _, v := range arr {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
