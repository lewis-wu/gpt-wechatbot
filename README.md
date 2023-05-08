# wechatbot
此ChatGPT微信机器人是基于[wechatbot](https://github.com/djun/wechatbot) `fork`出来做的二次开发，原项目只能完成了ChatPrompt功能，即简单的提示完成功能，本人增加一些其他的功能
### 目前实现了以下功能
- 带有聊天上下文的群聊@回复 私聊回复功能(使用OpenAI的”v1/chat/completions"接口替代原来的“v1/completions”)
- 文本编辑完善 (基于OpenAI的“v1/edit”接口开发，对中文支持不是很好，有点鸡肋)
- 图片生成
- 图片变体
 
# 注册openai
chatGPT注册可以参考[这里](https://juejin.cn/post/7173447848292253704)

# 安装使用
``` txt
# 获取项目
git clone https://github.com/lewis-wu/gpt-wechatbot.git

# 进入项目目录
cd gpt-wechatbot

# 复制配置文件
copy config.dev.json config.json

# 启动项目
go run main.go

启动前需替换config中的api_key（若不改其他的配置则使用默认配置，项目仍可启动）
```

# 配置文件说明
 - api_key : OpenAI的访问token
 - auto_pass : 是否自动通过添加好友申请
 - proxy: 网络代理，默认值为空(国内一般无法直接访问opeanai的接口，这一项国内用户基本都需要配置。格式`http(s)://xxx.xxx:port`)
 - chat_max_context: 保存的最大聊天上下文记录数，默认值`2` (不建议设置太大，此参数太大将导致token数消耗过快)
 - chat_ttl_time: 上下文保存的分钟数，默认值`10`(超过此时间，聊天上下文将被清空)
 - gpt_time_out: openai的接口请求超时分钟数，默认值为`60`
 - generate_image_keyword: 图片生成的关键词前缀，默认值是`[P]`
 - text_edit_keyword: 文本编辑完善的关键词前缀,默认值是`[TE]`
 - text_edit_separator: 文本编辑完善的原文本与完善建议的分隔符,默认值是`[TES]`
 - gpt_limit_per_minute: 每人每分钟最多与聊天机器人的互动次数,默认值是`3`
 - image_variation_keyword: 图片变体的关键词，默认值是`[PV]`
 - image_variation_chat_ttl: 在发送图片变体的关键词后，发送源图片的最大间隔时间秒数，默认值是`60`

**<font color="green">Note: </font>**
> 1. 若更换了登录人，或登录失效，需要重新登录时需要将`storage.json`这个文件删除,否则无法扫码登录（这个文件是保存微信登录信息的，在登录时自动生成）
> 2. 如果在机器人已登录后重启，可能会拉取到之前的历史消息，导致重复接收到ChatGPT回复的消息(重复请求OpenAI接口)。这个也可以通过删除`storage.json`文件然后重新扫码登录来解决。
