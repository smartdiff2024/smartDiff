# used in Windows

import os
import sys
from envi_variable import set_variable
from split.splitjson import split_json2bytecode

# set the path of IDA and diaphora
ida_path = ""
diaphora_path = ""


def exec_diff(dir_path):
  dirlist = os.listdir(dir_path)
  if "ori.bytecode" not in dirlist or "modi.bytecode" not in dirlist:
    return
  if "result.db" in dirlist:
    return
  
  file1_path = dir_path + "\\" + "ori.bytecode"
  file2_path = dir_path + "\\" + "modi.bytecode"
  file1_export_path = dir_path + "\\" + "ori.db"
  file2_export_path = dir_path + "\\" + "modi.db"

  set_variable.export_with_IDA(file1_export_path)
  os.system('\"'  + ida_path  + '\"' + " -A -c -S" + diaphora_path + " " + file1_path)


  set_variable.export_with_IDA(file2_export_path)
  os.system('\"'  + ida_path  + '\"' + " -A -c -S" + diaphora_path + " " + file2_path)

  set_variable.diff_in_batch_mode(dir_path)
  os.system("python " + diaphora_path)

  set_variable.reset_variable()

if __name__ == "__main__":
  set_variable.set_timeout_limit(300)
  json_dir = ""
  bytecode_dir = ""
  split_order_csv = ""
  if not os.path.exists(bytecode_dir):
    os.mkdir(bytecode_dir)
  # root work space path
  split_json2bytecode(json_dir, bytecode_dir, split_order_csv)
  os.chdir(bytecode_dir)
  proxy_addr_list = os.listdir(os.getcwd())
  for addr in proxy_addr_list:
    proxy_dir = bytecode_dir + "\\" + addr
    impl_list = os.listdir(proxy_dir)
    for impl in impl_list:
      impl_dir = bytecode_dir + "\\" + addr + "\\" + impl
      exec_diff(impl_dir)
      
