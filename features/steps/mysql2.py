#!/usr/bin/env python3

from behave import *
from hamcrest import *
import json


@given('mysql 执行')
def step_impl(context):
    with context.mysql_conn.cursor() as cursor:
        cursor.execute(context.text)
    context.mysql_conn.commit()


@then('mysql 检查 "{sql:str}"')
def step_impl(context, sql):
    obj = json.loads(context.text)
    with context.mysql_conn.cursor() as cursor:
        cursor.execute(sql)
        result = cursor.fetchone()
        for key in obj:
            assert_that(result[key], equal_to(obj[key]))
