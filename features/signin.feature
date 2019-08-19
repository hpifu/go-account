Feature: signin 登陆测试

    Scenario Outline: 登陆成功
        Given mysqldb.accounts 创建用户, phone: "<phone>", email: "<email>", password: "<password>", firstname: "<firstname>", lastname: "<lastname>", birthday: "<birthday>", gender: <gender>
        When 请求 /signin, username: "<username>", password: "<password>"
        Then 检查状态码 res.status_code: <status>
        Then 检查登陆返回 cookie
        Then 检查登陆返回包体 res.body, valid: <valid>, tokenlen: <tokenlen>
        Then 检查 rediscache.token, 存在记录 phone: "<phone>", email: "<email>", firstname: "<firstname>", lastname: "<lastname>", birthday: "<birthday>", gender: <gender>
        Examples:
            | username               | password | status | valid | tokenlen | phone       | email                  | password | firstname | lastname | birthday   | gender |
            | 13112345678            | 12345678 | 200    | true  | 32       | 13112345678 | hatlonely1@foxmail.com | 12345678 | 悟空      | 孙       | 1992-01-01 | 1      |
            | hatlonely1@foxmail.com | 12345678 | 200    | true  | 32       | 13112345678 | hatlonely1@foxmail.com | 12345678 | 悟空      | 孙       | 1992-01-01 | 1      |

    Scenario Outline: 登陆失败
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "<username>", password: "<password>"
        Then 检查状态码 res.status_code: <status>
        Then 检查登陆返回包体 res.body, valid: <valid>, tokenlen: <tokenlen>
        Examples:
            | username                  | password       | status | valid | tokenlen |
            | notexistsuser@foxmail.com | 12345678       | 200    | false | 0        |
            | hatlonely1@foxmail.com    | wrong_password | 200    | false | 0        |

    Scenario Outline: 异常登陆
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "<username>", password: "<password>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | username | password | status | body     |
            | N/A      | 12345678 | 400    | 必要字段 |
            | N/A      | N/A      | 400    | 必要字段 |
