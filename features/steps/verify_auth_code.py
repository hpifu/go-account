#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /verifyauthcode, type: "email", email: "{email:str}", code: "{code:str}"')
def step_impl(context, email, code):
    res = requests.get("{}/verifyauthcode".format(context.config["url"]), params={
        "type": "email",
        "email": email,
        "code": code,
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


@when('请求 /verifyauthcode, type: "phone", phone: "{email:str}", code: "{code:str}')
def step_impl(context, phone, code):
    res = requests.get("{}/verifyauthcode".format(context.config["url"]), params={
        "type": "phone",
        "phone": phone,
        "code": code,
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


@then('检查 verifyauthcode 返回包体 res.body, ok: {ok:bool}, tip: "{tip:str}"')
def step_impl(context, ok, tip):
    assert_that(context.res["ok"], equal_to(ok))
    assert_that(context.res["tip"], equal_to(tip))
