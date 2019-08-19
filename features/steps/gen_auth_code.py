#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /genauthcode, type: "email", email: "{email:str}", firstname: "{firstname:str}", lastname: "{lastname:str}"')
def step_impl(context, email, firstname, lastname):
    res = requests.post("{}/genauthcode".format(context.config["url"]), json={
        "type": "email",
        "email": email,
        "firstName": firstname,
        "lastName": lastname,
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


@when('请求 /genauthcode, type: "phone", phone: "{phone}", firstname: "{firstname:str}", lastname: "{lastname:str}"')
def step_impl(context, phone, firstname, lastname):
    res = requests.post("{}/genauthcode".format(context.config["url"]), json={
        "type": "phone",
        "phone": phone,
        "firstName": firstname,
        "lastName": lastname,
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


@then('检查 genauthcode 返回包体 res.body, ok: {ok:bool}')
def step_impl(context, ok):
    assert_that(context.res["ok"], equal_to(ok))
