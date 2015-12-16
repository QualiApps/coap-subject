#!/usr/bin/python

from os import sep
from json import dumps
import requests


class CoapResource(object):
    __slots__ = ("_api_path", "_sub_path", "_resources")

    def __init__(self, api_url):
        self._api_path = api_url
        self._sub_path = {
            "add": "resources/add/{}/{}",
            "remove": "resources/delete/{}",
            "event": "resources/event/{}"
        }
        self._resources = {}

    def add_resource(self, name, path, observable="1"):
        '''
            Adds a new resource
            :param str name - resource name
            :param tuple path - resource path
            :param str observable - observable flag (0 | 1)
            :return bool
        '''
        status = False
        try:
            url = self._init_url("add", observable, sep.join(path))

            resp = requests.post(url)
            if resp.status_code == requests.codes.created:
                self._resources[name] = {
                    "observable": observable,
                    "path": tuple(path)
                }
                status = True
        except requests.ConnectionError as mess:
            print(mess)

        return status

    def remove_resource(self, name):
        '''
            Removes the resource by name
            :param str name - res name
            :return bool
        '''
        status = False
        try:
            path = self._resources.get(name).get("path")
            if path:
                url = self._init_url(
                    "remove",
                    sep.join(path)
                )

                resp = requests.delete(url)
                if resp.status_code == requests.codes.ok:
                    del self._resources[name]
                    status = True
        except requests.ConnectionError as mess:
            print(mess)

        return status

    def send_event(self, name, data):
        '''
            Sends an event to the coap server
            :param str name - res name
            :param dict data - payload
            :return bool
        '''
        status = False
        try:
            path = self._resources.get(name).get("path")
            if path:
                url = self._init_url(
                    "event",
                    sep.join(path)
                )

                resp = requests.post(url, dumps(data))
                if resp.status_code == requests.codes.ok:
                    status = True
        except requests.ConnectionError as mess:
            print(mess)

        return status

    def get_resources(self):
        '''
            Retrieves the all resources
            :return dict
        '''
        return self._resources

    def _init_url(self, name, *args):
        '''
            Init a full api url
            :param str name - the name of sub path (add|remove|event)
            :param tuple args - additional params
            :return str - full path
        '''
        path = self._sub_path.get(name).format(*args)
        return sep.join((self._api_path, path))
