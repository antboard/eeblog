# eeblog
一个为电子工程师打造的博客

# 缘起
    我本来有个blog, 但是写起来比较不舒服. 最近看见postgresql比较好,就想要不要搞点什么事情.
    我是一个码农,也是一个电工.所以对互联网技术有着不同的需求. 
    我比较懒, 但是懒人改变世界. 突然有个想法,那就写个博客吧.
    博客是个好东西,只要你写了, 你就会入门一个新语言,比如golang...
    当然我现在写这个博客的时候,可不是入门者,而是一个极客.
    我想,只要能跟我玩下去,必然很脑洞打开~~
    大家一起来.有机会建个wechat群.

# 设计基调

这个是本项目的指导思想, 估计很多人会不喜欢

* db: 选择postgresql, 但是为了简单,暂时是用xorm来存取.
* golang: 本项目用golang开发. 
* gin: 简单稳定够用即可.
* wasm: 本项目前端使用golang编译成wasm来替代js.
* go-vue: https://github.com/norunners/vue 使用这个golang版本的vue来完成基本版面.
* markdown: 这个应该是博客的核心. 也就是说我们要搞定一个EE扩展版本的markdown编辑器.具体需求暂时还不定.等基本页面完成再说
* svg: 这个是服务markdown的, 很多EE的东西,需要用svg图形展现.
* [临时] 初期为了更快的展现,暂时使用jq和bootstrap来展现.第二阶段切换到wasm中去

# 测试用例

请到docker网站查看测试用部署:
https://hub.docker.com/r/icecut/eeblog

线上测试版本,等markdown功能添加完成后开放.