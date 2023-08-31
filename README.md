Flatted算法golang版本

原库实现了JavaScript和PHP版本，本库实现了golang版本
具有大量相同value的json将会获得更好的压缩效率

算法原理：
将json串中的value进行索引编号，并用序号替换相应的value。相同的value可以被合并，只保存一份副本，达到压缩的效果。

使用说明：
func FlattedFromStruct(input interface{}) (string, error)
提供了将结构体序列化后转化为flatted串的能力，输入参数为需要被压缩的结构体，返回被压缩后的字符串
func UnFlattedToStruct(str string, output interface{}) error
提供了将flatted串转化为结构体的能力，输入参数为flatted字符串与存放结果的结构体，实参需要传递引用
func Flatted(data string) string
提供了将序列化后的json串压缩为flatted串的能力
func UnFlatted(data string) string
提供了将flatted串解压缩为序列化后的json串的能力

不压缩name只压缩value是为了与原库兼容
