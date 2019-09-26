Feature: POST /authcode

    Scenario: case
        When http 请求 POST /authcode/email
            """
            {
                "json": {
                    "email": "hatlonely@foxmail.com",
                    "firstName": "爽",
                    "lastName": "郑"
                }
            }
            """
        Then http 检查 201
        Then redis get "ac_hatlonely@foxmail.com"
            """
            {}
            """
        Given redis del "ac_hatlonely@foxmail.com"
