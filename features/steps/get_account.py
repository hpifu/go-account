#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /getaccount, token: "{token:str}"')
def step_impl(context, token):
    res = requests.get("{}/getaccount".format(context.config["url"]), params={
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


@then('检查 token 返回包体 res.body, ok: {ok:bool}, phone: "{phone:str}", email: "{email:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
def step_impl(context, ok, phone, email, firstname, lastname, birthday, gender):
    assert_that(context.res["ok"], equal_to(ok))
    assert_that(context.res["account"]["phone"], equal_to(phone))
    assert_that(context.res["account"]["email"], equal_to(email))
    assert_that(context.res["account"]["firstName"], equal_to(firstname))
    assert_that(context.res["account"]["lastName"], equal_to(lastname))
    assert_that(context.res["account"]["birthday"], equal_to(birthday))
    assert_that(context.res["account"]["gender"], equal_to(gender))


@then('检查 token 返回包体 res.body, ok: false')
def step_impl(context):
    assert_that(context.res["ok"], equal_to(False))
