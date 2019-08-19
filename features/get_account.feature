Feature: getaccount 获取账号测试

    Scenario Outline: 获取成功
        Given rediscache.token 创建 token: "<token>", phone: "<phone>", email: "<email>", firstname: "<firstname>", lastname: "<lastname>", birthday: "<birthday>", gender: <gender>
        When 请求 /getaccount, token: "<token>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 token 返回包体 res.body, ok: <ok>, phone: "<phone>", email: "<email>", firstname: "<firstname>", lastname: "<lastname>", birthday: "<birthday>", gender: <gender>
        Examples:
            | status | token                            | ok   | phone       | email                  | firstname | lastname | birthday   | gender |
            | 200    | d571bda90c2d4e32a793b8a1ff4ff984 | true | 13145678901 | hatlonely1@foxmail.com | 悟空      | 孙       | 1992-01-01 | 1      |


    Scenario Outline: 获取失败
        When 请求 /getaccount, token: "<token>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 token 返回包体 res.body, ok: false
        Examples:
            | status | token         | ok    |
            | 200    | notexisttoken | false |
