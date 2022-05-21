from django.contrib import admin
from django.urls import path
from . import responce


urlpatterns = [
    path('admin/', admin.site.urls),
    path('', responce.process)
]
