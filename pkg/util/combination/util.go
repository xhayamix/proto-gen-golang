package combination

type Item interface {
	GetValue() int32
	GetWeight() int32
}

type Items []Item

func Search(items Items, capacity, minValue, maxValue int32) []Items {
	combinations := findCombinations(items, capacity, minValue, maxValue)

	result := make([]Items, 0, len(items)*len(items))
	for _, combination := range combinations {
		its := make(Items, 0, len(combination))
		for _, idx := range combination {
			its = append(its, items[idx])
		}
		result = append(result, its)
	}

	return result
}

func findCombinations(items Items, capacity, minValue, maxValue int32) [][]int32 {
	result := make([][]int32, 0, len(items))
	var find func(i int32, combination []int32)
	find = func(i int32, combination []int32) {
		if i == 0 {
			var totalWeight int32
			var totalValue int32
			for _, itemIndex := range combination {
				totalWeight += items[itemIndex].GetWeight()
				totalValue += items[itemIndex].GetValue()
			}
			if totalWeight <= capacity && totalValue >= minValue && totalValue <= maxValue {
				result = append(result, combination)
			}
			return
		}

		newCombination := make([]int32, len(combination)+1)
		copy(newCombination, combination)
		newCombination[len(combination)] = i - 1
		find(i-1, newCombination)
		find(i-1, combination)
	}

	find(int32(len(items)), []int32{})

	return result
}
