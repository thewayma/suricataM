package utils

import (
	"fmt"
	"sort"
	"strings"
)

// key升序排序
func KeysOfMap(m map[string]string) []string {
    keys := make(sort.StringSlice, len(m))
    i := 0
    for key := range m {
        keys[i] = key
        i++
    }

    keys.Sort()
    return []string(keys)
}

// 以key排序, 返回用","拼接而成的字符串
func SortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	if size == 1 {
		for k, v := range tags {
			return fmt.Sprintf("%s=%s", k, v)
		}
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	ret := make([]string, size)
	for j, key := range keys {
		ret[j] = fmt.Sprintf("%s=%s", key, tags[key])
	}

	return strings.Join(ret, ",")
}

// 以","拼接的字符串(tag1=v1,tag2=v2,tag3=v3...)转换成map
func DictedTagstring(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}
	s = strings.Replace(s, " ", "", -1)

	tag_dict := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {						//!< tag1=v1=v2, 切分成 key:tag1, value:v1=v2
			tag_dict[tag_pair[0]] = tag_pair[1]
		}
	}
	return tag_dict
}

// 以","拼接的字符串(tag1=v1,tag2=v2,tag3=v3...)转换成map
func SplitTagsString(s string) (err error, tags map[string]string) {
	err = nil
	tags = make(map[string]string)

	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		return
	}

	tagSlice := strings.Split(s, ",")
	for _, tag := range tagSlice {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {
			tags[tag_pair[0]] = tag_pair[1]
		} else {									//!< 若存在tag1=v1=v2, 则返回错误
			err = fmt.Errorf("bad tag %s", tag)
			return
		}
	}

	return
}
