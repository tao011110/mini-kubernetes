from django.http import HttpResponse
import json
from . import handler


def process(request):
    body = json.loads("{}")
    if len(request.body):
        body = json.loads(request.body.decode(encoding="utf-8", errors="strict"))
    ret = handler.handler(request.GET, body)
    return HttpResponse(ret)
