#!/usr/bin/env python3


from behave import *
from hamcrest import *
import requests
import json


@when('http 请求 {method:str} {path:str}')
def step_impl(context, method, path):
    if context.text:
        obj = json.loads(context.text)
    else:
        obj = {}
    if "params" not in obj:
        obj["params"] = {}
    if method == "GET":
        context.res = requests.get(
            "{}{}".format(context.config["url"], path),
            params=obj["params"]
        )
    if method == "PUT":
        context.res = requests.put(
            "{}{}".format(context.config["url"], path),
            json=obj["json"]
        )


@then('http 检查 {status:int}')
def step_impl(context, status):
    res = context.res
    if context.text:
        obj = json.loads(context.text)
    else:
        obj = {}
    assert_that(res.status_code, equal_to(status))
    if "json" in obj:
        result = json.loads(res.text)
        for key in obj["json"]:
            assert_that(result[key], equal_to(obj["json"][key]))
