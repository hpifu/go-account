FORMAT: 1A
HOST: https://api.account.hatlonely.com/

# account

账户接口

## 获取账户 [/account/{token}]

### 获取账户 [GET]

+ Response 200 (application/json)

        {
            "id": 1,
            "email": "hatlonely@foxmail.com",
            "phone": "13810245740",
            "firstName": "乐",
            "lastName": "贺",
            "birthday": "1992-05-15",
            "password": "s11209768hls8",
            "gender": 1,
            "avatar": "hatlonely-christmas.png"
        }

+ Response 204
+ Response 400

## 创建账户 [/account]

### 创建账户 [PUT]

+ Request (application/json)

        {
            "email": "hatlonely@foxmail.com",
            "phone": "13810245740",
            "firstName": "乐",
            "lastName": "贺",
            "birthday": "1992-05-15",
            "password": "12345678",
            "gender": 1,
            "avatar": "hatlonely-christmas.png"
        }

+ Response 201