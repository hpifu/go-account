#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /verify, field: "{field:str}", value: "{value:str}"')
def step_impl(context, field, value):
    res = requests.get("{}/verify".format(context.config["url"]), params={
        "field": field,
        "value": value,
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


@then('检查 verify 返回包体 res.body, ok: {ok:bool}, tip: "{tip:str}"')
def step_impl(context, ok, tip):
    assert_that(context.res["ok"], equal_to(ok))
    assert_that(context.res["tip"], equal_to(tip))
