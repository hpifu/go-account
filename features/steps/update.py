#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /update, firstName: "{firstname:str}", lastName: "{lastname:str}"')
def step_impl(context, firstname, lastname):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "name",
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


@when('请求 /update, phone: "{phone:str}"')
def step_impl(context, phone):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "phone",
        "phone": phone,
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


@when('请求 /update, email: "{email:str}"')
def step_impl(context, email):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "email",
        "email": email,
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


@when('请求 /update, birthday: "{birthday:str}"')
def step_impl(context, birthday):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "birthday",
        "birthday": birthday,
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


@when('请求 /update, gender: "{gender:int}"')
def step_impl(context, gender):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "gender",
        "gender": gender,
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


@when('请求 /update, password: "{password:str}", oldpassword: "{oldpassword:str}"')
def step_impl(context, password, oldpassword):
    res = requests.post("{}/update".format(context.config["url"]), json={
        "token": context.token,
        "field": "password",
        "password": password,
        "oldPassword": oldpassword,
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


@then('检查 update 返回包体 res.body, ok: {ok:bool}, err: "{err:str}"')
def step_impl(context, ok, err):
    assert_that(context.res["ok"], equal_to(ok))
    assert_that(context.res["err"], contains_string(err))
