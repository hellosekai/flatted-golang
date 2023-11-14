package flatted

import (
	"strconv"
)

type flatted struct {
	index_data  map[string]string
	data_index  map[string]string
	indexGlobal int // 全局index
	number_data map[int]string
}

func NewFlattedData() *flatted {
	return &flatted{
		index_data:  make(map[string]string),
		data_index:  make(map[string]string),
		indexGlobal: 0,
		number_data: make(map[int]string),
	}
}

func Flatted(data string) string {
	flat := NewFlattedData()

	if data[0] == '[' {
		flat.manualList(&data, flat.getIndex(&data))
	} else {
		flat.manualStruct(&data, flat.getIndex(&data))
	}
	ans := ""
	for i := 0; i < len(flat.index_data); i++ {
		ans += flat.index_data[strconv.Itoa(i)]
		if i < len(flat.index_data)-1 {
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

func (f *flatted) findMatchRightBrace(input *string, fr int, en int, left byte, right byte) int {
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

func (f *flatted) findNextElement(input *string, fr int, en int, ch byte) int {
	for fr < en {
		// 考虑转义
		if (*input)[fr] == ch && (*input)[fr-1] != '\\' {
			return fr
		}
		fr++
	}
	return -1
}

func (f *flatted) getIndex(data *string) string {
	// if ind, ok := data_index[*data]; ok {
	// 	return ind
	// }
	temp := strconv.Itoa(f.indexGlobal)
	f.indexGlobal++
	f.data_index[*data] = temp
	return temp
}

// 处理列表。传入的index是当前处理data的index
func (f *flatted) manualList(data *string, index string) {
	// 使用tempAns拼出当前处理的data
	tempAns := ""
	// 处理空list
	if *data == "[]" {
		f.index_data[index] = *data
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
			en := f.findMatchRightBrace(data, i+1, len(*data), '{', '}')
			tempData := (*data)[i : en+1]
			tempIndex := f.getIndex(&tempData)
			f.manualStruct(&tempData, tempIndex)
			tempAns += "\"" + tempIndex + "\""
			i = en + 1
		case ',':
			tempAns += ","
			i++
		// 递归处理list
		case '[':
			en := f.findMatchRightBrace(data, i+1, len(*data), '[', ']')
			tempData := (*data)[i : en+1]
			tempIndex := f.getIndex(&tempData)
			tempAns += "\"" + tempIndex + "\""
			f.manualList(&tempData, tempIndex)
			i = en + 1
		case ']':
			i++
		case '"':
			i++
			en := f.findNextElement(data, i, len(*data), '"')
			tempData := (*data)[i:en]
			tempIndex := f.getIndex(&tempData)
			f.index_data[tempIndex] = "\"" + tempData + "\""
			tempAns += "\"" + tempIndex + "\""
			i = en + 1
		default:
			en := f.findNextElement(data, i, len(*data), ',')
			if en == -1 {
				en = f.findNextElement(data, i, len(*data), ']')
			}
			tempAns += (*data)[i:en]
			i = en
		}
	}
	tempAns = "[" + tempAns + "]"
	f.index_data[index] = tempAns
}

// 处理结构体
func (f *flatted) manualStruct(data *string, index string) {
	tempAns := ""
	if *data == "{}" {
		f.index_data[index] = *data
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
				en := f.findNextElement(data, i, len(*data), '"')
				tempData := (*data)[i:en]
				tempIndex := f.getIndex(&tempData)
				f.index_data[tempIndex] = "\"" + tempData + "\""
				tempAns += "\"" + tempIndex + "\""
				i = en + 1
			case '{':
				en := f.findMatchRightBrace(data, i+1, len(*data), '{', '}')
				tempData := (*data)[i : en+1]
				tempIndex := f.getIndex(&tempData)
				f.manualStruct(&tempData, tempIndex)
				tempAns += "\"" + tempIndex + "\""
				i = en + 1
			case '[':
				en := f.findMatchRightBrace(data, i+1, len(*data), '[', ']')
				tempData := (*data)[i : en+1]
				tempIndex := f.getIndex(&tempData)
				tempAns += "\"" + tempIndex + "\""
				f.manualList(&tempData, tempIndex)
				i = en + 1
			// 处理非字符串
			default:
				en := f.findNextElement(data, i, len(*data), ',')
				if en == -1 {
					en = f.findNextElement(data, i, len(*data), '}')
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
			en := f.findNextElement(data, i, len(*data), '"')
			tempAns += "\"" + (*data)[i:en] + "\""
			i = en + 1
		case '}':
			i++
		}
	}
	tempAns = "{" + tempAns + "}"
	f.index_data[index] = tempAns
}

func UnFlatted(data string) string {
	flat := NewFlattedData()
	flat.getIndexDataMap(&data)
	return flat.getAns(0)
}

func (f *flatted) getIndexDataMap(input *string) {
	for i := 1; i < len(*input)-1; {
		switch (*input)[i] {
		case '{':
			en := f.findMatchRightBrace(input, i+1, len(*input), '{', '}')
			en++
			f.number_data[f.indexGlobal] = (*input)[i:en]
			f.indexGlobal++
			i = en
		case '[':
			en := f.findMatchRightBrace(input, i+1, len(*input), '[', ']')
			en++
			f.number_data[f.indexGlobal] = (*input)[i:en]
			f.indexGlobal++
			i = en
		case '"':
			en := f.findNextElement(input, i+1, len(*input), '"')
			en++
			f.number_data[f.indexGlobal] = (*input)[i:en]
			f.indexGlobal++
			i = en
		case ',':
			i++
		}
	}
}

func (f *flatted) getAns(index int) string {
	ans := ""
	temp := f.number_data[index]
	if temp[0] == '"' && temp[len(temp)-1] == '"' {
		if _, err := strconv.Atoi(temp[1 : len(temp)-1]); err == nil {
			return temp
		}
	}
	for i := 0; i < len(temp); {
		switch temp[i] {
		case '"':
			if i == 0 || temp[i-1] != '\\' {
				en := f.findNextElement(&temp, i+1, len(temp), '"')
				maybeNumber := temp[i+1 : en]
				number, err := strconv.Atoi(maybeNumber)
				if err == nil {
					ans += f.getAns(number)
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
