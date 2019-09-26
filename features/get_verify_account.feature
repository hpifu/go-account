Feature: POST /verify/account

    Scenario: case nonexist phone
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "phone",
                    "value": "13112345679"
                }
            }
            """
        Then http 检查 200
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """


    Scenario: case exist phone
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "phone",
                    "value": "13112345678"
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "电话号码已存在"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """

    Scenario: case exist email
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "email",
                    "value": "hatlonely1@foxmail.com"
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "邮箱已存在"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """

    Scenario: case exist username
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "username",
                    "value": "hatlonely1@foxmail.com"
                }
            }
            """
        Then http 检查 200
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "username",
                    "value": "13112345678"
                }
            }
            """
        Then http 检查 200
        When http 请求 GET /verify/account
            """
            {
                "params": {
                    "field": "username",
                    "value": "13112345679"
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "账号不存在"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """