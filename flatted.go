package flatted

import (
	"strconv"
)

var index_data map[string]string
var data_index map[string]string
var indexGlobal int // 全局index

func Flatted(data string) string {
	index_data = make(map[string]string)
	data_index = make(map[string]string)
	indexGlobal = 0
	if data[0] == '[' {
		manualList(&data, getIndex(&data))
	} else {
		manualStruct(&data, getIndex(&data))
	}
	ans := ""
	for i := 0; i < len(index_data); i++ {
		ans += index_data[strconv.Itoa(i)]
		if i < len(index_data)-1 {
			ans += ","
		}
	}
	ans = "[" + ans + "]"
	// file, err := os.OpenFile("./test_go.txt", os.O_WRONLY, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = file.WriteString(ans)
	// if err != nil {
	// 	panic(err)
	// }
	return ans
}

func findMatchRightBrace(input *string, fr int, en int, left byte, right byte) int {
	cnt := 1
	flag := false
	for fr < en {
		// flag判断当前是否在字符串内，若在字符串内忽略匹配
		if (*input)[fr] == '"' && (*input)[fr-1] != '\\' {
			flag = !flag
		}
		if flag {
			fr++
			continue
		}
		if (*input)[fr] == left {
			cnt++
		} else if (*input)[fr] == right {
			cnt--
		}
		if cnt == 0 {
			return fr
		}
		fr++
	}
	return -1
}

func findNextElement(input *string, fr int, en int, ch byte) int {
	for fr < en {
		// 考虑转义
		if (*input)[fr] == ch && (*input)[fr-1] != '\\' {
			return fr
		}
		fr++
	}
	return -1
}

func getIndex(data *string) string {
	// if ind, ok := data_index[*data]; ok {
	// 	return ind
	// }
	temp := strconv.Itoa(indexGlobal)
	indexGlobal++
	data_index[*data] = temp
	return temp
}

// 处理列表。传入的index是当前处理data的index
func manualList(data *string, index string) {
	// 使用tempAns拼出当前处理的data
	tempAns := ""
	// 处理空list
	if *data == "[]" {
		index_data[index] = *data
		return
	}
	for i := 1; i < len(*data); {
		// 兼容未压缩的json
		for (*data)[i] == ' ' {
			i++
		}
		switch (*data)[i] {
		// 递归处理struct
		case '{':
			en := findMatchRightBrace(data, i+1, len(*data), '{', '}')
			tempData := (*data)[i : en+1]
			tempIndex := getIndex(&tempData)
			manualStruct(&tempData, tempIndex)
			tempAns += "\"" + tempIndex + "\""
			i = en + 1
		case ',':
			tempAns += ","
			i++
		// 递归处理list
		case '[':
			en := findMatchRightBrace(data, i+1, len(*data), '[', ']')
			tempData := (*data)[i : en+1]
			tempIndex := getIndex(&tempData)
			tempAns += "\"" + tempIndex + "\""
			manualList(&tempData, tempIndex)
			i = en + 1
		case ']':
			i++
		case '"':
			i++
			en := findNextElement(data, i, len(*data), '"')
			tempData := (*data)[i:en]
			tempIndex := getIndex(&tempData)
			index_data[tempIndex] = "\"" + tempData + "\""
			tempAns += "\"" + tempIndex + "\""
			i = en + 1
		default:
			en := findNextElement(data, i, len(*data), ',')
			if en == -1 {
				en = findNextElement(data, i, len(*data), ']')
			}
			tempAns += (*data)[i:en]
			i = en
		}
	}
	tempAns = "[" + tempAns + "]"
	index_data[index] = tempAns
}

// 处理结构体
func manualStruct(data *string, index string) {
	tempAns := ""
	if *data == "{}" {
		index_data[index] = *data
		return
	}
	for i := 1; i < len(*data); {
		for (*data)[i] == ' ' {
			i++
		}
		// 状态机
		switch (*data)[i] {
		// ":"后面是结构体内的data
		case ':':
			tempAns += ":"
			i++
			for (*data)[i] == ' ' {
				i++
			}
			switch (*data)[i] {
			case '"':
				i++
				en := findNextElement(data, i, len(*data), '"')
				tempData := (*data)[i:en]
				tempIndex := getIndex(&tempData)
				index_data[tempIndex] = "\"" + tempData + "\""
				tempAns += "\"" + tempIndex + "\""
				i = en + 1
			case '{':
				en := findMatchRightBrace(data, i+1, len(*data), '{', '}')
				tempData := (*data)[i : en+1]
				tempIndex := getIndex(&tempData)
				manualStruct(&tempData, tempIndex)
				tempAns += "\"" + tempIndex + "\""
				i = en + 1
			case '[':
				en := findMatchRightBrace(data, i+1, len(*data), '[', ']')
				tempData := (*data)[i : en+1]
				tempIndex := getIndex(&tempData)
				tempAns += "\"" + tempIndex + "\""
				manualList(&tempData, tempIndex)
				i = en + 1
			// 处理非字符串
			default:
				en := findNextElement(data, i, len(*data), ',')
				if en == -1 {
					en = findNextElement(data, i, len(*data), '}')
				}
				tempAns += (*data)[i:en]
				i = en
			}
		case ',':
			i++
			tempAns += ","
		// 结构体内的tag都是"\""开头
		case '"':
			i++
			en := findNextElement(data, i, len(*data), '"')
			tempAns += "\"" + (*data)[i:en] + "\""
			i = en + 1
		case '}':
			i++
		}
	}
	tempAns = "{" + tempAns + "}"
	index_data[index] = tempAns
}

var number_data map[int]string

func UnFlatted(data string) string {
	number_data = make(map[int]string)
	indexGlobal = 0
	getIndexDataMap(&data)
	return getAns(0)
}

func getIndexDataMap(input *string) {
	for i := 1; i < len(*input)-1; {
		switch (*input)[i] {
		case '{':
			en := findMatchRightBrace(input, i+1, len(*input), '{', '}')
			en++
			number_data[indexGlobal] = (*input)[i:en]
			indexGlobal++
			i = en
		case '[':
			en := findMatchRightBrace(input, i+1, len(*input), '[', ']')
			en++
			number_data[indexGlobal] = (*input)[i:en]
			indexGlobal++
			i = en
		case '"':
			en := findNextElement(input, i+1, len(*input), '"')
			en++
			number_data[indexGlobal] = (*input)[i:en]
			indexGlobal++
			i = en
		case ',':
			i++
		}
	}
}

func getAns(index int) string {
	ans := ""
	temp := number_data[index]
	if temp[0] == '"' && temp[len(temp)-1] == '"' {
		if _, err := strconv.Atoi(temp[1 : len(temp)-1]); err == nil {
			return temp
		}
	}
	for i := 0; i < len(temp); {
		switch temp[i] {
		case '"':
			if i == 0 || temp[i-1] != '\\' {
				en := findNextElement(&temp, i+1, len(temp), '"')
				maybeNumber := temp[i+1 : en]
				number, err := strconv.Atoi(maybeNumber)
				if err == nil {
					ans += getAns(number)
					i = en
				} else {
					ans += temp[i : en+1]
					i = en
				}
			}
		default:
			ans += string(temp[i])
		}
		i++
	}
	return ans
}
