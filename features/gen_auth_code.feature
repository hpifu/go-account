Feature: genautocode 生成验证码测试

    Scenario Outline: genauthcode 成功
        When 请求 /genauthcode, type: "email", email: "<email>", firstname: "<firstname>", lastname: "<lastname>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 genauthcode 返回包体 res.body, ok: <ok>
        Then 检查 rediscache.authcode, 存在记录 key: "<email>"
        Examples:
            | email                 | firstname | lastname | status | ok   |
            | hatlonely@foxmail.com | 爽        | 郑       | 200    | true |
