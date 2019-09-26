#!/usr/bin/env python3

from behave import *
from hamcrest import *
import json


@given('redis set "{key:str}"')
def step_impl(context, key):
    obj = json.loads(context.text)
    context.redis_client.set(key, json.dumps(obj))


@given('redis del "{key:str}"')
def step_impl(context, key):
    context.redis_client.delete(key)


@then('redis get "{key:str}"')
def step_impl(context, key):
    obj = json.loads(context.text)
    val = context.redis_client.get(key)
    print(val)
    result = json.loads(val)
    for key in obj:
        assert_that(result[key], equal_to(obj[key]))
