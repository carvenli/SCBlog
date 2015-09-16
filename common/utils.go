package common

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"github.com/ylqjgm/SCBlog/models"
	"gopkg.in/mgo.v2/bson"
)

// 生成Gravatar头像URL
func Gravatar(email string, size int) string {
	// 如果大小小于1
	if size < 1 {
		// 设置大小为80
		size = 80
	}

	// 将email地址去除空格并转换为小写
	email = strings.ToLower(strings.TrimSpace(email))
	// 创建一个MD5
	hash := md5.New()
	// 将email转换为MD5
	hash.Write([]byte(email))

	// 返回生成的Gravatar头像URL
	return fmt.Sprintf("https://secure.gravatar.com/avatar/%x?s=%d", hash.Sum(nil), size)
}

// 过滤Html
func filterHtml(str string) string {
	// 将Html标签全部转换为小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllStringFunc(str, strings.ToLower)

	// 去除Style
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	str = re.ReplaceAllString(str, "")

	// 去除Script
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	str = re.ReplaceAllString(str, "")

	// 去除所有尖括号内的Html代码, 并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "\n")

	// 去除连续的换行符
	re, _ = regexp.Compile("\\S\\s{2,}")
	str = re.ReplaceAllString(str, "\n")

	return str
}

// 获取预览内容
func Preview(str string, length int) string {
	// 先过滤Html
	str = filterHtml(str)
	// 将字符串转换为rune列表
	rs := []rune(str)
	// 获取长度
	rl := len(rs)

	// 如果截取长度大于字符串长度
	if length > rl {
		// 截取长度等于字符串长度
		str = string(rs[0:rl])
	} else {
		str = string(rs[0:length]) + "..."
	}

	return strings.Replace(str, "\n", "", -1)
}

// 获取Id
func GetId(id bson.ObjectId) string {
	return id.Hex()
}

// 获取Tag列表
func GetTagSlug(caption string) string {
	var tag models.SC_Tag

	models.GetOneByQuery(models.DbTag, bson.M{"caption": caption}, &tag)

	return tag.Slug
}

// 自动获取slug
func GetSlug(str string, isslug bool) string {
	slug := Convert(str, isslug)
	re, _ := regexp.Compile(`\\pP|\\pS|\\u3002|\\uff1b|\\uff0c|\\uff1a|\\u201c|\\u201d|\\uff08|\\uff09|\\u3001|\\uff1f|\\u300a|\\u300b|–|\\p{P}|\\f|\\n|\\r|\\t|\\v|\\x85|\\p{Z}]`)
	slug = re.ReplaceAllString(slug, "")
	slug = strings.TrimSpace(slug)
	if isslug {
		slug = strings.Replace(slug, " ", "", -1)
		slug = strings.Replace(slug, "--", "-", -1)
	} else {
		slug = strings.Replace(slug, " ", "", -1)
		slug = strings.Replace(slug, "-", "", -1)
	}

	if slug[len(slug)-1:] == "-" {
		slug = slug[:len(slug)-1]
	}

	return strings.ToLower(slug)
}

// 替换关键字
func ReplaceKeys(str string) string {
	// 获取关键字列表
	keys := models.Keys
	// 定义一个正则, 用以确认关键词是否已经存在链接
	re, _ := regexp.Compile(`(?is)<a\b[^>]*>(.*?)</a>|<a\b[^>]*"(.*?)"*>`)
	// 在文本中搜索所有链接
	mc := re.FindAllStringSubmatch(str, -1)

	// 对匹配结果进行循环处理
	for _, m := range mc {
		// 循环关键字
		for _, k := range keys {
			// 如果关键字名称与链接名称相同
			if strings.ToLower(k.Caption) == m[1] {
				// 去除此关键字
				delete(keys, k.Caption)
			}
		}
	}

	// 对关键字进行循环
	for _, k := range keys {
		// 定义一个正则, 忽略大小写
		r, _ := regexp.Compile("(?is)(^'\"" + k.Caption + ")")
		// 如果是标签
		if k.IsTag {
			str = r.ReplaceAllString(str, fmt.Sprintf(`<a href="/tag/%s" title="%s">%s</a>`, k.Slug, "$0", "$0"))
		} else {
			str = r.ReplaceAllString(str, fmt.Sprintf(`<a href="/go/%s" target="_blank" title="%s">%s</a>`, k.Caption, "$0", "$0"))
		}
	}

	return str
}
