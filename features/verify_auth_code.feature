Feature: verifyauthcode 校验测试

    Scenario Outline: 校验成功
        Given rediscache.authcode 创建 key: "<email>", code: "<rcode>"
        When 请求 /verifyauthcode, type: "email", email: "<email>", code: "<vcode>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 verify 返回包体 res.body, ok: <ok>, tip: "<tip>"
        Examples:
            | email                  | rcode  | vcode  | status | ok    | tip      |
            | hatlonely1@foxmail.com | 123123 | 123123 | 200    | true  | N/A      |
            | hatlonely1@foxmail.com | 123123 | 123124 | 200    | false | 验证失败 |
