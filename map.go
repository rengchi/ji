package ji

// MergeMap 递归合并多个 map[string]any
func MergeMap(mObj ...map[string]any) map[string]any {
	// 预估容量，减少 map 扩容次数，使用 totalSize 来估算容量
	newObj := make(map[string]any, 16)

	for _, m := range mObj {
		for k, v := range m {
			if existingVal, exists := newObj[k]; exists {
				// 处理 key 冲突情况
				newObj[k] = mergeValues(existingVal, v)
			} else {
				newObj[k] = v
			}
		}
	}

	return newObj
}

// mergeValues 递归合并值
func mergeValues(oldVal, newVal any) any {
	switch oldValTyped := oldVal.(type) {
	case map[string]any:
		// 如果新值也是 map，递归合并
		if newValTyped, ok := newVal.(map[string]any); ok {
			return MergeMap(oldValTyped, newValTyped)
		}
		// 如果类型不匹配，直接替换
		return newVal

	case []any:
		// 如果新值是 slice，合并数组
		if newValTyped, ok := newVal.([]any); ok {
			return append(oldValTyped, newValTyped...)
		}
		// 如果类型不匹配，直接替换
		return newVal
	}

	// 对于其他类型，直接覆盖原值
	return newVal
}
