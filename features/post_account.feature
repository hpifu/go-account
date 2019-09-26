Feature: account [POST] 新建账号测试

    Scenario: account case1
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """
        When http 请求 POST /account
            """
            {
                "json": {
                    "phone": "13112341234",
                    "email": "hatlonely@foxmail.com",
                    "password": "12341234"
                }
            }
            """
        Then http 检查 201
        Then mysql 检查 "SELECT * FROM accounts WHERE email='hatlonely@foxmail.com'"
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "password": "12341234"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """

    Scenario: account case2
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """
        When http 请求 POST /account
            """
            {
                "json": {
                    "email": "hatlonely@foxmail.com",
                    "password": "12341234"
                }
            }
            """
        Then http 检查 201
        Then mysql 检查 "SELECT * FROM accounts WHERE email='hatlonely@foxmail.com'"
            """
            {
                "email": "hatlonely@foxmail.com",
                "password": "12341234"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """

    Scenario: account case3
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """
        When http 请求 POST /account
            """
            {
                "json": {
                    "email": "hatlonely@foxmail.com",
                    "password": "12341234",
                    "lastName": "郑",
                    "firstName": "郑",
                    "birthday": "1992-05-15",
                    "gender": 1
                }
            }
            """
        Then http 检查 201
        Then mysql 检查 "SELECT * FROM accounts WHERE email='hatlonely@foxmail.com'"
            """
            {
                "email": "hatlonely@foxmail.com",
                "password": "12341234",
                "last_name": "郑",
                "first_name": "郑",
                "gender": 1
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts
            WHERE email="hatlonely@foxmail.com"
            """