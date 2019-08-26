Feature: update 登陆测试

    Scenario Outline: update name
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, firstName: "<firstname>", lastName: "<lastname>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", firstname: "<firstname>", lastname: "<lastname>", birthday: "1992-01-01", gender: 1
        Then 检查 mysqldb.accounts, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "<firstname>", lastname: "<lastname>", birthday: "1992-01-01", gender: 1
        Examples:
            | firstname | lastname | status | ok   | err |
            | 八戒      | 猪       | 200    | true | N/A |
            | 僧        | 沙       | 200    | true | N/A |

    Scenario Outline: update phone
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, phone: "<phone>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "<phone>", email: "hatlonely1@foxmail.com", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Then 检查 mysqldb.accounts, 存在记录 phone: "<phone>", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Examples:
            | phone       | status | ok   | err |
            | 13111112222 | 200    | true | N/A |

    Scenario Outline: update email
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, email: "<email>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "13112345678", email: "<email>", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Then 检查 mysqldb.accounts, 存在记录 phone: "13112345678", email: "<email>", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Examples:
            | email                  | status | ok   | err |
            | hatlonely2@foxmail.com | 200    | true | N/A |

    Scenario Outline: update birthday
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, birthday: "<birthday>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", firstname: "孙", lastname: "悟空", birthday: "<birthday>", gender: 1
        Then 检查 mysqldb.accounts, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "<birthday>", gender: 1
        Examples:
            | birthday   | status | ok   | err |
            | 1991-11-11 | 200    | true | N/A |

    Scenario Outline: update gender
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, gender: "<gender>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: <gender>
        Then 检查 mysqldb.accounts, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: <gender>
        Examples:
            | gender | status | ok   | err |
            | 0      | 200    | true | N/A |
            | 1      | 200    | true | N/A |
            | 2      | 200    | true | N/A |

    Scenario Outline: update password
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, password: "<password>", oldpassword: "<oldpassword>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 update 返回包体 res.body, ok: <ok>, err: "<err>"
        Then 检查 rediscache.token, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Then 检查 mysqldb.accounts, 存在记录 phone: "13112345678", email: "hatlonely1@foxmail.com", password: "<newpassword>", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        Examples:
            | password | oldpassword | newpassword | status | ok    | err      |
            | 11112222 | 12345678    | 11112222    | 200    | true  | N/A      |
            | 11112222 | 12341234    | 12345678    | 200    | false | 密码错误 |

    Scenario Outline: update name 异常
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, firstname: "<firstname>", lastname: "<lastname>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | firstname                         | lastname                          | status | body         |
            | N/A                               | hatlonely                         | 400    | 必要字段     |
            | hatlonely                         | N/A                               | 400    | 必要字段     |
            | N/A                               | N/A                               | 400    | 必要字段     |
            | hatlonley                         | 123456789012345678901234567890123 | 400    | 至多32个字符 |
            | 123456789012345678901234567890123 | hatlonely                         | 400    | 至多32个字符 |

    Scenario Outline: update phone 异常
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, phone: "<phone>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | phone     | status | body           |
            | N/A       | 400    | 必要字段       |
            | hatlonley | 400    | 无效的电话号码 |

    Scenario Outline: update email 异常
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, email: "<email>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | email     | status | body       |
            | N/A       | 400    | 必要字段   |
            | hatlonley | 400    | 无效的邮箱 |

    Scenario Outline: update birthday 异常
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, birthday: "<birthday>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | birthday   | status | body         |
            | N/A        | 400    | 必要字段     |
            | 19920515   | 400    | 日期格式错误 |
            | 1992/05/15 | 400    | 日期格式错误 |
            | 2049-10-01 | 400    | 日期超过范围 |
            | 1849-10-01 | 400    | 日期超过范围 |

    Scenario Outline: update password 异常
        Given mysqldb.accounts 创建用户, phone: "13112345678", email: "hatlonely1@foxmail.com", password: "12345678", firstname: "孙", lastname: "悟空", birthday: "1992-01-01", gender: 1
        When 请求 /signin, username: "13112345678", password: "12345678"
        When 请求 /update, password: "<password>", oldpassword: "<oldpassword>"
        Then 检查状态码 res.status_code: <status>
        Then 检查返回包体 res.body，包含字符串 "<body>"
        Examples:
            | password | oldpassword | status | body        |
            | N/A      | 12345678    | 400    | 必要字段    |
            | 123456   | 12345678    | 400    | 至少8个字符 |
