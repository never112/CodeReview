项目主要功能
1.项目golang语言实现。
2.通过webhook监听github的pullrequest事件，在监听到pullrequest事件后，git clone 当前项目代码。
3.在git clone的代码目录下，通过执行 /review 命名对代码进行review。
4.review完成后，将review结果通过调用github的api，将结果写入到github对应的pullrequest的 comment 上。