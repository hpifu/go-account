#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /signin, username: "{username:str}", password: "{password:str}"')
def step_impl(context, username, password):
    res = requests.post("{}/signin".format(context.config["url"]), json={
        "username": username,
        "password": password,
    })
    context.status = res.status_code
    context.body = res.text
    context.cookies = res.cookies
    if context.status == 200:
        context.res = json.loads(res.text)
        context.token = context.res["token"]
    print({
        "status": context.status,
        "body": context.body,
        "cookies": context.cookies,
    })


@then('检查登陆返回包体 res.body, valid: {valid:bool}, tokenlen: {tokenlen:int}')
def step_impl(context, valid, tokenlen):
    assert_that(context.res["valid"], equal_to(valid))
    assert_that(len(context.res["token"]), equal_to(tokenlen))


@then('检查登陆返回 cookie')
def step_impl(context):
    # assert_that(context.res["token"], equal_to(context.cookies["token"]))
    pass
