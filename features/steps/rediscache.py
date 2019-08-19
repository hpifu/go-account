#!/usr/bin/env python3

from behave import *
from hamcrest import *
import json


@given('rediscache.token 创建 token: "{token:str}", phone: "{phone:str}", email: "{email:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
def step_impl(context, token, phone, email, firstname, lastname, birthday, gender):
    account = {
        "phone": phone,
        "email": email,
        "firstName": firstname,
        "lastName": lastname,
        "birthday": birthday,
        "gender": gender,
    }
    context.redis_client.set(token, json.dumps(account))


@then('检查 rediscache.token, 存在记录 phone: "{phone:str}", email: "{email:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
def step_impl(context, phone, email, firstname, lastname, birthday, gender):
    res = context.redis_client.get(context.token)
    account = json.loads(res)
    assert_that(account["phone"], equal_to(phone))
    assert_that(account["email"], equal_to(email))
    assert_that(account["firstName"], equal_to(firstname))
    assert_that(account["lastName"], equal_to(lastname))
    assert_that(account["birthday"], equal_to(birthday))
    assert_that(account["gender"], equal_to(gender))


@then('检查 rediscache.token, 不存在记录 token: "{token:str}"')
def step_impl(context, token):
    res = context.redis_client.get(token)
    print(res)
    assert_that(res, equal_to(None))


@given('rediscache.authcode 创建 key: "{key:str}", code: "{code:str}"')
def step_impl(context, key, code):
    context.redis_client.set("ac_"+key, code)


@then('检查 rediscache.authcode, 存在记录 key: "{key:str}"')
def step_impl(context, key):
    res = context.redis_client.get("ac_" + key)
    print(res)
    assert_that(len(res) == 6)
