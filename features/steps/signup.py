#!/usr/bin/env python3

from behave import *
from hamcrest import *
import requests
import json


@when('请求 /signup, phone: "{phone:str}", email: "{email:str}", password: "{password:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
def step_impl(context, phone, email, password, firstname, lastname, birthday, gender):
    context.cleanup = {
        "sql": "DELETE FROM accounts WHERE phone='{}' OR email='{}'".format(
            phone, email
        )
    }
    context.phone = phone
    context.email = email
    context.password = password
    context.firstname = firstname
    context.lastname = lastname
    context.birthday = birthday
    context.gender = gender
    res = requests.post("{}/signup".format(context.config["url"]), json={
        "phone": phone,
        "email": email,
        "password": password,
        "firstName": firstname,
        "lastName": lastname,
        "birthday": birthday,
        "gender": gender
    })
    context.status = res.status_code
    context.body = res.text
    print(res)
    if context.status == 200:
        context.res = json.loads(res.text)
    print({
        "status": context.status,
        "body": context.body,
    })


@then('检查注册返回包体 res.body, success: {success:bool}')
def step_impl(context, success):
    assert_that(context.res["success"], equal_to(success))
