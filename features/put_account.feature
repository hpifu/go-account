Feature: PUT /account/:token/:field

    Scenario: case update phone
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/phone
            """
            {
                "json": {
                    "phone": "13112341234"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "phone": "13112341234"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "phone": "13112341234"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case email
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/email
            """
            {
                "json": {
                    "email": "hatlonely2@foxmail.com"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "email": "hatlonely2@foxmail.com"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "email": "hatlonely2@foxmail.com"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"


    Scenario: case birthday
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/birthday
            """
            {
                "json": {
                    "birthday": "1994-12-31"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "birthday": "1994-12-31"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "birthday": "1994-12-31"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case gender
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/gender
            """
            {
                "json": {
                    "gender": 0
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "gender": 0
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "gender": 0
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case name
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/name
            """
            {
                "json": {
                    "firstName": "八戒",
                    "lastName": "猪"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "first_name": "八戒",
                "last_name": "猪"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "firstName": "八戒",
                "lastName": "猪"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case password
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1,
                "password": "12345678"
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/password
            """
            {
                "json": {
                    "password": "11112222",
                    "oldPassword": "12345678"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "password": "11112222"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "password": "11112222"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case wrong password
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1 OR email="hatlonely2@foxmail.com"
            """
        Given mysql 执行
            """
            INSERT INTO accounts (id, phone, email, password, first_name, last_name, birthday, gender)
            VALUES (1, "13112345678", "hatlonely1@foxmail.com", "12345678", "悟空", "孙", "1992-01-01", 1)
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 1,
                "password": "12345678"
            }
            """
        When http 请求 PUT /account/d571bda90c2d4e32a793b8a1ff4ff984/password
            """
            {
                "json": {
                    "password": "11112222",
                    "oldPassword": "1234567890"
                }
            }
            """
        Then http 检查 403
        Then mysql 检查 "SELECT * FROM accounts WHERE id=1"
            """
            {
                "password": "12345678"
            }
            """
        Then redis get object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "password": "12345678"
            }
            """
        Given mysql 执行
            """
            DELETE FROM accounts WHERE id=1
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
