import cStringIO
import base64
import copy
import json
import httplib
import time
import re
import string
from django.conf import settings
from .states import JobState

POD_TEMPLATE = '''{
   "kind":"ReplicationController",
   "apiVersion":"$version",
   "metadata":{
      "name":"$id",
      "labels":{
         "name":"$id"
      }
   },
   "spec":{
      "replicas":1,
      "selector":{
         "name":"$id"
      },
      "template":{
         "name":"$id",
         "metadata":{
            "labels":{
               "name":"$id",
               "tier":"frontend",
               "environment":"production"
            }
         },
         "spec":{
            "containers":[
               {
                  "args":[
                    "start",
                    "web"
                  ],
                  "name":"$name",
                  "image":"$image",
                  "ports": [
                      {
                        "protocol": "TCP",
                        "containerPort": 5000
                      }
                  ],
                  "terminationMessagePath": "/dev/termination-log",
                  "imagePullPolicy": "IfNotPresent",
                  "capabilities": {}
               }
            ],
            "restartPolicy": "Always",
            "dnsPolicy": "ClusterFirst"
         }
      }
   }
}'''

RETRIES = 3
MATCH = re.compile(
    r'(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)?.(?P<c_num>[0-9]+)')

class KubeHTTPClient():

    def __init__(self, target, auth, options, pkey):
        self.target = settings.K8S_MASTER
        self.port = "8080"
        self.registry = settings.REGISTRY_HOST+":"+settings.REGISTRY_PORT
        self.apiversion = "v1beta3"
        self.conn = httplib.HTTPConnection(self.target+":"+self.port)
        #self.container_state = ""

    # container api

    def create(self, name, image, command, **kwargs):
        #self.container_state = "create"
        justName=name.split('.')[0].split('_')[0]
        l = {}
        l["id"]=justName
        l["version"]=self.apiversion
        l["image"]=self.registry+"/"+image
        l["name"]=l['id'].replace(".","-")
        template=string.Template(POD_TEMPLATE).substitute(l)
        js_template = json.loads(template)
        loc = locals().copy()
        loc.update(re.match(MATCH, name).groupdict())
        mem = kwargs.get('memory', {}).get(loc['c_type'])
        cpu = kwargs.get('cpu', {}).get(loc['c_type'])
        if mem or cpu :
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"] = {"limits":{}}
        if mem:
            mem = mem.lower()
            if mem[-2:-1].isalpha() and mem[-1].isalpha():
                mem = mem[:-1]
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"]["limits"]["memory"] = mem
        if cpu:
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"]["limits"]["cpu"] = cpu
        headers = {'Content-Type': 'application/json'}
        self.conn.request('POST', '/api/'+self.apiversion+'/namespaces/default/replicationcontrollers',
                  headers=headers, body=json.dumps(js_template))
        resp = self.conn.getresponse()
        data = resp.read()
        if not 200 <= resp.status <= 299:
            errmsg = "Failed to retrieve unit: {} {} - {}".format(
                resp.status, resp.reason, data)
            raise RuntimeError(errmsg)

    def start(self, name):
        """
        Start a container
        """
        #self.container_state = "start"
        actual_pod = {}
        for _ in xrange(30):
            self.conn.request('GET','/api/'+self.apiversion+'/namespaces/default/pods')
            resp = self.conn.getresponse()
            parsed_json =  json.loads(resp.read())

            for pod in parsed_json['items']:
                if pod['metadata']['generateName'] == name.replace("_",".")+'-':
                    actual_pod = pod
                    break
            if actual_pod and actual_pod['status']['phase'] == 'Running':
                return
            time.sleep(1)

    def stop(self, name):
        """
        Stop a container
        """
        return

    def destroy(self, name):
        """
        Destroy a container
        """
        return

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        # dump input into a json object for testing purposes
        return 0, json.dumps({'name': name,
                              'image': image,
                              'entrypoint': entrypoint,
                              'command': command})

    def _get_pod_state(self, name):
        try:
            self.conn.request('GET','/api/'+self.apiversion+'/namespaces/default/pods')
            resp = self.conn.getresponse()
            parsed_json =  json.loads(resp.read())
            actual_pod = {}
            for pod in parsed_json['items']:
                if pod['metadata']['generateName'] == name.replace("_",".")+'-':
                    actual_pod = pod
                    break
            if  actual_pod.get('status') and (actual_pod['status']['phase'] == 'Running' or actual_pod['status']['phase'] == 'Pending'):
                return JobState.created
            else:
                return JobState.up
        except:
            return JobState.destroyed

    def state(self, name):
        try:
            for _ in xrange(30):
                return self._get_pod_state(name)
                time.sleep(1)
            # FIXME (smothiki): should be able to send JobState.crashed
        except KeyError:
            return JobState.error
        except RuntimeError:
            return JobState.destroyed

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

SchedulerClient = KubeHTTPClient
