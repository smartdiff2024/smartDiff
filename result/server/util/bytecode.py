import json
import os
from path import some_path
from Errors import ErrorFiles, jsonBasepath
import Errors


def ifNoBytecode(proxyAddress, logfile, caller):
  print("jsonBasepath: ", Errors.jsonBasepath)
  jsons = os.listdir(jsonBasepath)
  proxyaddressname = ""
  for jsonname in jsons:
    if str.lower(proxyAddress + ".json") == str.lower(jsonname):
      proxyaddressname = jsonname
      break
  with open(jsonBasepath + "/" + proxyaddressname, "r", encoding="utf-8") as jsonFile:
    data = json.load(jsonFile)
    data_keys = data.keys()   
    if len(data_keys) < 2:
      if some_path(logfile, -2) not in ErrorFiles:
        ErrorFiles.append(some_path(logfile, -2))
      return True
    for key in list(data.keys()):
        if "createbin" not in data[key]:
          if caller ==  "Log":
            if some_path(logfile, -2) not in ErrorFiles:
              ErrorFiles.append(some_path(logfile, -2))
            return True
          elif caller == "NoTwoFile":
            if some_path(logfile, -2) not in ErrorFiles:
              ErrorFiles.append(some_path(logfile, -2))
            return True
          elif caller == "DtraceLength":
            if some_path(logfile, -2) not in ErrorFiles:
              ErrorFiles.append(some_path(logfile, -2))
            return True
          else:
            raise RuntimeError("Unknown Caller")
  return False
  
if __name__ == "__main__":
  jsonBasepath = "///"
  pass