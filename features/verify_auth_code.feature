Feature: GET /verify/authcode/:type

    Scenario:
        Given redis set string "ac_hatlonely@foxmail.com"
            """
            123123
            """
        When http 请求 GET /verify/authcode/email
            """
            {
                "params": {
                    "email": "hatlonely@foxmail.com",
                    "code": "123123"
                }
            }
            """
        Then http 检查 200
        Given redis del "ac_hatlonely@foxmail.com"

    Scenario:
        Given redis set string "ac_hatlonely@foxmail.com"
            """
            123123
            """
        When http 请求 GET /verify/authcode/email
            """
            {
                "params": {
                    "email": "hatlonely@foxmail.com",
                    "code": "123456"
                }
            }
            """
        Then http 检查 200
            """
            {
                "text": "验证失败"
            }
            """
        Given redis del "ac_hatlonely@foxmail.com"

    Scenario:
        Given redis del "ac_hatlonely@foxmail.com"
            """
            123123
            """
        When http 请求 GET /verify/authcode/email
            """
            {
                "params": {
                    "email": "hatlonely@foxmail.com",
                    "code": "123456"
                }
            }
            """
        Then http 检查 200
            """
            {
                "text": "验证码不存在"
            }
            """