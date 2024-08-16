package wordseg

import (
	"fmt"
	"testing"
)

func TestWordseg(t *testing.T) {
	fmt.Println(_default_word_seg.words['共'])
	fmt.Println(Split("这项技术的的原始版本是图灵当年和他的助手"))
	fmt.Println(Split("法轮大法好"))
	fmt.Println(Split("我操你大爷"))
	fmt.Println(Split("共产党傻屄"))
	fmt.Println(Split("你妈逼"))
	fmt.Println(Split("威胁的"))
	fmt.Println(Split("我来到北京清华大学"))
	fmt.Println(Split("永和服装饰品有限公司"))
	fmt.Println(Split("小明硕士毕业于中国科学院计算所后在日本京都大学深造"))
	fmt.Println(Split("you mother fucker"))
	fmt.Println(Split("901872340918723,zx,cv!@#$%^&*()_LKNAN<XZMNVC"))
	fmt.Println(Split("他称自己昨日给上海市纪委提供了6月9日视频的完整证据。"))
	fmt.Println(Split("工信处女干事每月经过下属科室都要亲口交代二十四口交换机等技术性器件的安装工作然后工信处女干事每月经过下属科室都要亲口交代二十四口交换机等技术性器件的安装工作"))
}
