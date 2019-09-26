Feature: GET /signout

    Scenario:
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "phone": "13112341234",
                "email": "hatlonely@foxmail.com",
                "firstName": "爽",
                "lastName": "郑"
            }
            """
        When http 请求 GET /signout/d571bda90c2d4e32a793b8a1ff4ff984
        Then http 检查 202
        Then redis not exist "d571bda90c2d4e32a793b8a1ff4ff984"