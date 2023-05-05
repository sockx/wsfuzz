## wsfuzz🎉🎉🎉

> 渗透测试中偶尔会遇到需要对`websocket`进行ID遍历，或者爆破等操作，就顺手做了一个基础工具。

### 使用方法
> go run . fuzz -f req.txt

### 其他说明
* 工具会自动忽略证书限制
* 工具仅使用数字ID作为爆破规则，如果需要其他字符需要对工具进行修改
* 例子中需要填充的位置使用了`{CG}`进行标记
* req.txt 中 `GET PATH HTTP/1.1` 与 `headers`需要一个空行，详情查看`ex_req.txt`

### 例子
![ex_img](/docs/image.png)
