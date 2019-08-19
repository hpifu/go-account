#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /signout, token: "{token:str}"')
def step_impl(context, token):
    res = requests.get("{}/signout".format(context.config["url"]), params={
        "token": token,
    })
    context.status = res.status_code
    context.body = res.text
    context.cookies = res.cookies
    if context.status == 200:
        context.res = json.loads(res.text)
    print({
        "status": context.status,
        "body": context.body,
        "cookies": context.cookies,
    })


@then('检查 signout 返回包体 res.body, ok: {ok:bool}')
def step_impl(context, ok):
    assert_that(context.res["ok"], equal_to(ok))
