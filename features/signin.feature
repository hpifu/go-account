Feature: POST /signin

    Scenario: case1
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        When http 请求 POST /signin
            """
            {
                "json": {
                    "username": "13112345678",
                    "password": "12345678"
                }
            }
            """
        Then http 检查 200
            """
            {
                "cookies": {
                    "token": "exist"
                }
            }
            """
        When http 请求 POST /signin
            """
            {
                "json": {
                    "username": "hatlonely1@foxmail.com",
                    "password": "12345678"
                }
            }
            """
        Then http 检查 200
            """
            {
                "cookies": {
                    "token": "exist"
                }
            }
            """
        When http 请求 POST /signin
            """
            {
                "json": {
                    "username": "hatlonely1@foxmail.com",
                    "password": "xxxxxxxx"
                }
            }
            """
        Then http 检查 403
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
