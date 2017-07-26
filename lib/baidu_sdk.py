# -*- coding: utf-8 -*-

from aip import AipSpeech
import sys
import json

# python ../lib/baidu_sdk.py 9926043 dE3SNoLcH5KuhhwVc8hhL5fL-
# fwMleRZAsW9vmMf2GcmlVexXNooTVHGa ~/xxx/xxx.amr
APP_ID = sys.argv[1]
API_KEY = sys.argv[2]
SECRET_KEY = sys.argv[3]
path = sys.argv[4]

aipSpeech = AipSpeech(APP_ID, API_KEY, SECRET_KEY)


def get_file_content(filePath):
    with open(filePath, 'rb') as fp:
        return fp.read()


result = aipSpeech.asr(get_file_content(path), 'amr', 8000, {
    'lan': 'zh',
})

result = json.dumps(result)
print result.encode('ascii', 'ignore')
