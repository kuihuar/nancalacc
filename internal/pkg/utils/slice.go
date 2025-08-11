package utils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesBy[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	seen := make(map[K]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		key := keyFunc(item)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// 使用示例：
// nums := []int{1, 2, 3, 4, 5}
// fmt.Println(Contains(nums, 3)) // true
// fmt.Println(Contains(nums, 6)) // false
