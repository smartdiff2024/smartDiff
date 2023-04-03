import csv

import sys
sys.path.append("./util")
from util import findReplayError, interesting, trace
from util.path import *
from util.json_read import json_data_read
from util.findReplayError import *
from util import Errors
from jsonlines import InvalidLineError

keys = []
baseResult = []
errorFiles = []
ValidNum = 0

def resultAnalysis(result):
  global baseResult, resultwriter
  for key, value in list(result.items()):    #type(value) -> dict
    if key == "dictionary_item_removed":
      for slot, slotvalue in value.items():
        writeResult = baseResult.copy()
        writeResult.append(key)
        writeResult.append(str.split(slot, "'")[1])
        writeResult.append(slotvalue)
        writeResult.append(" ")
        resultwriter.writerow(writeResult)
        if len(writeResult) != 9:
          raise RuntimeError("Error writeResult!")
    elif key == "dictionary_item_added":   
      for slot, slotvalue in value.items():
        writeResult = baseResult.copy()
        writeResult.append(key)
        writeResult.append(str.split(slot, "'")[1])
        writeResult.append("")
        writeResult.append(slotvalue)
        resultwriter.writerow(writeResult)
        if len(writeResult) != 9:
          raise RuntimeError("Error writeResult!")
    elif key == "values_changed":
      for slot, slotvalue in value.items():
        writeResult = baseResult.copy()
        writeResult.append(key)
        writeResult.append(str.split(slot, "'")[1])
        writeResult.append(slotvalue["old_value"])
        writeResult.append(slotvalue["new_value"])
        resultwriter.writerow(writeResult)
        if len(writeResult) != 9:
          raise RuntimeError("Error writeResult!")
    else:
      raise RuntimeError("new diff type: ", key)
          
def addBaseResult(json, path):
  global baseResult
  proxyaddress = some_path(path, -2)
  baseResult.append(proxyaddress)
  baseResult.append(str(int(json["txblocknum"], 16)))
  if "txpositionnum" not in json:
    json["txpositionnum"] = 0
  baseResult.append(str(json["txpositionnum"]))
  baseResult.append(json["preimpladdress"])
  baseResult.append(json["postimpladdress"])

from deepdiff import DeepDiff  # pip install deepdiff
import os
def analy(json1, json2, dirpath):
  global baseResult, resultwriter
  baseResult = []
  addBaseResult(json2, dirpath)
  if ("result" in json1) ^ ("result" in json2):
    baseResult.append("warning: trace result different in path:  " + dirpath)
    resultwriter.writerow(baseResult)
    return True
  StorageResult = DeepDiff(json1["storage"], json2["storage"] , ignore_order=True, verbose_level=2, cutoff_distance_for_pairs=1)
  ProxyStorageResult =  DeepDiff(json1["proxystorage"], json2["proxystorage"] , ignore_order=True, verbose_level=2, cutoff_distance_for_pairs=1)
  StorageResultLength, ProxyStorageResultLength = len(StorageResult), len(ProxyStorageResult)
  hasDiff = False
  if StorageResultLength != 0:
    resultAnalysis(StorageResult)
    hasDiff = True
    
  if ProxyStorageResultLength != 0:
    resultAnalysis(ProxyStorageResult)
    hasDiff = True
  return hasDiff

def readJsonS(dirpath):
  print("jsonpath: ", dirpath)
  global ValidNum
  files = os.listdir(dirpath)
  if len(files) != 2:
    if dirpath not in findReplayError.ErrorFiles:
      proxyaddress = some_path(dirpath, -2)
      notAll = True
      try:
        if ifNoBytecode(proxyaddress, dirpath, "NoTwoFile"):
          notAll = False
        if interesting.onlyOrigin(dirpath):
          notAll = False
        if trace.errorFirstImpe(proxyaddress):
          notAll = False
      except FileNotFoundError as e:
        print("FileNotFoundError: ", dirpath)
      finally:
        if notAll:
          print("Warning: not two files in dirpath: ", dirpath)
    return
  json1 = json_data_read(dirpath + "/" + "origin.json")
  json2 = json_data_read(dirpath + "/" + "modify.json")
  jsonLength = len(json2) if len(json2) < len(json1) else len(json1)
  if len(json2) != len(json1):
    if dirpath not in findReplayError.ErrorFiles:
      proxyaddress = some_path(dirpath, -2)
      if not ifNoBytecode(proxyaddress, dirpath, "DtraceLength"):
        if not interesting.moreModify(len(json1), len(json2), dirpath):
          print("Warning: Different Dtrace Length", dirpath)
    return
  ValidNum += 1
  hasError = False
  for index in range(0, jsonLength):
    if analy(json1[index], json2[index], dirpath):
      hasError = True
  if not hasError:
    print("exactly the same: ", dirpath)
      

import os
def check_if_dir(file_path):
    temp_list = os.listdir(file_path)    #put file name from file_path in temp_list
    for temp_list_each in temp_list:
        if os.path.isfile(file_path + '/' + temp_list_each):
            temp_path = file_path + '/' + temp_list_each
            if os.path.splitext(temp_path)[-1] == '.json': 
                try:
                  readJsonS(file_path)
                except InvalidLineError:
                  print("skip path: ", file_path)
                  break
                break
            else:
                continue
        else:
            check_if_dir(file_path + '/' + temp_list_each)    #loop traversal
if __name__ == "__main__":
  trace.readtrace("order.csv")
  result = open("result.csv", "w", encoding="utf-8")
  resultwriter = csv.writer(result)
  header = ["proxy address", "blocknum", "positionnum", "ori impl", "modi impl", \
    "diff type", "slot", "ori value", "modi value"]
  findReplayError.check_if_dir("./log directory")
  print("log errorFiles: ", findReplayError.ErrorFiles, "length: ", len(findReplayError.ErrorFiles))
  resultwriter.writerow(header)
  check_if_dir("./result directory")
  print("dtrace errorFiles: ", findReplayError.ErrorFiles, "length: ", len(findReplayError.ErrorFiles))
  print("ValidNum: ", ValidNum)
  result.close()