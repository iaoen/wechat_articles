import requests


# 登录,会在同目录生成qrcode.jpg,请用微信扫描登录
# res = requests.get("http://127.0.0.1:12312/login")
# print(res.text)
# 成功返回{"code":0,"data":"rand_ixxxxx657","msg":"success"}
# data值就是登录后的cookie


# 搜公众号,一般来说,名称一致,搜到的就在第一个
# data = {
#     "cookie": "rand_ixxxxx657",
#     "query": "浙江发布",
# }
# res = requests.post("http://127.0.0.1:12312/search", json=data)
# print(res.text)
# 成功返回{"code":200,"data":[{"alias":"公众号代码","fakeid":"公众号id，查文章要用到","nickname":"公众号名称","round_head_img":"logo图片","service_type":0,"signature":"描述"}],"msg":"success"}
# "fakeid":"公众号id，查文章要用到"


# 查文章列表
data = {
    "cookie": "rand_ixxxxx657",
    "fakeid": "MzA4ODY3MjkxNA==",  # 搜公众号返回的id
    "page": "2",  # 第几页
}
res = requests.post("http://127.0.0.1:12312/appmsg", json=data)
print(res.text)
# 成功返回文章列表,link链接,time发布时间,title标题
