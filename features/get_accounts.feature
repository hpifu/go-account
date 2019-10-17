Feature: GET /accounts

    Scenario: accounts
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id IN (1,2,3)
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345671", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (2, "13112345672", "hatlonely2@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (3, "13112345673", "hatlonely3@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set string "godtoken"
            """
            1c15b6b0b18aa0d3a5d2de37484f992c
            """
        When http 请求 GET /accounts
            """
            {
                "params": {
                    "ids": [
                        1,
                        2,
                        3
                    ]
                },
                "header": {
                    "Authorization": "1c15b6b0b18aa0d3a5d2de37484f992c"
                }
            }
            """
        Then http 检查 200
            """
            {
                "json": [
                    {
                        "phone": "13112345671",
                        "email": "hatlonely1@foxmail.com",
                        "firstName": "悟空",
                        "lastName": "孙"
                    },
                    {
                        "phone": "13112345672",
                        "email": "hatlonely2@foxmail.com",
                        "firstName": "悟空",
                        "lastName": "孙"
                    },
                    {
                        "phone": "13112345673",
                        "email": "hatlonely3@foxmail.com",
                        "firstName": "悟空",
                        "lastName": "孙"
                    }
                ]
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id IN (1,2,3)
            """
        Given redis del "1c15b6b0b18aa0d3a5d2de37484f992c"


    Scenario: accounts
        Given redis set string "godtoken"
            """
            1c15b6b0b18aa0d3a5d2de37484f992c
            """
        When http 请求 GET /accounts
            """
            {
                "params": {
                    "ids": [
                        1,
                        2,
                        3
                    ]
                },
                "header": {
                    "Authorization": "wrong token"
                }
            }
            """
        Then http 检查 401
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id IN (1,2,3)
            """
        Given redis del "1c15b6b0b18aa0d3a5d2de37484f992c"