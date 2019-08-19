Feature: verify 校验测试

    Scenario Outline: 校验成功
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /verify, field: "<field>", value: "<value>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 verify 返回包体 res.body, ok: <ok>, tip: "<tip>"
        Examples:
            | field    | value                  | status | ok    | tip            |
            | phone    | 13112345678            | 200    | false | 电话号码已存在 |
            | email    | hatlonely1@foxmail.com | 200    | false | 邮箱已存在     |
            | username | hatlonely2@foxmail.com | 200    | false | 账号不存在     |
            | username | 13811111111            | 200    | false | 账号不存在     |


    Scenario Outline: 异常校验
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /verify, field: "<field>", value: "<value>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | field    | value    | status | body     |
            | N/A      | 12345678 | 400    | 必要字段 |
            | phone    | N/A      | 400    | 必要字段 |
            | password | 12345678 | 400    | 必须在   |
