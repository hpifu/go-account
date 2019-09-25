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
