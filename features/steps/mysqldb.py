#!/usr/bin/env python3

from behave import *
from hamcrest import *


@given('mysqldb.accounts 创建用户, phone: "{phone:str}", email: "{email:str}", password: "{password:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
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
    with context.mysql_conn.cursor() as cursor:
        cursor.execute(context.cleanup["sql"])
        cursor.execute(
            "INSERT INTO accounts (phone, email, password, first_name, last_name, birthday, gender) VALUES (%s, %s, %s, %s, %s, %s, %s)",
            (phone, email, password, firstname, lastname, birthday, gender)
        )
    context.mysql_conn.commit()


@given('mysqldb.accounts 删除用户, email: "{email:str}"')
def step_impl(context, email):
    context.cleanup = {
        "sql": "DELETE FROM accounts WHERE email='{}'".format(
            email
        )
    }
    with context.mysql_conn.cursor() as cursor:
        cursor.execute(context.cleanup["sql"])
    context.mysql_conn.commit()


@then('检查 mysqldb.accounts, 存在记录 phone: "{phone:str}", email: "{email:str}", password: "{password:str}", firstname: "{firstname:str}", lastname: "{lastname:str}", birthday: "{birthday:str}", gender: {gender:int}')
def step_impl(context, phone, email, password, firstname, lastname, birthday, gender):
    with context.mysql_conn.cursor() as cursor:
        cursor.execute(
            "SELECT * FROM accounts WHERE email='{}'".format(email)
        )
        res = cursor.fetchall()
        assert_that(len(res), equal_to(1))
        account = res[0]
        assert_that(phone, equal_to(account["phone"]))
        assert_that(email, equal_to(account["email"]))
        assert_that(password, equal_to(account["password"]))
        assert_that(firstname, equal_to(account["first_name"]))
        assert_that(lastname, equal_to(account["last_name"]))
        assert_that(birthday, equal_to(
            account["birthday"].strftime("%Y-%m-%d")))
        assert_that(gender, equal_to(account["gender"]))
