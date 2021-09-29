package main

import (
	"bufio"
	//"encoding/json"
	"log"
	"os"
)

type filterInfo struct{
	//Str string `json:"str"`
	SubFilterList map[rune]*filterInfo	`json:"sub_filter_list"`
	IsEnd bool	`json:"is_end"`
}

var (
	err error
	fillRune rune = rune('*')
	filterList map[rune]*filterInfo
)

func main() {
	initFilter()

	oldStr := "怀疑人生"
	newStr := checkFilter(oldStr)
	log.Println(oldStr, newStr)

	oldStr = "xiaoen123456"
	newStr = checkFilter(oldStr)
	log.Println(oldStr, newStr)
	
	oldStr = "K　粉123"
	newStr = checkFilter(oldStr)
	log.Println(oldStr, newStr)
}

func initFilter() {
	filterList = make(map[rune]*filterInfo)
	f, err := os.Open("./filter.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		txRunes := []rune(text)

		//log.Println(text)

		fi, ok := filterList[txRunes[0]]
		if !ok {
			fi = &filterInfo{
				//Str:string(txRunes[0]),
				SubFilterList:make(map[rune]*filterInfo),
			}
			filterList[txRunes[0]] = fi
		}

		isEnd := initFilterFor(fi, txRunes, 1)
		if !fi.IsEnd && isEnd {
			fi.IsEnd = true
		}
	}

	/*for k, v := range filterList{
		log.Println(string(k))
		b, _ := json.Marshal(v)
		log.Println(string(b))
	}*/
	
}

func initFilterFor(pfi *filterInfo, runes []rune, index int) bool {
	//log.Println(string(runes), index)
	if len(runes) <= index {
		return true
	}

	fi, ok := pfi.SubFilterList[runes[index]]
	if !ok {
		//log.Println(string(runes), string(runes[index]))
		fi = &filterInfo{
			//Str:string(runes[index]),
			SubFilterList:make(map[rune]*filterInfo),
		}
	    pfi.SubFilterList[runes[index]] = fi
	}

	//同节点如果已经有过IsEnd==true,则不能再被改变成false
	isEnd := initFilterFor(fi, runes, index + 1)
	if !fi.IsEnd && isEnd {
		fi.IsEnd = true
	}
	return false
}

func checkFilter(cs string) string {
	checkRunes := []rune(cs)
	newRunes := make([]rune, len(checkRunes))
	for i := 0; i < len(checkRunes); i++{
		if fi, ok := filterList[checkRunes[i]]; ok {
			has, offset := checkFilterFor(fi, checkRunes[i:], 0)
			//log.Println(cs, i, offset, has, string(checkRunes[i]))
			if has {
				//3个字符的才替换成***,否则一个字符替换一个*
				if offset >= 2 {
				    newRunes = append(newRunes, fillRune, fillRune, fillRune)
				} else {
					for n := 0; n <= offset; n++ {
						newRunes = append(newRunes, fillRune)
					}
				}
			} else {
				for n := 0; n <= offset; n++ {
		    		newRunes = append(newRunes, checkRunes[i+n])
				}
			}
			i += offset
		} else {
		    newRunes = append(newRunes, checkRunes[i])
		}
	}
	return string(newRunes)
}

func checkFilterFor(fi *filterInfo, runes []rune, index int) (bool, int) {
	if index + 1 < len(runes) {
    	if sfi, ok := fi.SubFilterList[runes[index + 1]]; ok {
			//log.Println(string(runes[index+1]), sfi)
	    	return checkFilterFor(sfi, runes, index + 1)
		}
	}

	if fi.IsEnd {
		//log.Println(string(runes), string(runes[index]), fi)
		return true, index
	}
	return false, index
}
