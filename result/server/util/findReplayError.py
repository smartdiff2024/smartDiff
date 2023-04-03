import os
from path import some_path
from bytecode import ifNoBytecode, ErrorFiles, jsonBasepath

def ifError(path):
  logfiles = os.listdir(path)
  if len(logfiles) > 1:
      print("multi files::path: ", path, "logfiles: ", logfiles)
  for logfile in logfiles:
    with open(path + "/" + logfile, "r", encoding="utf-8") as log:
        log.seek(0)
        lines = log.readlines()
        for line in lines:
            if "Something wrong about" in line:
                _, proxyAddress = os.path.split(path)
                if not ifNoBytecode(proxyAddress, path + "/" + logfile, "Log"):
                    print("Log error path: ", path + "/" + logfile)
                    if some_path(path, -1) not in ErrorFiles:
                        ErrorFiles.append(some_path(path, -1))
                break
    
path = ""
def check_if_dir(file_path):
    temp_list = os.listdir(file_path)    #put file name from file_path in temp_list
    for temp_list_each in temp_list:
        if os.path.isfile(file_path + '/' + temp_list_each):
            temp_path = file_path + '/' + temp_list_each
            if os.path.splitext(temp_path)[-1] == '.txt': 
                ifError(file_path)
                break
            else:
                continue
        else:
            check_if_dir(file_path + '/' + temp_list_each)    #loop traversal
            
if __name__ == "__main__":
    check_if_dir("./log2")
    print("ErrorFiles: ", ErrorFiles)