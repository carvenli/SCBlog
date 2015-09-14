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
		length = rl
	}

	return string(rs[0:length])
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
	// 定义一个Redirect列表
	var links []models.SC_Redirect

	// 获取所有内链
	models.GetAllByQuery(models.DbRedirect, nil, &links)

	// 定义一个标签列表
	var tags []models.SC_Tag

	// 获取所有标签
	models.GetAllByQuery(models.DbTag, nil, &tags)

	// 定义一个正则, 用以确认关键词是否已经存在链接
	re, _ := regexp.Compile("(?is)<a\b[^>]*>(.*?)</a>")
	// 在文本中搜索所有链接
	mc := re.FindAllStringSubmatch(str, -1)

	// 对匹配结果进行循环处理
	for _, m := range mc {
		// 循环内链
		for i := 0; i < len(links); i++ {
			// 如果内链内容与链接内容相同
			if strings.ToLower(links[i].Caption) == m[1] {
				// 去除此内链
				links = append(links[:i], links[i+1:]...)
			}
		}

		// 循环标签
		for i := 0; i < len(tags); i++ {
			// 如果标签内容与链接内容相同
			if strings.ToLower(tags[i].Caption) == m[1] {
				// 去除标签
				tags = append(tags[:i], tags[i+1:]...)
			}
		}
	}

	// 对内链进行循环
	for _, l := range links {
		// 循环标签
		for i := 0; i < len(tags); i++ {
			// 如果名称与内链名称相同
			if tags[i].Caption == l.Caption {
				// 从标签中去除此项
				tags = append(tags[:i], tags[i+1:]...)
			}
		}
	}

	return ""
}
