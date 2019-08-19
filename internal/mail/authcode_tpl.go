package mail

import "fmt"

var AuthCodeTpl = `<html>
<style>
    body {
        background-color: #fafafa;
    }

    .paper {
        background-color: #fff;
        width: 400px;
        border: 1px solid rgba(0, 0, 0, 0.12);
        border-radius: 8px;
        padding: 20px;
        margin: auto;
    }

    .authcode {
        font-weight: bold
    }
</style>

<body>

    <div class="paper">
        <p>
            您好，如果 %v %v 不是您的 hatlonely 账户，请不要点击此邮件中的任何内容！
        </p>

        <p>
            以下是您的验证码：
        </p>
        <p class="authcode">
            %v
        </p>

        <p>
            %v %v，您好！
        </p>

        <p>
            我们收到了来自您的 hatlonely 账户的安全请求。请使用上面的验证码验证您的账号归属。
        </p>

        <p>
            请注意：该验证码将在10分钟后过期，请尽快验证！
        </p>

        <p>
            享受您的历险！
        </p>
        <p>
            hatlonely 客服团队
        </p>

    </div>
</body>

</html>`

func NewAuthCodeTpl(firstName string, lastName string, authCode string) string {
	return fmt.Sprintf(AuthCodeTpl, lastName, firstName, authCode, lastName, firstName)
}
