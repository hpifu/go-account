Feature: GET /account/:token

    Scenario: account
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "firstName": "爽",
                "lastName": "郑"
            }
            """
        When http 请求 GET /account/d571bda90c2d4e32a793b8a1ff4ff984
        Then http 检查 200
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "firstName": "爽",
                "lastName": "郑"
            }
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: account
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "firstName": "爽",
                "lastName": "郑"
            }
            """
        When http 请求 GET /account/d571bda90c2d4e32a793b8a1ff4ff983
        Then http 检查 401
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "firstName": "爽",
                "lastName": "郑"
            }
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
