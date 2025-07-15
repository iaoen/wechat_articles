# wechat_articles 公众号文章采集爬虫
稳定工作2年的微信公众号爬虫 Based on golang and python 微信公众号采集 Python爬虫 公众号采集 公众号文章爬虫

## 必备条件
需要申请个人公众号

## 使用流程
1.下载编译好的exe文件
2.使用test.py测试接口可用性

```python
接口功能介绍

1.在软件目录下生成登录二维码，使用微信扫码登录
get请求http://127.0.0.1:12312/login接口

2.搜公众号,一般来说,名称一致,搜到的就在第一个
post请求http://127.0.0.1:12312/search
请求体json格式{
    "cookie": "rand_ixxxxx657",
    "query": "浙江发布",
}

3.查历史文章列表
post请求http://127.0.0.1:12312/appmsg
请求体json格式{
    "cookie": "rand_ixxxxx657",
    "fakeid": "MzA4ODY3MjkxNA==",  # 搜公众号返回的id
    "page": "2",  # 第几页
}
```
