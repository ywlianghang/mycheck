package Algorithm

import "fmt"

func BinaryFind(arr *[6]int, leftIndex int,rightIndex int,findVal int){
	//如果 leftIndex > rightIndex,说明是递归的退出条件
	if leftIndex > rightIndex {
		fmt.Println("找不到了...")
		return
	}
	// 先找到中间的下标
	middleIndex := (leftIndex + rightIndex) / 2
	if (*arr)[middleIndex] > findVal {
		rightIndex = middleIndex - 1
		BinaryFind(arr, leftIndex, rightIndex, findVal)
	} else if (*arr)[middleIndex] < findVal {
		leftIndex = middleIndex + 1
		BinaryFind(arr, leftIndex, rightIndex, findVal)
	} else {
		fmt.Printf("找到了，值为%v, 下标为%v...", findVal,middleIndex)
	}
}
