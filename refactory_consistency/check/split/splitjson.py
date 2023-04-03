import json
import os
import csv

def split_json2bytecode(json_dir, bytecode_dir, split_order_csv):
  json_dir = json_dir + "\\"
  bytecode_dir = bytecode_dir + "\\"
  jsonlist = os.listdir(json_dir)
  listLength = len(jsonlist)
  with open(split_order_csv, encoding='utf-8') as f:
    for row in csv.reader(f, skipinitialspace=True):
      fileNum = 0
      for i in range(8, len(row) - 5, 5):
        fileNum += 1
        savePath = bytecode_dir + row[5] + '/' + str(fileNum)
        os.makedirs(savePath)
        for j in range(0, listLength, 1):
          if str.lower(row[5] + '.json') == str.lower(jsonlist[j]):
            with open(json_dir+jsonlist[j], encoding='utf-8') as a:
              file = json.load(a)
              print("len: ", len(file.items()))
              for key, value in file.items():
                print("key: ", key, "row[i]: ", row[i])
                if str.lower(key) == str.lower(row[i]) and "runtimebin" in value:
                  with open(savePath +'/ori.bytecode','w+') as f:
                    runtimebin = value['runtimebin']
                    f.write(str(runtimebin))
                if str.lower(key) == str.lower(row[i + 5]) and "runtimebin" in value:
                  with open(savePath +'/modi.bytecode','w+') as f:
                    runtimebin = value['runtimebin']
                    f.write(str(runtimebin))
