from Errors import ErrorFiles
from json_read import json_data_read
from deepdiff import DeepDiff
from path import some_path

def onlyOrigin(dirpath):
  json = json_data_read(dirpath + "/" + "origin.json")
  jsonLength = len(json)
  if jsonLength % 2 != 0:
    return False
  for index in range(0, jsonLength, 2):
    StorageResult = DeepDiff(json[index], json[index + 1] , ignore_order=True, verbose_level=2, cutoff_distance_for_pairs=1)
    if len(StorageResult) != 0:
      return False
  proxyaddress = some_path(dirpath, -2)
  if proxyaddress not in ErrorFiles:
    ErrorFiles.append(proxyaddress)
  return True

def moreModify(origin, modify, dirpath):
  if modify > origin:
    proxyaddress = some_path(dirpath, -2)
    if proxyaddress not in ErrorFiles:
      ErrorFiles.append(proxyaddress)
    return True
  return False